// Package throttle provides concurrency limiting for cron job execution.
// It ensures that no more than a configured number of jobs run simultaneously,
// preventing resource exhaustion under high load.
package throttle

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrThrottled is returned when the concurrency limit has been reached.
var ErrThrottled = errors.New("throttle: concurrency limit reached")

// Throttle limits the number of concurrent job executions.
type Throttle struct {
	mu      sync.Mutex
	sem     chan struct{}
	limit   int
	active  int
}

// New creates a new Throttle with the given concurrency limit.
// limit must be greater than zero.
func New(limit int) (*Throttle, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("throttle: limit must be greater than zero, got %d", limit)
	}
	return &Throttle{
		sem:   make(chan struct{}, limit),
		limit: limit,
	}, nil
}

// Acquire attempts to acquire a slot. It blocks until a slot is available
// or the context is cancelled. Returns ErrThrottled if the context is done.
func (t *Throttle) Acquire(ctx context.Context) error {
	select {
	case t.sem <- struct{}{}:
		t.mu.Lock()
		t.active++
		t.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ErrThrottled
	}
}

// TryAcquire attempts to acquire a slot without blocking.
// Returns ErrThrottled immediately if no slot is available.
func (t *Throttle) TryAcquire() error {
	select {
	case t.sem <- struct{}{}:
		t.mu.Lock()
		t.active++
		t.mu.Unlock()
		return nil
	default:
		return ErrThrottled
	}
}

// Release frees a previously acquired slot.
func (t *Throttle) Release() {
	select {
	case <-t.sem:
		t.mu.Lock()
		t.active--
		t.mu.Unlock()
	default:
	}
}

// Active returns the number of currently active (acquired) slots.
func (t *Throttle) Active() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.active
}

// Limit returns the configured concurrency limit.
func (t *Throttle) Limit() int {
	return t.limit
}
