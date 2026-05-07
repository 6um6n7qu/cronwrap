// Package throttle provides a concurrency limiter for cron job execution.
//
// It uses a buffered channel as a semaphore to cap the number of jobs
// that may run at the same time. This prevents resource exhaustion when
// many scheduled jobs fire simultaneously.
//
// Basic usage:
//
//	th, err := throttle.New(4) // allow up to 4 concurrent jobs
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Non-blocking attempt:
//	if err := th.TryAcquire(); err != nil {
//		// job is throttled — skip or queue
//	}
//	defer th.Release()
//
//	// Blocking attempt with context:
//	if err := th.Acquire(ctx); err != nil {
//		// context cancelled before slot was available
//	}
//	defer th.Release()
package throttle
