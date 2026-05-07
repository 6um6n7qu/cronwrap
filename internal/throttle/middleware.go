package throttle

import (
	"context"
	"fmt"
)

// JobFunc is the signature for a runnable cron job.
type JobFunc func(ctx context.Context) error

// Wrap returns a new JobFunc that applies throttling before invoking fn.
// If a slot cannot be acquired (context cancelled or limit reached),
// the job is skipped and ErrThrottled is returned without calling fn.
//
// Example integration with the scheduler:
//
//	th, _ := throttle.New(4)
//	wrapped := th.Wrap(myJobFunc)
//	scheduler.Add(job.Name, job.Schedule, wrapped)
func (t *Throttle) Wrap(fn JobFunc) JobFunc {
	return func(ctx context.Context) error {
		if err := t.Acquire(ctx); err != nil {
			return fmt.Errorf("%w: job skipped due to concurrency limit", ErrThrottled)
		}
		defer t.Release()
		return fn(ctx)
	}
}

// WrapNonBlocking returns a new JobFunc that uses TryAcquire so the job
// is dropped immediately if the limit is already reached, rather than
// waiting for a free slot.
func (t *Throttle) WrapNonBlocking(fn JobFunc) JobFunc {
	return func(ctx context.Context) error {
		if err := t.TryAcquire(); err != nil {
			return fmt.Errorf("%w: job dropped, concurrency limit reached", ErrThrottled)
		}
		defer t.Release()
		return fn(ctx)
	}
}
