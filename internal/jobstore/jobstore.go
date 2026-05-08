package jobstore

import (
	"errors"
	"sync"
	"time"
)

// Status represents the current state of a tracked job.
type Status string

const (
	StatusPending  Status = "pending"
	StatusRunning  Status = "running"
	StatusSuccess  Status = "success"
	StatusFailed   Status = "failed"
)

// Entry holds runtime metadata for a single job execution.
type Entry struct {
	Name       string
	Status     Status
	StartedAt  time.Time
	FinishedAt time.Time
	ExitCode   int
	Error      string
}

// Store is a thread-safe in-memory registry of job execution states.
type Store struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New returns an initialised Store.
func New() *Store {
	return &Store{entries: make(map[string]*Entry)}
}

// Start records a job as running. Returns an error if the name is empty.
func (s *Store) Start(name string) error {
	if name == "" {
		return errors.New("jobstore: job name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[name] = &Entry{
		Name:      name,
		Status:    StatusRunning,
		StartedAt: time.Now(),
	}
	return nil
}

// Finish marks a job as succeeded or failed based on exitCode.
func (s *Store) Finish(name string, exitCode int, jobErr error) error {
	if name == "" {
		return errors.New("jobstore: job name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[name]
	if !ok {
		e = &Entry{Name: name}
		s.entries[name] = e
	}
	e.FinishedAt = time.Now()
	e.ExitCode = exitCode
	if jobErr != nil {
		e.Status = StatusFailed
		e.Error = jobErr.Error()
	} else {
		e.Status = StatusSuccess
	}
	return nil
}

// Get returns a copy of the entry for the given job name.
func (s *Store) Get(name string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[name]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of every tracked entry.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, *e)
	}
	return out
}
