package jobruncount_test

import (
	"testing"

	"github.com/cronwrap/internal/jobruncount"
)

func TestRegister_AddsJob(t *testing.T) {
	s := jobruncount.New()
	if err := s.Register("backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected job to be registered")
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	s := jobruncount.New()
	if err := s.Register(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_DuplicateReturnsError(t *testing.T) {
	s := jobruncount.New()
	_ = s.Register("sync")
	if err := s.Register("sync"); err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestRecordSuccess_IncrementsCounters(t *testing.T) {
	s := jobruncount.New()
	_ = s.Register("etl")
	_ = s.RecordSuccess("etl")
	_ = s.RecordSuccess("etl")
	e, _ := s.Get("etl")
	if e.TotalRuns != 2 {
		t.Errorf("expected TotalRuns=2, got %d", e.TotalRuns)
	}
	if e.SuccessRuns != 2 {
		t.Errorf("expected SuccessRuns=2, got %d", e.SuccessRuns)
	}
	if e.FailureRuns != 0 {
		t.Errorf("expected FailureRuns=0, got %d", e.FailureRuns)
	}
}

func TestRecordFailure_IncrementsCounters(t *testing.T) {
	s := jobruncount.New()
	_ = s.Register("report")
	_ = s.RecordFailure("report")
	e, _ := s.Get("report")
	if e.TotalRuns != 1 {
		t.Errorf("expected TotalRuns=1, got %d", e.TotalRuns)
	}
	if e.FailureRuns != 1 {
		t.Errorf("expected FailureRuns=1, got %d", e.FailureRuns)
	}
}

func TestRecordSuccess_UnknownJobReturnsError(t *testing.T) {
	s := jobruncount.New()
	if err := s.RecordSuccess("ghost"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestRecordFailure_UnknownJobReturnsError(t *testing.T) {
	s := jobruncount.New()
	if err := s.RecordFailure("ghost"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGet_UnknownJobReturnsFalse(t *testing.T) {
	s := jobruncount.New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected ok=false for unknown job")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s := jobruncount.New()
	_ = s.Register("a")
	_ = s.Register("b")
	_ = s.RecordSuccess("a")
	entries := s.All()
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}
