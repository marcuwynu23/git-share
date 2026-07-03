package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/markwayne/git-share/internal/discovery"
	"github.com/markwayne/git-share/internal/git"
	"github.com/markwayne/git-share/internal/ui"
)

type Handler struct {
	Repo   *git.RepoInfo
	Config *ServerConfig
}

type HealthResponse struct {
	Status string `json:"status"`
}

type InfoResponse struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Bare       bool   `json:"bare"`
	Port       int    `json:"port"`
	ReadOnly   bool   `json:"readonly"`
	CloneURL   string `json:"clone_url"`
	LANURL     string `json:"lan_url"`
}

func NewHandler(repo *git.RepoInfo, config *ServerConfig) *Handler {
	return &Handler{
		Repo:   repo,
		Config: config,
	}
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	if isGitRequest(r) {
		h.Git(w, r)
		return
	}
	ui.Dashboard(w, r)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	addrs := discovery.Discover()
	portStr := fmt.Sprintf("%d", h.Config.Port)
	cloneURL := "http://localhost:" + portStr + "/"
	lanURL := ""
	if addrs.LAN != "" {
		lanURL = "http://" + addrs.LAN + ":" + portStr + "/"
	}

	resp := InfoResponse{
		Repository: h.Repo.Root,
		Branch:     h.Repo.Branch,
		Bare:       h.Repo.IsBare,
		Port:       h.Config.Port,
		ReadOnly:   h.Config.ReadOnly,
		CloneURL:   cloneURL,
		LANURL:     lanURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Git(w http.ResponseWriter, r *http.Request) {
	if h.Config.ReadOnly {
		if r.Method == "POST" {
			bodyType := r.Header.Get("Content-Type")
			if strings.Contains(bodyType, "git-receive-pack") {
				http.Error(w, "push is disabled (read-only mode)", http.StatusForbidden)
				return
			}
		}
	}

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/git")
	if path == "" {
		path = "/"
	}

	r.URL.Path = path

	repoName := filepath.Base(h.Repo.Root)

	git.HandleSmartHTTP(h.Repo, w, r, repoName, h.Config.ReadOnly)
}

func isGitRequest(r *http.Request) bool {
	ua := r.UserAgent()
	if strings.Contains(ua, "git/") {
		return true
	}

	if strings.Contains(r.URL.RawQuery, "service=git-") {
		return true
	}

	path := r.URL.Path
	if strings.HasSuffix(path, "/info/refs") ||
		strings.HasSuffix(path, "/HEAD") ||
		strings.HasSuffix(path, "/git-upload-pack") ||
		strings.HasSuffix(path, "/git-receive-pack") {
		return true
	}

	if strings.Contains(path, ".git") {
		return true
	}

	return false
}
