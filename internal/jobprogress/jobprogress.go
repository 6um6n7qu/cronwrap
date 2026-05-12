// Package jobprogress tracks the progress and completion percentage of running jobs.
package jobprogress

import (
	"errors"
	"sync"
	"time"
)

// Entry holds progress information for a single job.
type Entry struct {
	Name      string    `json:"name"`
	Total     int64     `json:"total"`
	Done      int64     `json:"done"`
	Percent   float64   `json:"percent"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store tracks progress for registered jobs.
type Store struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New creates a new Store.
func New() *Store {
	return &Store{entries: make(map[string]*Entry)}
}

// Register initialises progress tracking for a job with a known total.
func (s *Store) Register(name string, total int64) error {
	if name == "" {
		return errors.New("jobprogress: name must not be empty")
	}
	if total <= 0 {
		return errors.New("jobprogress: total must be greater than zero")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[name] = &Entry{
		Name:      name,
		Total:     total,
		UpdatedAt: time.Now(),
	}
	return nil
}

// Update sets the number of completed units for a job.
func (s *Store) Update(name string, done int64) error {
	if name == "" {
		return errors.New("jobprogress: name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[name]
	if !ok {
		return errors.New("jobprogress: unknown job: " + name)
	}
	if done < 0 {
		done = 0
	}
	if done > e.Total {
		done = e.Total
	}
	e.Done = done
	e.Percent = float64(done) / float64(e.Total) * 100
	e.UpdatedAt = time.Now()
	return nil
}

// Get returns the progress entry for a job.
func (s *Store) Get(name string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[name]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all progress entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, *e)
	}
	return out
}
