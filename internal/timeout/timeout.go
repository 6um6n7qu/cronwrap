// Package timeout provides utilities for enforcing execution time limits
// on cron job commands, integrating with the runner and scheduler.
package timeout

import (
	"context"
	"fmt"
	"time"
)

// ErrDeadlineExceeded is returned when a job exceeds its allowed duration.
var ErrDeadlineExceeded = fmt.Errorf("job exceeded deadline")

// Policy defines how a timeout is applied to a job execution.
type Policy struct {
	// Duration is the maximum allowed execution time.
	Duration time.Duration
	// GracePeriod is additional time allowed for cleanup after the deadline.
	GracePeriod time.Duration
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		Duration:    30 * time.Minute,
		GracePeriod: 5 * time.Second,
	}
}

// WithContext returns a derived context that is cancelled after the policy
// Duration elapses. The caller must invoke the returned cancel function.
func (p Policy) WithContext(parent context.Context) (context.Context, context.CancelFunc) {
	if p.Duration <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, p.Duration)
}

// Wrap executes fn within a context governed by the policy. It returns
// ErrDeadlineExceeded when the timeout fires before fn returns.
func (p Policy) Wrap(parent context.Context, fn func(ctx context.Context) error) error {
	ctx, cancel := p.WithContext(parent)
	defer cancel()

	type result struct {
		err error
	}

	ch := make(chan result, 1)
	go func() {
		ch <- result{err: fn(ctx)}
	}()

	select {
	case res := <-ch:
		return res.err
	case <-ctx.Done():
		if p.GracePeriod > 0 {
			select {
			case res := <-ch:
				return res.err
			case <-time.After(p.GracePeriod):
			}
		}
		return fmt.Errorf("%w: limit=%s", ErrDeadlineExceeded, p.Duration)
	}
}
