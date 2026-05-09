package jobpause_test

import (
	"testing"

	"github.com/cronwrap/internal/jobpause"
)

func TestRegister_AddsJob(t *testing.T) {
	s := jobpause.New()
	if err := s.Register("backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsPaused("backup") {
		t.Error("newly registered job should not be paused")
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	s := jobpause.New()
	if err := s.Register(""); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestPause_MarksJobPaused(t *testing.T) {
	s := jobpause.New()
	_ = s.Register("cleanup")
	if err := s.Pause("cleanup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.IsPaused("cleanup") {
		t.Error("job should be paused")
	}
}

func TestPause_UnknownJobReturnsError(t *testing.T) {
	s := jobpause.New()
	if err := s.Pause("ghost"); err == nil {
		t.Error("expected error for unknown job")
	}
}

func TestResume_ClearsPauseState(t *testing.T) {
	s := jobpause.New()
	_ = s.Register("report")
	_ = s.Pause("report")
	if err := s.Resume("report"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsPaused("report") {
		t.Error("job should no longer be paused after resume")
	}
}

func TestResume_UnknownJobReturnsError(t *testing.T) {
	s := jobpause.New()
	if err := s.Resume("ghost"); err == nil {
		t.Error("expected error for unknown job")
	}
}

func TestIsPaused_UnknownJobReturnsFalse(t *testing.T) {
	s := jobpause.New()
	if s.IsPaused("nonexistent") {
		t.Error("unknown job should not be considered paused")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := jobpause.New()
	_ = s.Register("job-a")
	_ = s.Register("job-b")
	_ = s.Pause("job-a")

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if !all["job-a"].Paused {
		t.Error("job-a should be paused in snapshot")
	}
	if all["job-b"].Paused {
		t.Error("job-b should not be paused in snapshot")
	}
}

func TestPause_RecordsPausedAt(t *testing.T) {
	s := jobpause.New()
	_ = s.Register("timed")
	_ = s.Pause("timed")
	all := s.All()
	if all["timed"].PausedAt.IsZero() {
		t.Error("PausedAt should be set after pausing")
	}
}
