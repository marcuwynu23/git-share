package git

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func newRequest(method, path string) *http.Request {
	return &http.Request{
		Method: method,
		URL:    &url.URL{Path: path, RawQuery: ""},
		Proto:  "HTTP/1.1",
		Header: http.Header{},
	}
}

func TestBuildCGIEnv(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		gitDir   string
		repoName string
		isBare   bool
		checks   map[string]string
	}{
		{
			name:     "non-bare with repo prefix",
			path:     "/myrepo.git/info/refs",
			gitDir:   "/home/user/projects/myrepo/.git",
			repoName: "myrepo",
			isBare:   false,
			checks: map[string]string{
				"GIT_PROJECT_ROOT": "/home/user/projects/myrepo",
				"PATH_INFO":        "/info/refs",
				"SCRIPT_NAME":      "/myrepo.git",
			},
		},
		{
			name:     "non-bare without repo prefix",
			path:     "/info/refs",
			gitDir:   "/home/user/projects/myrepo/.git",
			repoName: "myrepo",
			isBare:   false,
			checks: map[string]string{
				"GIT_PROJECT_ROOT": "/home/user/projects/myrepo",
				"PATH_INFO":        "/info/refs",
			},
		},
		{
			name:     "non-bare HEAD request",
			path:     "/HEAD",
			gitDir:   "/home/user/projects/myrepo/.git",
			repoName: "myrepo",
			isBare:   false,
			checks: map[string]string{
				"PATH_INFO": "/HEAD",
			},
		},
		{
			name:     "non-bare git-upload-pack",
			path:     "/myrepo.git/git-upload-pack",
			gitDir:   "/home/user/projects/myrepo/.git",
			repoName: "myrepo",
			isBare:   false,
			checks: map[string]string{
				"PATH_INFO": "/git-upload-pack",
			},
		},
		{
			name:     "bare with repo prefix",
			path:     "/myrepo.git/info/refs",
			gitDir:   "/home/user/repos/myrepo.git",
			repoName: "myrepo",
			isBare:   true,
			checks: map[string]string{
				"GIT_PROJECT_ROOT": "/home/user/repos",
				"PATH_INFO":        "/myrepo.git/info/refs",
			},
		},
		{
			name:     "bare without repo prefix",
			path:     "/info/refs",
			gitDir:   "/home/user/repos/myrepo.git",
			repoName: "myrepo",
			isBare:   true,
			checks: map[string]string{
				"PATH_INFO": "/myrepo.git/info/refs",
			},
		},
		{
			name:     "with /git prefix stripped",
			path:     "/git/myrepo.git/info/refs",
			gitDir:   "/home/user/projects/myrepo/.git",
			repoName: "myrepo",
			isBare:   false,
			checks: map[string]string{
				"PATH_INFO": "/info/refs",
			},
		},
		{
			name:     "non-bare path with repo.git suffix exact match",
			path:     "/myrepo.git",
			gitDir:   "/home/user/projects/myrepo/.git",
			repoName: "myrepo",
			isBare:   false,
			checks: map[string]string{
				"PATH_INFO": "/",
			},
		},
		{
			name:     "non-bare windows paths",
			path:     "/myrepo.git/info/refs",
			gitDir:   `C:\Users\test\projects\myrepo\.git`,
			repoName: "myrepo",
			isBare:   false,
			checks: map[string]string{
				"GIT_PROJECT_ROOT": `C:\Users\test\projects\myrepo`,
				"PATH_INFO":        "/info/refs",
			},
		},
		{
			name:     "bare with repo prefix windows",
			path:     "/myrepo.git/info/refs",
			gitDir:   `C:\Users\test\repos\myrepo.git`,
			repoName: "myrepo",
			isBare:   true,
			checks: map[string]string{
				"PATH_INFO": "/myrepo.git/info/refs",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRequest("GET", tt.path)
			env := buildCGIEnv(r, tt.gitDir, tt.repoName, tt.isBare, false)

			envMap := make(map[string]string)
			for _, e := range env {
				if strings.HasPrefix(e, "GIT_") || strings.HasPrefix(e, "PATH_") || strings.HasPrefix(e, "SCRIPT_") || strings.HasPrefix(e, "REQUEST_") || strings.HasPrefix(e, "CONTENT_") || strings.HasPrefix(e, "SERVER_") || strings.HasPrefix(e, "QUERY_") {
					parts := strings.SplitN(e, "=", 2)
					if len(parts) == 2 {
						envMap[parts[0]] = parts[1]
					}
				}
			}

			for key, expected := range tt.checks {
				got, ok := envMap[key]
				if !ok {
					t.Errorf("missing env var %s", key)
					continue
				}
				// Normalize path separators for comparison
				gotNorm := filepath.ToSlash(got)
				expectedNorm := filepath.ToSlash(expected)
				if gotNorm != expectedNorm {
					t.Errorf("%s = %s, want %s", key, got, expected)
				}
			}
		})
	}
}

func TestBuildCGIEnvStandardVars(t *testing.T) {
	r := newRequest("POST", "/repo.git/git-upload-pack")
	r.Header.Set("Content-Type", "application/x-git-upload-pack-request")
	r.Header.Set("Content-Length", "100")
	r.URL.RawQuery = "service=git-upload-pack"

	env := buildCGIEnv(r, "/home/user/projects/repo/.git", "repo", false, false)
	envMap := make(map[string]string)
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	checks := map[string]string{
		"REQUEST_METHOD":   "POST",
		"CONTENT_TYPE":     "application/x-git-upload-pack-request",
		"CONTENT_LENGTH":   "100",
		"QUERY_STRING":     "service=git-upload-pack",
		"SERVER_PROTOCOL":  "HTTP/1.1",
		"SERVER_SOFTWARE":  "git-share/0.1.0",
		"GIT_HTTP_EXPORT_ALL": "1",
		"GIT_TERMINAL_PROMPT": "0",
	}

	for key, expected := range checks {
		got, ok := envMap[key]
		if !ok {
			t.Errorf("missing env var %s", key)
			continue
		}
		if got != expected {
			t.Errorf("%s = %s, want %s", key, got, expected)
		}
	}

	if _, ok := envMap["GIT_NAMESPACE"]; ok {
		t.Error("GIT_NAMESPACE should not be set")
	}
}

