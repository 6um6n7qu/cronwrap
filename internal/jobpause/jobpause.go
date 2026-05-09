// Package jobpause provides a mechanism to pause and resume scheduled jobs
// at runtime without modifying the underlying cron schedule.
package jobpause

import (
	"errors"
	"sync"
	"time"
)

// ErrUnknownJob is returned when an operation targets a job that has not been registered.
var ErrUnknownJob = errors.New("jobpause: unknown job")

// ErrEmptyName is returned when an empty job name is provided.
var ErrEmptyName = errors.New("jobpause: job name must not be empty")

// State holds the pause state for a single job.
type State struct {
	Paused    bool      `json:"paused"`
	PausedAt  time.Time `json:"paused_at,omitempty"`
	ResumedAt time.Time `json:"resumed_at,omitempty"`
}

// Store manages pause states for a collection of jobs.
type Store struct {
	mu     sync.RWMutex
	states map[string]*State
}

// New creates an empty Store.
func New() *Store {
	return &Store{states: make(map[string]*State)}
}

// Register adds a job to the store in the unpaused state.
func (s *Store) Register(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.states[name]; !ok {
		s.states[name] = &State{}
	}
	return nil
}

// Pause marks the named job as paused.
func (s *Store) Pause(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.states[name]
	if !ok {
		return ErrUnknownJob
	}
	st.Paused = true
	st.PausedAt = time.Now().UTC()
	return nil
}

// Resume marks the named job as active.
func (s *Store) Resume(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.states[name]
	if !ok {
		return ErrUnknownJob
	}
	st.Paused = false
	st.ResumedAt = time.Now().UTC()
	return nil
}

// IsPaused reports whether the named job is currently paused.
// Returns false for unknown jobs.
func (s *Store) IsPaused(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.states[name]
	if !ok {
		return false
	}
	return st.Paused
}

// All returns a snapshot of every registered job's state.
func (s *Store) All() map[string]State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]State, len(s.states))
	for k, v := range s.states {
		out[k] = *v
	}
	return out
}
