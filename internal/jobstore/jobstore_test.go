package jobstore

import (
	"errors"
	"testing"
	"time"
)

func TestStart_RecordsRunningEntry(t *testing.T) {
	s := New()
	if err := s.Start("backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Status != StatusRunning {
		t.Errorf("expected running, got %s", e.Status)
	}
	if e.StartedAt.IsZero() {
		t.Error("expected StartedAt to be set")
	}
}

func TestStart_EmptyNameReturnsError(t *testing.T) {
	s := New()
	if err := s.Start(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestFinish_SuccessEntry(t *testing.T) {
	s := New()
	_ = s.Start("report")
	if err := s.Finish("report", 0, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := s.Get("report")
	if e.Status != StatusSuccess {
		t.Errorf("expected success, got %s", e.Status)
	}
	if e.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", e.ExitCode)
	}
	if e.FinishedAt.IsZero() {
		t.Error("expected FinishedAt to be set")
	}
}

func TestFinish_FailureEntry(t *testing.T) {
	s := New()
	_ = s.Start("sync")
	if err := s.Finish("sync", 1, errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := s.Get("sync")
	if e.Status != StatusFailed {
		t.Errorf("expected failed, got %s", e.Status)
	}
	if e.Error != "disk full" {
		t.Errorf("unexpected error string: %s", e.Error)
	}
}

func TestGet_UnknownJobReturnsFalse(t *testing.T) {
	s := New()
	_, ok := s.Get("ghost")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New()
	_ = s.Start("a")
	_ = s.Start("b")
	_ = s.Finish("b", 0, nil)
	entries := s.All()
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestFinish_SetsFinishedAtAfterStart(t *testing.T) {
	s := New()
	_ = s.Start("cleanup")
	time.Sleep(time.Millisecond)
	_ = s.Finish("cleanup", 0, nil)
	e, _ := s.Get("cleanup")
	if !e.FinishedAt.After(e.StartedAt) {
		t.Error("expected FinishedAt to be after StartedAt")
	}
}
