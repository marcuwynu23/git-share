package git_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/marcuwynu23/git-share/internal/git"
)

func TestCheckGit(t *testing.T) {
	err := git.CheckGit()
	if err != nil {
		t.Skip("git not found in PATH")
	}
}

func TestFindRepository(t *testing.T) {
	if err := git.CheckGit(); err != nil {
		t.Skip("git not found in PATH")
	}

	repo, err := git.FindRepository()
	if err != nil {
		t.Fatalf("FindRepository() error = %v", err)
	}
	if repo.Root == "" {
		t.Error("Root should not be empty")
	}
	if repo.Branch == "" {
		t.Error("Branch should not be empty")
	}
	if repo.GitDir == "" {
		t.Error("GitDir should not be empty")
	}
}

func TestValidate(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	os.MkdirAll(gitDir, 0755)

	repo := &git.RepoInfo{GitDir: gitDir}
	if err := git.Validate(repo); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestValidateNonExistent(t *testing.T) {
	repo := &git.RepoInfo{GitDir: "/nonexistent/path/.git"}
	if err := git.Validate(repo); err == nil {
		t.Error("Validate() should error on non-existent path")
	}
}

func TestValidateNonDir(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "notadir")
	os.WriteFile(file, []byte("data"), 0644)

	repo := &git.RepoInfo{GitDir: file}
	if err := git.Validate(repo); err == nil {
		t.Error("Validate() should error on non-directory")
	}
}

func TestHandleSmartHTTPNonBare(t *testing.T) {
	if err := git.CheckGit(); err != nil {
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
	info := &git.RepoInfo{
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
			git.HandleSmartHTTP(info, w, req, "repo", false)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("status = %d, want %d\nbody: %s", resp.StatusCode, tt.wantStatus, body)
			}
		})
	}
}

func TestHandleSmartHTTPBare(t *testing.T) {
	if err := git.CheckGit(); err != nil {
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

	info := &git.RepoInfo{
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
			git.HandleSmartHTTP(info, w, req, "repo", false)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("status = %d, want %d\nbody: %s", resp.StatusCode, tt.wantStatus, body)
			}
		})
	}
}

func TestErrorValues(t *testing.T) {
	if git.ErrGitNotFound == nil {
		t.Error("ErrGitNotFound should not be nil")
	}
	if git.ErrNotARepo == nil {
		t.Error("ErrNotARepo should not be nil")
	}
	if git.ErrNotOnBranch == nil {
		t.Error("ErrNotOnBranch should not be nil")
	}
}
