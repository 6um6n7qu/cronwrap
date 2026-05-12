package jobstatus_test

import (
	"testing"

	"github.com/cronwrap/internal/jobstatus"
)

func TestRegister_AddsJob(t *testing.T) {
	s := jobstatus.New()
	if err := s.Register("backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Status != jobstatus.StatusIdle {
		t.Errorf("expected idle, got %s", e.Status)
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	s := jobstatus.New()
	if err := s.Register(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_DuplicateReturnsError(t *testing.T) {
	s := jobstatus.New()
	_ = s.Register("job1")
	if err := s.Register("job1"); err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestSetRunning_UpdatesStatus(t *testing.T) {
	s := jobstatus.New()
	_ = s.Register("job1")
	if err := s.SetRunning("job1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := s.Get("job1")
	if e.Status != jobstatus.StatusRunning {
		t.Errorf("expected running, got %s", e.Status)
	}
	if e.StartedAt.IsZero() {
		t.Error("expected StartedAt to be set")
	}
}

func TestSetSuccess_UpdatesStatus(t *testing.T) {
	s := jobstatus.New()
	_ = s.Register("job1")
	_ = s.SetRunning("job1")
	if err := s.SetSuccess("job1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := s.Get("job1")
	if e.Status != jobstatus.StatusSuccess {
		t.Errorf("expected success, got %s", e.Status)
	}
}

func TestSetFailed_SetsExitCode(t *testing.T) {
	s := jobstatus.New()
	_ = s.Register("job1")
	if err := s.SetFailed("job1", 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := s.Get("job1")
	if e.Status != jobstatus.StatusFailed {
		t.Errorf("expected failed, got %s", e.Status)
	}
	if e.ExitCode != 2 {
		t.Errorf("expected exit code 2, got %d", e.ExitCode)
	}
}

func TestSetRunning_UnknownJobReturnsError(t *testing.T) {
	s := jobstatus.New()
	if err := s.SetRunning("ghost"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s := jobstatus.New()
	_ = s.Register("a")
	_ = s.Register("b")
	_ = s.Register("c")
	if len(s.All()) != 3 {
		t.Errorf("expected 3 entries, got %d", len(s.All()))
	}
}

func TestGet_UnknownJobReturnsFalse(t *testing.T) {
	s := jobstatus.New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}
