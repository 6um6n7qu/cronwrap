// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently cron jobs can be triggered or retried.
package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Limiter controls the rate at which a named job may execute.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     time.Duration // minimum interval between executions
	maxBurst int
}

type bucket struct {
	tokens    int
	lastRefil time.Time
}

// New creates a Limiter that allows at most maxBurst executions per job
// within any window of `rate` duration.
func New(rate time.Duration, maxBurst int) (*Limiter, error) {
	if rate <= 0 {
		return nil, fmt.Errorf("ratelimit: rate must be positive, got %v", rate)
	}
	if maxBurst < 1 {
		return nil, fmt.Errorf("ratelimit: maxBurst must be at least 1, got %d", maxBurst)
	}
	return &Limiter{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		maxBurst: maxBurst,
	}, nil
}

// Allow reports whether the job identified by name is permitted to run now.
// It consumes one token if available.
func (l *Limiter) Allow(name string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[name]
	if !ok {
		b = &bucket{tokens: l.maxBurst, lastRefil: now}
		l.buckets[name] = b
	}

	// Refill tokens proportional to elapsed time.
	elapsed := now.Sub(b.lastRefil)
	newTokens := int(elapsed / l.rate)
	if newTokens > 0 {
		b.tokens += newTokens
		if b.tokens > l.maxBurst {
			b.tokens = l.maxBurst
		}
		b.lastRefil = now
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// Wait blocks until a token is available for the named job or the context
// is cancelled.
func (l *Limiter) Wait(ctx context.Context, name string) error {
	for {
		if l.Allow(name) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(l.rate / 2):
		}
	}
}

// Reset clears the token bucket for the named job.
func (l *Limiter) Reset(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, name)
}
