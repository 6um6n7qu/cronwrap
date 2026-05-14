// Package jobruncount tracks how many times each job has been executed.
package jobruncount

import (
	"errors"
	"sync"
)

// Entry holds run count statistics for a single job.
type Entry struct {
	Name       string `json:"name"`
	TotalRuns  int64  `json:"total_runs"`
	SuccessRuns int64 `json:"success_runs"`
	FailureRuns int64 `json:"failure_runs"`
}

// Store tracks run counts for registered jobs.
type Store struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		entries: make(map[string]*Entry),
	}
}

// Register adds a job to the store. Returns an error if the name is empty or
// the job is already registered.
func (s *Store) Register(name string) error {
	if name == "" {
		return errors.New("jobruncount: name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[name]; ok {
		return errors.New("jobruncount: job already registered: " + name)
	}
	s.entries[name] = &Entry{Name: name}
	return nil
}

// RecordSuccess increments the total and success counters for the named job.
func (s *Store) RecordSuccess(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[name]
	if !ok {
		return errors.New("jobruncount: unknown job: " + name)
	}
	e.TotalRuns++
	e.SuccessRuns++
	return nil
}

// RecordFailure increments the total and failure counters for the named job.
func (s *Store) RecordFailure(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[name]
	if !ok {
		return errors.New("jobruncount: unknown job: " + name)
	}
	e.TotalRuns++
	e.FailureRuns++
	return nil
}

// Get returns a copy of the Entry for the named job and whether it was found.
func (s *Store) Get(name string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[name]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, *e)
	}
	return out
}