func TestHandleSmartHTTPNonBare(t *testing.T) {
	if err := CheckGit(); err != nil {
		t.Skip("git not found in PATH")
	}

	dir := t.TempDir()

	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}

	run("init")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")
	run("commit", "--allow-empty", "-m", "initial")
	run("tag", "v1")

	gitDir := filepath.Join(dir, ".git")
	info := &RepoInfo{
		Root:     dir,
		Worktree: dir,
		Branch:   "main",
		IsBare:   false,
		GitDir:   gitDir,
	}

	tests := []struct {
		name       string
		method     string
		path       string
		query      string
		wantStatus int
	}{
		{
			name:       "HEAD",
			method:     "GET",
			path:       "/repo.git/HEAD",
			wantStatus: 200,
		},
		{
			name:       "info/refs without service",
			method:     "GET",
			path:       "/repo.git/info/refs",
			wantStatus: 200,
		},
		{
			name:       "info/refs with upload-pack",
			method:     "GET",
			path:       "/repo.git/info/refs",
			query:      "service=git-upload-pack",
			wantStatus: 200,
		},
		{
			name:       "info/refs with receive-pack",
			method:     "GET",
			path:       "/repo.git/info/refs",
			query:      "service=git-receive-pack",
			wantStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Method: tt.method,
				URL:    &url.URL{Path: tt.path, RawQuery: tt.query},
				Proto:  "HTTP/1.1",
				Header: http.Header{"User-Agent": []string{"git/2.x"}},
				Body:   http.NoBody,
			}
			w := httptest.NewRecorder()
			HandleSmartHTTP(info, w, req, "repo", false)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("status = %d, want %d\nbody: %s", resp.StatusCode, tt.wantStatus, body)
			}
		})
	}
}

func TestHandleSmartHTTPNonBareViaGitDir(t *testing.T) {
	if err := CheckGit(); err != nil {
		t.Skip("git not found in PATH")
	}

	dir := t.TempDir()

	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}

	run("init")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")
	run("commit", "--allow-empty", "-m", "initial")

	gitDir := filepath.Join(dir, ".git")
	info := &RepoInfo{
		Root:     dir,
		Worktree: dir,
		Branch:   "main",
		IsBare:   false,
		GitDir:   gitDir,
	}

	// Test using PATH_INFO = /.git/<action> with GIT_PROJECT_ROOT = working tree
	t.Run("info/refs with receive-pack via /.git prefix", func(t *testing.T) {
		req := &http.Request{
			Method: "GET",
			URL:    &url.URL{Path: "/repo.git/info/refs", RawQuery: "service=git-receive-pack"},
			Proto:  "HTTP/1.1",
			Header: http.Header{"User-Agent": []string{"git/2.x"}},
			Body:   http.NoBody,
		}
		w := httptest.NewRecorder()
		HandleSmartHTTP(info, w, req, "repo", false)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		t.Logf("status = %d, body = %s", resp.StatusCode, body)
		if resp.StatusCode != 200 {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})
}

func TestHandleSmartHTTPBare(t *testing.T) {
	if err := CheckGit(); err != nil {
		t.Skip("git not found in PATH")
	}

	dir := t.TempDir()
	bareDir := filepath.Join(dir, "repo.git")

	execGit := func(workdir string, args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = workdir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v in %s failed: %v\n%s", args, workdir, err, out)
		}
	}

	execGit(dir, "init", "--bare", bareDir)

	cloneDir := filepath.Join(dir, "clone")
	execGit(dir, "clone", bareDir, cloneDir)
	execGit(cloneDir, "config", "user.email", "test@test.com")
	execGit(cloneDir, "config", "user.name", "Test")
	execGit(cloneDir, "commit", "--allow-empty", "-m", "initial")
	execGit(cloneDir, "push", "origin", "main")
	execGit(cloneDir, "tag", "v1")
	execGit(cloneDir, "push", "origin", "--tags")

	info := &RepoInfo{
		Root:     bareDir,
		Worktree: bareDir,
		Branch:   "main",
		IsBare:   true,
		GitDir:   bareDir,
	}

	tests := []struct {
		name       string
		method     string
		path       string
		query      string
		wantStatus int
	}{
		{
			name:       "HEAD",
			method:     "GET",
			path:       "/repo.git/HEAD",
			wantStatus: 200,
		},
		{
			name:       "info/refs with receive-pack",
			method:     "GET",
			path:       "/repo.git/info/refs",
			query:      "service=git-receive-pack",
			wantStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Method: tt.method,
				URL:    &url.URL{Path: tt.path, RawQuery: tt.query},
				Proto:  "HTTP/1.1",
				Header: http.Header{"User-Agent": []string{"git/2.x"}},
				Body:   http.NoBody,
			}
			w := httptest.NewRecorder()
			HandleSmartHTTP(info, w, req, "repo", false)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("status = %d, want %d\nbody: %s", resp.StatusCode, tt.wantStatus, body)
			}
		})
	}
}
