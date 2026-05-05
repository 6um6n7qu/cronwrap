// Package retry provides configurable retry logic for cron job execution.
package retry

import (
	"context"
	"fmt"
	"time"
)

// Policy defines the retry behaviour for a job.
type Policy struct {
	// MaxAttempts is the total number of attempts (1 means no retries).
	MaxAttempts int
	// Delay is the wait duration between attempts.
	Delay time.Duration
	// Multiplier scales the delay after each failure (1.0 = constant delay).
	Multiplier float64
}

// DefaultPolicy returns a sensible default: 3 attempts, 5 s constant delay.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       5 * time.Second,
		Multiplier:  1.0,
	}
}

// Attempt represents the result of a single execution attempt.
type Attempt struct {
	Number int
	Err    error
	Took   time.Duration
}

// Do runs fn according to p, retrying on non-nil errors.
// It returns all attempts made and the final error (nil on success).
func Do(ctx context.Context, p Policy, fn func(ctx context.Context) error) ([]Attempt, error) {
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}
	if p.Multiplier < 1.0 {
		p.Multiplier = 1.0
	}

	var attempts []Attempt
	delay := p.Delay

	for i := 1; i <= p.MaxAttempts; i++ {
		if err := ctx.Err(); err != nil {
			return attempts, fmt.Errorf("retry aborted: %w", err)
		}

		start := time.Now()
		err := fn(ctx)
		attempts = append(attempts, Attempt{Number: i, Err: err, Took: time.Since(start)})

		if err == nil {
			return attempts, nil
		}

		if i < p.MaxAttempts {
			select {
			case <-ctx.Done():
				return attempts, fmt.Errorf("retry aborted during backoff: %w", ctx.Err())
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * p.Multiplier)
		}
	}

	return attempts, attempts[len(attempts)-1].Err
}
