package snapshot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_AllReturnsEntries(t *testing.T) {
	s, _ := New(tempPath(t))
	_ = s.Record(Entry{JobName: "nightly", Success: true, ExitCode: 0, StartedAt: time.Now()})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/snapshots", nil)
	Handler(s).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var result map[string]Entry
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if _, ok := result["nightly"]; !ok {
		t.Error("expected 'nightly' in response")
	}
}

func TestHandler_SingleJobFound(t *testing.T) {
	s, _ := New(tempPath(t))
	_ = s.Record(Entry{JobName: "report", Success: true})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/snapshots/report", nil)
	Handler(s).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var e Entry
	if err := json.NewDecoder(rec.Body).Decode(&e); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if e.JobName != "report" {
		t.Errorf("expected job_name 'report', got %q", e.JobName)
	}
}

func TestHandler_SingleJobNotFound(t *testing.T) {
	s, _ := New(tempPath(t))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/snapshots/ghost", nil)
	Handler(s).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	s, _ := New(tempPath(t))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/snapshots", nil)
	Handler(s).ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
