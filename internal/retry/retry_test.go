package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/cronwrap/internal/retry"
)

var errFail = errors.New("job failed")

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: 0, Multiplier: 1.0}
	attempts, err := retry.Do(context.Background(), p, func(_ context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(attempts))
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	p := retry.Policy{MaxAttempts: 3, Delay: 0, Multiplier: 1.0}
	attempts, err := retry.Do(context.Background(), p, func(_ context.Context) error {
		calls++
		if calls < 3 {
			return errFail
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success on 3rd attempt, got %v", err)
	}
	if len(attempts) != 3 {
		t.Fatalf("expected 3 attempts, got %d", len(attempts))
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	p := retry.Policy{MaxAttempts: 2, Delay: 0, Multiplier: 1.0}
	attempts, err := retry.Do(context.Background(), p, func(_ context.Context) error {
		return errFail
	})
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if len(attempts) != 2 {
		t.Fatalf("expected 2 attempts, got %d", len(attempts))
	}
}

func TestDo_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := retry.Policy{MaxAttempts: 3, Delay: 0, Multiplier: 1.0}
	_, err := retry.Do(ctx, p, func(_ context.Context) error {
		return errFail
	})
	if err == nil {
		t.Fatal("expected error due to cancelled context")
	}
}

func TestDo_MultiplierIncreasesDelay(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: 10 * time.Millisecond, Multiplier: 2.0}
	start := time.Now()
	retry.Do(context.Background(), p, func(_ context.Context) error { //nolint:errcheck
		return errFail
	})
	// 10 ms + 20 ms = 30 ms minimum between attempts
	if time.Since(start) < 25*time.Millisecond {
		t.Error("expected multiplier to increase total delay")
	}
}

func TestDefaultPolicy_Values(t *testing.T) {
	p := retry.DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts 3, got %d", p.MaxAttempts)
	}
	if p.Delay != 5*time.Second {
		t.Errorf("expected Delay 5s, got %v", p.Delay)
	}
}
