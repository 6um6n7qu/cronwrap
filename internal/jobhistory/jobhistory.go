// Package jobhistory maintains a bounded, in-memory ring buffer of completed
// job executions so that operators can inspect recent run outcomes via the
// HTTP handler without requiring an external datastore.
package jobhistory

import (
	"errors"
	"sync"
	"time"
)

// Entry describes a single completed job execution.
type Entry struct {
	JobName   string        `json:"job_name"`
	StartedAt time.Time     `json:"started_at"`
	FinishedAt time.Time    `json:"finished_at"`
	Duration  time.Duration `json:"duration_ms"`
	ExitCode  int           `json:"exit_code"`
	Success   bool          `json:"success"`
	Error     string        `json:"error,omitempty"`
}

// Store is a thread-safe ring buffer of job history entries.
type Store struct {
	mu      sync.RWMutex
	entries []Entry
	cap     int
	head    int
	size    int
}

// New creates a Store that retains at most capacity entries.
// Returns an error if capacity is less than 1.
func New(capacity int) (*Store, error) {
	if capacity < 1 {
		return nil, errors.New("jobhistory: capacity must be at least 1")
	}
	return &Store{
		entries: make([]Entry, capacity),
		cap:     capacity,
	}, nil
}

// Record appends an entry to the ring buffer, evicting the oldest entry when
// the buffer is full. Returns an error if JobName is empty.
func (s *Store) Record(e Entry) error {
	if e.JobName == "" {
		return errors.New("jobhistory: job name must not be empty")
	}
	if e.FinishedAt.IsZero() {
		e.FinishedAt = time.Now()
	}
	if !e.StartedAt.IsZero() {
		e.Duration = e.FinishedAt.Sub(e.StartedAt)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[s.head] = e
	s.head = (s.head + 1) % s.cap
	if s.size < s.cap {
		s.size++
	}
	return nil
}

// All returns a slice of all stored entries ordered from oldest to newest.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.size == 0 {
		return []Entry{}
	}
	out := make([]Entry, s.size)
	start := (s.head - s.size + s.cap) % s.cap
	for i := 0; i < s.size; i++ {
		out[i] = s.entries[(start+i)%s.cap]
	}
	return out
}

// ForJob returns all stored entries for the given job name.
func (s *Store) ForJob(name string) []Entry {
	all := s.All()
	var out []Entry
	for _, e := range all {
		if e.JobName == name {
			out = append(out, e)
		}
	}
	if out == nil {
		return []Entry{}
	}
	return out
}
