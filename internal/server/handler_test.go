package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/markwayne/git-share/internal/git"
)

func testHandler(t *testing.T) *Handler {
	t.Helper()
	repo := &git.RepoInfo{
		Root:   "/home/user/projects/myrepo",
		GitDir: "/home/user/projects/myrepo/.git",
		Branch: "main",
		IsBare: false,
	}
	config := &ServerConfig{
		Port:     8080,
		ReadOnly: true,
	}
	return NewHandler(repo, config)
}

func TestHealth(t *testing.T) {
	h := testHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	h.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var resp HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("Status = %q, want ok", resp.Status)
	}
}

func TestInfo(t *testing.T) {
	h := testHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/info", nil)
	h.Info(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}

	var resp InfoResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Repository != "/home/user/projects/myrepo" {
		t.Errorf("Repository = %q, want /home/user/projects/myrepo", resp.Repository)
	}
	if resp.Branch != "main" {
		t.Errorf("Branch = %q, want main", resp.Branch)
	}
	if resp.Bare {
		t.Error("Bare should be false")
	}
	if resp.Port != 8080 {
		t.Errorf("Port = %d, want 8080", resp.Port)
	}
	if !resp.ReadOnly {
		t.Error("ReadOnly should be true")
	}
	if resp.CloneURL != "http://localhost:8080/" {
		t.Errorf("CloneURL = %q, want http://localhost:8080/", resp.CloneURL)
	}
}

func TestIsGitRequest(t *testing.T) {
	tests := []struct {
		name      string
		userAgent string
		query     string
		path      string
		want      bool
	}{
		{"git user-agent", "git/2.40.0", "", "/", true},
		{"browser user-agent", "Mozilla/5.0", "", "/", false},
		{"service in query", "", "service=git-upload-pack", "/", true},
		{"info/refs path", "", "", "/repo.git/info/refs", true},
		{"HEAD path", "", "", "/repo.git/HEAD", true},
		{"git-upload-pack path", "", "", "/repo.git/git-upload-pack", true},
		{"git-receive-pack path", "", "", "/repo.git/git-receive-pack", true},
		{".git in path", "", "", "/repo.git/objects/abc123", true},
		{"normal browser path", "", "", "/", false},
		{"random path", "", "", "/some-page", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			req.Header.Set("User-Agent", tt.userAgent)
			req.URL.RawQuery = tt.query

			got := isGitRequest(req)
			if got != tt.want {
				t.Errorf("isGitRequest = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitReadOnlyBlocksPush(t *testing.T) {
	h := testHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/repo.git/git-receive-pack", nil)
	req.Header.Set("Content-Type", "application/x-git-receive-pack-request")
	h.Git(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403", rec.Code)
	}
}

func TestGitReadWriteAllowsPush(t *testing.T) {
	repo := &git.RepoInfo{
		Root:   "/home/user/projects/myrepo",
		GitDir: "/home/user/projects/myrepo/.git",
		Branch: "main",
		IsBare: false,
	}
	config := &ServerConfig{
		Port:     8080,
		ReadOnly: false,
	}
	h := NewHandler(repo, config)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/repo.git/git-receive-pack", nil)
	req.Header.Set("Content-Type", "application/x-git-receive-pack-request")

	h.Git(rec, req)

	if rec.Code == http.StatusForbidden {
		t.Error("push should not be blocked in read-write mode")
	}
}

func TestRootRoutesGitRequests(t *testing.T) {
	h := testHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/repo.git/info/refs", nil)
	req.Header.Set("User-Agent", "git/2.40.0")

	h.Root(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("Root should route git requests, got 404")
	}
}

func TestRootServesDashboard(t *testing.T) {
	h := testHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	h.Root(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html; charset=utf-8", ct)
	}
}
