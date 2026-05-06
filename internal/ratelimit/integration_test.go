package ratelimit_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cronwrap/internal/ratelimit"
)

// TestConcurrent_SafeUnderContention verifies that the limiter is
// goroutine-safe and that no more than maxBurst tokens are consumed in
// the initial burst window across concurrent callers.
func TestConcurrent_SafeUnderContention(t *testing.T) {
	const burst = 5
	l, err := ratelimit.New(time.Minute, burst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var allowed atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.Allow("concurrent-job") {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := allowed.Load(); got > burst {
		t.Errorf("allowed %d executions, want at most %d", got, burst)
	}
}

// TestRefill_TokensReplenishOverTime checks that tokens are refilled after
// the configured rate window elapses.
func TestRefill_TokensReplenishOverTime(t *testing.T) {
	l, err := ratelimit.New(50*time.Millisecond, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !l.Allow("refill-job") {
		t.Fatal("first call should be allowed")
	}
	if l.Allow("refill-job") {
		t.Fatal("second immediate call should be denied")
	}

	// Wait for one full rate window so a token is refilled.
	time.Sleep(60 * time.Millisecond)

	if !l.Allow("refill-job") {
		t.Fatal("call after rate window should be allowed")
	}
}
