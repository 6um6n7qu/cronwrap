package jobcooldown

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRegister_AddsJob(t *testing.T) {
	s := New(nil)
	if err := s.Register("backup", 5*time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	s := New(nil)
	if err := s.Register("", time.Minute); err != ErrEmptyName {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}
}

func TestAllow_NewJobPermitted(t *testing.T) {
	s := New(nil)
	_ = s.Register("sync", time.Minute)
	if err := s.Allow("sync"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_UnknownJobReturnsError(t *testing.T) {
	s := New(nil)
	if err := s.Allow("ghost"); err != ErrUnknownJob {
		t.Fatalf("expected ErrUnknownJob, got %v", err)
	}
}

func TestAllow_CooldownActiveAfterRecord(t *testing.T) {
	now := time.Now()
	s := New(fixedNow(now))
	_ = s.Register("export", 10*time.Minute)
	_ = s.Record("export")
	if err := s.Allow("export"); err != ErrCooldownActive {
		t.Fatalf("expected ErrCooldownActive, got %v", err)
	}
}

func TestAllow_PermittedAfterCooldownExpires(t *testing.T) {
	now := time.Now()
	s := New(fixedNow(now))
	_ = s.Register("export", 10*time.Minute)
	_ = s.Record("export")

	// Advance time past the cooldown.
	s.now = fixedNow(now.Add(11 * time.Minute))
	if err := s.Allow("export"); err != nil {
		t.Fatalf("expected nil after cooldown, got %v", err)
	}
}

func TestRecord_UnknownJobReturnsError(t *testing.T) {
	s := New(nil)
	if err := s.Record("ghost"); err != ErrUnknownJob {
		t.Fatalf("expected ErrUnknownJob, got %v", err)
	}
}

func TestRemaining_ZeroForFreshJob(t *testing.T) {
	s := New(nil)
	_ = s.Register("cleanup", 5*time.Minute)
	rem, err := s.Remaining("cleanup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rem != 0 {
		t.Fatalf("expected 0 remaining, got %v", rem)
	}
}

func TestRemaining_ReturnsPositiveDuringCooldown(t *testing.T) {
	now := time.Now()
	s := New(fixedNow(now))
	_ = s.Register("report", 10*time.Minute)
	_ = s.Record("report")

	s.now = fixedNow(now.Add(3 * time.Minute))
	rem, err := s.Remaining("report")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 7 * time.Minute
	if rem != expected {
		t.Fatalf("expected %v remaining, got %v", expected, rem)
	}
}

func TestRemaining_ZeroAfterCooldownExpires(t *testing.T) {
	now := time.Now()
	s := New(fixedNow(now))
	_ = s.Register("report", 5*time.Minute)
	_ = s.Record("report")

	s.now = fixedNow(now.Add(10 * time.Minute))
	rem, err := s.Remaining("report")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rem != 0 {
		t.Fatalf("expected 0 after expiry, got %v", rem)
	}
}
