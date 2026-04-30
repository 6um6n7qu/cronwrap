package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// ServerConfig holds configuration for the health-check HTTP server.
type ServerConfig struct {
	// Addr is the TCP address to listen on, e.g. ":8080".
	Addr string
	// ReadTimeout is the maximum duration for reading the request.
	ReadTimeout time.Duration
	// WriteTimeout is the maximum duration before timing out writes.
	WriteTimeout time.Duration
}

// DefaultServerConfig returns a ServerConfig with sensible defaults.
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}

// Server wraps an http.Server and a Checker to serve health status.
type Server struct {
	httpServer *http.Server
	checker    *Checker
}

// NewServer creates a Server using the provided Checker and config.
func NewServer(checker *Checker, cfg ServerConfig) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", checker.Handler())

	return &Server{
		checker: checker,
		httpServer: &http.Server{
			Addr:         cfg.Addr,
			Handler:      mux,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

// Start begins listening and serving. It blocks until the server stops.
// Returns nil if the server was shut down gracefully via the context.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	}
}
