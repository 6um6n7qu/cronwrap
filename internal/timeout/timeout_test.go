package timeout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/cronwrap/internal/timeout"
)

func TestDefaultPolicy_HasReasonableDefaults(t *testing.T) {
	p := timeout.DefaultPolicy()
	if p.Duration != 30*time.Minute {
		t.Errorf("expected 30m duration, got %s", p.Duration)
	}
	if p.GracePeriod != 5*time.Second {
		t.Errorf("expected 5s grace period, got %s", p.GracePeriod)
	}
}

func TestWrap_CompletesWithinDeadline(t *testing.T) {
	p := timeout.Policy{Duration: 500 * time.Millisecond, GracePeriod: 50 * time.Millisecond}
	err := p.Wrap(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWrap_ReturnsJobError(t *testing.T) {
	p := timeout.Policy{Duration: 500 * time.Millisecond}
	sentinel := errors.New("job failed")
	err := p.Wrap(context.Background(), func(ctx context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestWrap_ExceedsDeadline(t *testing.T) {
	p := timeout.Policy{Duration: 50 * time.Millisecond, GracePeriod: 10 * time.Millisecond}
	err := p.Wrap(context.Background(), func(ctx context.Context) error {
		time.Sleep(500 * time.Millisecond)
		return nil
	})
	if !errors.Is(err, timeout.ErrDeadlineExceeded) {
		t.Fatalf("expected ErrDeadlineExceeded, got %v", err)
	}
}

func TestWrap_ZeroDurationNoTimeout(t *testing.T) {
	p := timeout.Policy{Duration: 0}
	err := p.Wrap(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error with zero duration, got %v", err)
	}
}

func TestWrap_ContextCancelledExternally(t *testing.T) {
	p := timeout.Policy{Duration: 10 * time.Second}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	err := p.Wrap(ctx, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestWithContext_ReturnsCancellableContext(t *testing.T) {
	p := timeout.Policy{Duration: 1 * time.Second}
	ctx, cancel := p.WithContext(context.Background())
	defer cancel()
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected context to have a deadline")
	}
	if time.Until(deadline) > 1*time.Second {
		t.Error("deadline is further than expected")
	}
}
