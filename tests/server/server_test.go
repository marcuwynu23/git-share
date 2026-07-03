package server_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marcuwynu23/git-share/internal/git"
	"github.com/marcuwynu23/git-share/internal/server"
)

func testHandler(t *testing.T) *server.Handler {
	t.Helper()
	repo := &git.RepoInfo{
		Root:   "/home/user/projects/myrepo",
		GitDir: "/home/user/projects/myrepo/.git",
		Branch: "main",
		IsBare: false,
	}
	config := &server.ServerConfig{
		Port:     9720,
		ReadOnly: true,
	}
	return server.NewHandler(repo, config)
}

func TestNewHandler(t *testing.T) {
	h := testHandler(t)
	if h == nil {
		t.Fatal("NewHandler() returned nil")
	}
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

	var resp server.HealthResponse
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

	var resp server.InfoResponse
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
	if resp.Port != 9720 {
		t.Errorf("Port = %d, want 9720", resp.Port)
	}
	if !resp.ReadOnly {
		t.Error("ReadOnly should be true")
	}
	if resp.CloneURL != "http://localhost:9720/" {
		t.Errorf("CloneURL = %q, want http://localhost:9720/", resp.CloneURL)
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
	config := &server.ServerConfig{
		Port:     9720,
		ReadOnly: false,
	}
	h := server.NewHandler(repo, config)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/repo.git/git-receive-pack", nil)
	req.Header.Set("Content-Type", "application/x-git-receive-pack-request")

	h.Git(rec, req)

	if rec.Code == http.StatusForbidden {
		t.Error("push should not be blocked in read-write mode")
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

func TestNewServer(t *testing.T) {
	repo := &git.RepoInfo{
		Root:   "/home/user/projects/myrepo",
		GitDir: "/home/user/projects/myrepo/.git",
		Branch: "main",
		IsBare: false,
	}
	config := &server.ServerConfig{
		Port:     0,
		ReadOnly: false,
	}
	srv := server.New(repo, config)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNewServerWithHostname(t *testing.T) {
	repo := &git.RepoInfo{
		Root:   "/home/user/projects/myrepo",
		GitDir: "/home/user/projects/myrepo/.git",
		Branch: "main",
		IsBare: false,
	}
	config := &server.ServerConfig{
		Port:     9720,
		ReadOnly: true,
		Hostname: "192.168.1.1",
	}
	srv := server.New(repo, config)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestServerUptime(t *testing.T) {
	repo := &git.RepoInfo{
		Root:   "/home/user/projects/myrepo",
		GitDir: "/home/user/projects/myrepo/.git",
		Branch: "main",
		IsBare: false,
	}
	config := &server.ServerConfig{
		Port:     0,
		ReadOnly: false,
	}
	srv := server.New(repo, config)
	uptime := srv.Uptime()
	if uptime < 0 {
		t.Error("Uptime should be non-negative")
	}
}

func TestServerStop(t *testing.T) {
	repo := &git.RepoInfo{
		Root:   "/home/user/projects/myrepo",
		GitDir: "/home/user/projects/myrepo/.git",
		Branch: "main",
		IsBare: false,
	}
	config := &server.ServerConfig{
		Port:     0,
		ReadOnly: false,
	}
	srv := server.New(repo, config)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start(context.Background())
	}()

	time.Sleep(10 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil && err != http.ErrServerClosed {
		t.Fatalf("Stop() error = %v", err)
	}
}
