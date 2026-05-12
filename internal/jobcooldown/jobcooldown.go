// Package jobcooldown enforces a minimum wait period between successive
// executions of the same job, preventing rapid re-runs after completion.
package jobcooldown

import (
	"errors"
	"sync"
	"time"
)

// ErrCooldownActive is returned when a job is requested before its cooldown
// period has elapsed.
var ErrCooldownActive = errors.New("jobcooldown: job is in cooldown period")

// ErrEmptyName is returned when an empty job name is provided.
var ErrEmptyName = errors.New("jobcooldown: job name must not be empty")

// ErrUnknownJob is returned when a job has not been registered.
var ErrUnknownJob = errors.New("jobcooldown: unknown job")

// entry holds cooldown state for a single job.
type entry struct {
	cooldown  time.Duration
	lastRanAt time.Time
}

// Store tracks per-job cooldown state.
type Store struct {
	mu   sync.RWMutex
	jobs map[string]*entry
	now  func() time.Time
}

// New creates a new Store. The now function is used to determine the current
// time; pass nil to use time.Now.
func New(now func() time.Time) *Store {
	if now == nil {
		now = time.Now
	}
	return &Store{
		jobs: make(map[string]*entry),
		now:  now,
	}
}

// Register adds a job with the given cooldown duration.
func (s *Store) Register(name string, cooldown time.Duration) error {
	if name == "" {
		return ErrEmptyName
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[name] = &entry{cooldown: cooldown}
	return nil
}

// Allow returns nil if the job may run, or ErrCooldownActive if it must wait.
// Returns ErrUnknownJob if the job has not been registered.
func (s *Store) Allow(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	s.mu.RLock()
	e, ok := s.jobs[name]
	s.mu.RUnlock()
	if !ok {
		return ErrUnknownJob
	}
	if e.lastRanAt.IsZero() {
		return nil
	}
	elapsed := s.now().Sub(e.lastRanAt)
	if elapsed < e.cooldown {
		return ErrCooldownActive
	}
	return nil
}

// Record marks a job as having just run, resetting its cooldown window.
func (s *Store) Record(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.jobs[name]
	if !ok {
		return ErrUnknownJob
	}
	e.lastRanAt = s.now()
	return nil
}

// Remaining returns the time left in the cooldown period for the given job.
// Returns 0 if the job may run immediately.
func (s *Store) Remaining(name string) (time.Duration, error) {
	if name == "" {
		return 0, ErrEmptyName
	}
	s.mu.RLock()
	e, ok := s.jobs[name]
	s.mu.RUnlock()
	if !ok {
		return 0, ErrUnknownJob
	}
	if e.lastRanAt.IsZero() {
		return 0, nil
	}
	remaining := e.cooldown - s.now().Sub(e.lastRanAt)
	if remaining < 0 {
		return 0, nil
	}
	return remaining, nil
}
