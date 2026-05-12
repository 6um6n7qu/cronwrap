package jobnotify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func buildNotifyStore(t *testing.T) *Store {
	t.Helper()
	s := New()
	s.Record("backup", EventFailed, "email")
	s.Record("backup", EventFinished, "slack")
	s.Record("report", EventStarted, "webhook")
	return s
}

func TestHandler_AllHistoryReturns200(t *testing.T) {
	s := buildNotifyStore(t)
	req := httptest.NewRequest(http.MethodGet, "/jobnotify", nil)
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestHandler_SingleJobFound(t *testing.T) {
	s := buildNotifyStore(t)
	req := httptest.NewRequest(http.MethodGet, "/jobnotify/backup", nil)
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for backup, got %d", len(entries))
	}
	for _, e := range entries {
		if e.JobName != "backup" {
			t.Errorf("unexpected job name: %s", e.JobName)
		}
	}
}

func TestHandler_SingleJobNotFound(t *testing.T) {
	s := buildNotifyStore(t)
	req := httptest.NewRequest(http.MethodGet, "/jobnotify/ghost", nil)
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	s := New()
	req := httptest.NewRequest(http.MethodPost, "/jobnotify", nil)
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
