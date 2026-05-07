package throttle_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/throttle"
)

func TestNew_InvalidLimit(t *testing.T) {
	_, err := throttle.New(0)
	if err == nil {
		t.Fatal("expected error for zero limit")
	}
	_, err = throttle.New(-1)
	if err == nil {
		t.Fatal("expected error for negative limit")
	}
}

func TestNew_ValidLimit(t *testing.T) {
	th, err := throttle.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.Limit() != 3 {
		t.Errorf("expected limit 3, got %d", th.Limit())
	}
}

func TestTryAcquire_WithinLimit(t *testing.T) {
	th, _ := throttle.New(2)
	if err := th.TryAcquire(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.Active() != 1 {
		t.Errorf("expected 1 active, got %d", th.Active())
	}
}

func TestTryAcquire_ExceedsLimit(t *testing.T) {
	th, _ := throttle.New(1)
	_ = th.TryAcquire()
	err := th.TryAcquire()
	if err != throttle.ErrThrottled {
		t.Errorf("expected ErrThrottled, got %v", err)
	}
}

func TestRelease_DecrementsActive(t *testing.T) {
	th, _ := throttle.New(2)
	_ = th.TryAcquire()
	_ = th.TryAcquire()
	th.Release()
	if th.Active() != 1 {
		t.Errorf("expected 1 active after release, got %d", th.Active())
	}
}

func TestAcquire_BlocksUntilSlotAvailable(t *testing.T) {
	th, _ := throttle.New(1)
	_ = th.TryAcquire()

	go func() {
		time.Sleep(20 * time.Millisecond)
		th.Release()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	if err := th.Acquire(ctx); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAcquire_ContextCancellation(t *testing.T) {
	th, _ := throttle.New(1)
	_ = th.TryAcquire()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := th.Acquire(ctx)
	if err != throttle.ErrThrottled {
		t.Errorf("expected ErrThrottled on context cancel, got %v", err)
	}
}

func TestConcurrent_SafeUnderContention(t *testing.T) {
	th, _ := throttle.New(5)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()
			if err := th.Acquire(ctx); err == nil {
				time.Sleep(5 * time.Millisecond)
				th.Release()
			}
		}()
	}
	wg.Wait()
	if th.Active() != 0 {
		t.Errorf("expected 0 active after all goroutines, got %d", th.Active())
	}
}
