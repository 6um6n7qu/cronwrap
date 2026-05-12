// Package jobstatus tracks the current execution status of registered cron jobs.
package jobstatus

import (
	"errors"
	"sync"
	"time"
)

// Status represents the lifecycle state of a job.
type Status string

const (
	StatusIdle     Status = "idle"
	StatusRunning  Status = "running"
	StatusSuccess  Status = "success"
	StatusFailed   Status = "failed"
)

// Entry holds the current status and timing information for a single job.
type Entry struct {
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	StartedAt time.Time `json:"started_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	ExitCode  int       `json:"exit_code,omitempty"`
}

// Store manages job status entries.
type Store struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New creates a new empty Store.
func New() *Store {
	return &Store{entries: make(map[string]*Entry)}
}

// Register adds a job to the store with an initial idle status.
func (s *Store) Register(name string) error {
	if name == "" {
		return errors.New("jobstatus: name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.entries[name]; exists {
		return errors.New("jobstatus: job already registered: " + name)
	}
	s.entries[name] = &Entry{Name: name, Status: StatusIdle, UpdatedAt: time.Now()}
	return nil
}

// SetRunning marks a job as currently running.
func (s *Store) SetRunning(name string) error {
	return s.update(name, func(e *Entry) {
		e.Status = StatusRunning
		e.StartedAt = time.Now()
		e.ExitCode = 0
	})
}

// SetSuccess marks a job as successfully completed.
func (s *Store) SetSuccess(name string) error {
	return s.update(name, func(e *Entry) {
		e.Status = StatusSuccess
		e.ExitCode = 0
	})
}

// SetFailed marks a job as failed with the given exit code.
func (s *Store) SetFailed(name string, exitCode int) error {
	return s.update(name, func(e *Entry) {
		e.Status = StatusFailed
		e.ExitCode = exitCode
	})
}

// Get returns the status entry for the named job.
func (s *Store) Get(name string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[name]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all registered job entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, *e)
	}
	return out
}

func (s *Store) update(name string, fn func(*Entry)) error {
	if name == "" {
		return errors.New("jobstatus: name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[name]
	if !ok {
		return errors.New("jobstatus: unknown job: " + name)
	}
	fn(e)
	e.UpdatedAt = time.Now()
	return nil
}
