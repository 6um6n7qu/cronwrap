package scheduler

import (
	"context"
	"testing"
	"time"
)

func TestAdd_ValidJob(t *testing.T) {
	s := New()
	err := s.Add(Job{
		Name:     "backup",
		Schedule: "@hourly",
		Command:  "/usr/bin/backup.sh",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if got := len(s.Jobs()); got != 1 {
		t.Fatalf("expected 1 job, got %d", got)
	}
}

func TestAdd_InvalidSchedule(t *testing.T) {
	s := New()
	err := s.Add(Job{
		Name:     "bad",
		Schedule: "not-a-cron",
		Command:  "echo hi",
	})
	if err == nil {
		t.Fatal("expected error for invalid schedule, got nil")
	}
}

func TestAdd_MissingName(t *testing.T) {
	s := New()
	err := s.Add(Job{Schedule: "@daily", Command: "echo"})
	if err == nil {
		t.Fatal("expected error for missing name, got nil")
	}
}

func TestAdd_MissingCommand(t *testing.T) {
	s := New()
	err := s.Add(Job{Name: "noop", Schedule: "@daily"})
	if err == nil {
		t.Fatal("expected error for missing command, got nil")
	}
}

func TestStart_CancelStopsScheduler(t *testing.T) {
	s := New()
	_ = s.Add(Job{
		Name:     "ping",
		Schedule: "@every 1s",
		Command:  "echo",
		Args:     []string{"ping"},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	s.Start(ctx)
	<-ctx.Done()
	// If we reach here without hanging, the scheduler stopped cleanly.
}

func TestJobs_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Add(Job{Name: "a", Schedule: "@daily", Command: "true"})

	jobs := s.Jobs()
	jobs[0].Name = "mutated"

	if s.Jobs()[0].Name != "a" {
		t.Fatal("Jobs() should return a copy, not a reference")
	}
}
