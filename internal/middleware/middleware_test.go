package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/cronwrap/internal/middleware"
)

func newTestLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, nil))
}

func TestRequestLogger_LogsRequestFields(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	handler := middleware.RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	log := buf.String()
	if !strings.Contains(log, "http request") {
		t.Errorf("expected log to contain 'http request', got: %s", log)
	}
	if !strings.Contains(log, "/health") {
		t.Errorf("expected log to contain path '/health', got: %s", log)
	}
	if !strings.Contains(log, "200") {
		t.Errorf("expected log to contain status 200, got: %s", log)
	}
}

func TestRequestLogger_CapturesNon200Status(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	handler := middleware.RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !strings.Contains(buf.String(), "404") {
		t.Errorf("expected log to contain status 404, got: %s", buf.String())
	}
}

func TestRecovery_CatchesPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	handler := middleware.Recovery(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	}))

	req := httptest.NewRequest(http.MethodGet, "/crash", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
	if !strings.Contains(buf.String(), "panic recovered") {
		t.Errorf("expected panic log, got: %s", buf.String())
	}
}

func TestChain_AppliesMiddlewareInOrder(t *testing.T) {
	var order []string

	mk := func(label string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, label)
				next.ServeHTTP(w, r)
			})
		}
	}

	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler := middleware.Chain(base, mk("first"), mk("second"))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if len(order) != 2 || order[0] != "first" || order[1] != "second" {
		t.Errorf("unexpected middleware order: %v", order)
	}
}
