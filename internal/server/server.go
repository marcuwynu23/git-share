package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/marcuwynu23/git-share/internal/discovery"
	"github.com/marcuwynu23/git-share/internal/git"
	"github.com/marcuwynu23/git-share/internal/middleware"
)

type ServerConfig struct {
	Port     int
	ReadOnly bool
	Hostname string
	Timeout  time.Duration
}

type Server struct {
	config  *ServerConfig
	http    *http.Server
	handler *Handler
	started time.Time
}

func New(repo *git.RepoInfo, config *ServerConfig) *Server {
	handler := NewHandler(repo, config)
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.Root)
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/info", handler.Info)
	mux.HandleFunc("/git/", handler.Git)

	var h http.Handler = mux
	h = middleware.Chain(h, middleware.Recoverer, middleware.Logger, middleware.CORS)

	addr := fmt.Sprintf(":%d", config.Port)
	if config.Hostname != "" {
		addr = net.JoinHostPort(config.Hostname, fmt.Sprintf("%d", config.Port))
	}

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
	}

	return &Server{
		config:  config,
		http:    httpServer,
		handler: handler,
		started: time.Now(),
	}
}

func (s *Server) Start(ctx context.Context) error {
	addrs := discovery.Discover()
	repoName := filepath.Base(s.handler.Repo.Root)

	lis, err := net.Listen("tcp", s.http.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	actualPort := lis.Addr().(*net.TCPAddr).Port

	fmt.Printf("Repository: %s\n\n", repoName)
	fmt.Printf("Listening:\n")
	fmt.Printf("  http://localhost:%d\n", actualPort)
	if addrs.LAN != "" {
		fmt.Printf("\nLAN:\n")
		fmt.Printf("  http://%s:%d\n", addrs.LAN, actualPort)
	}
	cloneAddr := addrs.LAN
	if cloneAddr == "" {
		cloneAddr = "localhost"
	}
	fmt.Printf("\nClone:\n\n")
	fmt.Printf("  git clone http://%s:%d/%s.git\n\n", cloneAddr, actualPort, repoName)
	fmt.Printf("Press Ctrl+C to stop.\n")

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.http.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.http.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func (s *Server) Uptime() time.Duration {
	return time.Since(s.started)
}
