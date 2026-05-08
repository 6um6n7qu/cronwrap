package jobhistory

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func buildHistory(t *testing.T) *History {
	t.Helper()
	h, err := New(10)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return h
}

func TestHandler_AllJobsReturns200(t *testing.T) {
	h := buildHistory(t)
	_ = h.Record("backup", true, time.Second, 0)
	_ = h.Record("cleanup", false, 2*time.Second, 1)

	req := httptest.NewRequest(http.MethodGet, "/jobhistory", nil)
	rec := httptest.NewRecorder()
	Handler(h)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestHandler_SingleJobFound(t *testing.T) {
	h := buildHistory(t)
	_ = h.Record("sync", true, 500*time.Millisecond, 0)

	req := httptest.NewRequest(http.MethodGet, "/jobhistory/sync", nil)
	rec := httptest.NewRecorder()
	Handler(h)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	if entries[0].JobName != "sync" {
		t.Errorf("expected job name 'sync', got %q", entries[0].JobName)
	}
}

func TestHandler_SingleJobNotFound(t *testing.T) {
	h := buildHistory(t)

	req := httptest.NewRequest(http.MethodGet, "/jobhistory/unknown", nil)
	rec := httptest.NewRecorder()
	Handler(h)(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := buildHistory(t)

	req := httptest.NewRequest(http.MethodPost, "/jobhistory", nil)
	rec := httptest.NewRecorder()
	Handler(h)(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandler_EmptyHistoryReturnsEmptyArray(t *testing.T) {
	h := buildHistory(t)

	req := httptest.NewRequest(http.MethodGet, "/jobhistory", nil)
	rec := httptest.NewRecorder()
	Handler(h)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty array, got %d entries", len(entries))
	}
}
