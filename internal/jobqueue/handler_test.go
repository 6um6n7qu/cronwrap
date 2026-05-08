package jobqueue

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_AllEntriesReturns200(t *testing.T) {
	q, _ := New(8)
	_ = q.Enqueue("backup", 1)
	_ = q.Enqueue("report", 2)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/queue", nil)
	Handler(q)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestHandler_SingleEntryFound(t *testing.T) {
	q, _ := New(8)
	_ = q.Enqueue("cleanup", 3)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/queue/cleanup", nil)
	Handler(q)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entry Entry
	if err := json.NewDecoder(rec.Body).Decode(&entry); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if entry.Name != "cleanup" {
		t.Errorf("expected name 'cleanup', got %q", entry.Name)
	}
}

func TestHandler_SingleEntryNotFound(t *testing.T) {
	q, _ := New(8)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/queue/missing", nil)
	Handler(q)(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	q, _ := New(8)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/queue", nil)
	Handler(q)(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandler_EmptyQueueReturnsEmptyArray(t *testing.T) {
	q, _ := New(4)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/queue", nil)
	Handler(q)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty array, got %d entries", len(entries))
	}
}
