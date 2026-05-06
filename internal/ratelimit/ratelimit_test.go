package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/cronwrap/internal/ratelimit"
)

func TestNew_InvalidRate(t *testing.T) {
	_, err := ratelimit.New(0, 1)
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
}

func TestNew_InvalidBurst(t *testing.T) {
	_, err := ratelimit.New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero burst")
	}
}

func TestAllow_ConsumesTokens(t *testing.T) {
	l, err := ratelimit.New(time.Second, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 3; i++ {
		if !l.Allow("myjob") {
			t.Fatalf("expected Allow to return true on call %d", i+1)
		}
	}
	// Burst exhausted — next call should be denied.
	if l.Allow("myjob") {
		t.Fatal("expected Allow to return false after burst exhausted")
	}
}

func TestAllow_IndependentBuckets(t *testing.T) {
	l, _ := ratelimit.New(time.Second, 1)

	if !l.Allow("job-a") {
		t.Fatal("job-a should be allowed")
	}
	// job-b has its own bucket and should still be allowed.
	if !l.Allow("job-b") {
		t.Fatal("job-b should be allowed independently")
	}
	// job-a burst is now exhausted.
	if l.Allow("job-a") {
		t.Fatal("job-a should be denied after burst")
	}
}

func TestReset_RestoresTokens(t *testing.T) {
	l, _ := ratelimit.New(time.Second, 1)

	l.Allow("myjob") // consume the single token
	if l.Allow("myjob") {
		t.Fatal("expected denial before reset")
	}
	l.Reset("myjob")
	if !l.Allow("myjob") {
		t.Fatal("expected allow after reset")
	}
}

func TestWait_ReturnsOnContextCancel(t *testing.T) {
	l, _ := ratelimit.New(10*time.Second, 1)

	l.Allow("myjob") // exhaust burst

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := l.Wait(ctx, "myjob")
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestWait_SucceedsWhenTokenAvailable(t *testing.T) {
	l, _ := ratelimit.New(time.Millisecond, 2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := l.Wait(ctx, "fastjob"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
