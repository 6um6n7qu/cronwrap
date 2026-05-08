package jobhistory

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero capacity")
	}
}

func TestRecord_StoresEntry(t *testing.T) {
	s, _ := New(10)
	err := s.Record(Entry{JobName: "backup", Success: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	all := s.All()
	if len(all) != 1 || all[0].JobName != "backup" {
		t.Fatalf("expected 1 entry, got %v", all)
	}
}

func TestRecord_EmptyNameReturnsError(t *testing.T) {
	s, _ := New(10)
	if err := s.Record(Entry{}); err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestRecord_RingBufferEvictsOldest(t *testing.T) {
	s, _ := New(3)
	for i, name := range []string{"a", "b", "c", "d"} {
		_ = i
		s.Record(Entry{JobName: name, Success: true})
	}
	all := s.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
	if all[0].JobName != "b" {
		t.Errorf("expected oldest entry to be 'b', got %q", all[0].JobName)
	}
	if all[2].JobName != "d" {
		t.Errorf("expected newest entry to be 'd', got %q", all[2].JobName)
	}
}

func TestRecord_SetsDuration(t *testing.T) {
	s, _ := New(5)
	now := time.Now()
	s.Record(Entry{
		JobName:    "sync",
		StartedAt:  now.Add(-2 * time.Second),
		FinishedAt: now,
	})
	all := s.All()
	if all[0].Duration < 2*time.Second {
		t.Errorf("expected duration >= 2s, got %v", all[0].Duration)
	}
}

func TestForJob_FiltersEntries(t *testing.T) {
	s, _ := New(10)
	s.Record(Entry{JobName: "alpha", Success: true})
	s.Record(Entry{JobName: "beta", Success: false})
	s.Record(Entry{JobName: "alpha", Success: false})

	result := s.ForJob("alpha")
	if len(result) != 2 {
		t.Fatalf("expected 2 entries for 'alpha', got %d", len(result))
	}
}

func TestHandler_AllEntriesReturns200(t *testing.T) {
	s, _ := New(10)
	s.Record(Entry{JobName: "job1", Success: true})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jobhistory", nil)
	Handler(s)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestHandler_SingleJobFound(t *testing.T) {
	s, _ := New(10)
	s.Record(Entry{JobName: "cleanup", Success: true})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jobhistory/cleanup", nil)
	Handler(s)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestHandler_SingleJobNotFound(t *testing.T) {
	s, _ := New(10)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jobhistory/missing", nil)
	Handler(s)(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	s, _ := New(10)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/jobhistory", nil)
	Handler(s)(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
