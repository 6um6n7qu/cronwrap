// Package ratelimit implements a token-bucket rate limiter for cronwrap jobs.
//
// It prevents a single job from being triggered too frequently — useful when
// a scheduler misconfiguration or external trigger could cause a burst of
// concurrent executions that would overwhelm downstream systems.
//
// # Design
//
// Each unique job name gets its own independent token bucket. Tokens refill at
// a fixed interval (the "period" per token), and a configurable burst size
// controls how many executions can occur back-to-back before throttling kicks
// in.
//
// # Basic usage
//
// Allow at most 3 executions per minute for any job:
//
//	l, err := ratelimit.New(20*time.Second, 3)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if !l.Allow(job.Name) {
//		log.Printf("job %s rate-limited, skipping", job.Name)
//		return
//	}
//
// # Blocking usage
//
// For blocking behaviour, use Wait, which respects context cancellation:
//
//	if err := l.Wait(ctx, job.Name); err != nil {
//		log.Printf("context cancelled while waiting for rate-limit: %v", err)
//	}
package ratelimit
