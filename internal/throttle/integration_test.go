package throttle_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/throttle"
)

// TestMaxConcurrency_NeverExceedsLimit verifies that the number of
// goroutines holding a slot never exceeds the configured limit.
func TestMaxConcurrency_NeverExceedsLimit(t *testing.T) {
	const limit = 3
	const workers = 30

	th, _ := throttle.New(limit)

	var (
		wg      sync.WaitGroup
		current int64
		peak    int64
		mu      sync.Mutex
	)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			if err := th.Acquire(ctx); err != nil {
				return
			}
			defer th.Release()

			now := atomic.AddInt64(&current, 1)
			mu.Lock()
			if now > peak {
				peak = now
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)
			atomic.AddInt64(&current, -1)
		}()
	}

	wg.Wait()

	if peak > int64(limit) {
		t.Errorf("peak concurrency %d exceeded limit %d", peak, limit)
	}
}

// TestRelease_AllowsQueuedAcquire ensures that releasing a slot unblocks
// a waiting Acquire call.
func TestRelease_AllowsQueuedAcquire(t *testing.T) {
	th, _ := throttle.New(1)
	_ = th.TryAcquire() // fill the only slot

	acquired := make(chan struct{})
	go func() {
		ctx := context.Background()
		_ = th.Acquire(ctx)
		close(acquired)
	}()

	time.Sleep(10 * time.Millisecond)
	th.Release()

	select {
	case <-acquired:
		// success
	case <-time.After(200 * time.Millisecond):
		t.Error("timed out waiting for queued Acquire to unblock")
	}
}
