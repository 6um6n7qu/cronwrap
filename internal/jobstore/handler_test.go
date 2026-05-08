package jobstore

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_AllJobsReturns200(t *testing.T) {
	s := New()
	_ = s.Start("alpha")
	_ = s.Finish("alpha", 0, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	Handler(s).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestHandler_SingleJobFound(t *testing.T) {
	s := New()
	_ = s.Start("beta")
	_ = s.Finish("beta", 1, errors.New("oops"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jobs/beta", nil)
	Handler(s).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var e Entry
	if err := json.NewDecoder(rec.Body).Decode(&e); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if e.Name != "beta" {
		t.Errorf("expected name beta, got %s", e.Name)
	}
	if e.Status != StatusFailed {
		t.Errorf("expected failed, got %s", e.Status)
	}
}

func TestHandler_SingleJobNotFound(t *testing.T) {
	s := New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jobs/ghost", nil)
	Handler(s).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	s := New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/jobs", nil)
	Handler(s).ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
