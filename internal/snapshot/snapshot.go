package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry holds a point-in-time record of a job's last execution.
type Entry struct {
	JobName   string        `json:"job_name"`
	StartedAt time.Time     `json:"started_at"`
	FinishedAt time.Time    `json:"finished_at"`
	Duration  time.Duration `json:"duration_ns"`
	ExitCode  int           `json:"exit_code"`
	Success   bool          `json:"success"`
}

// Store persists the most-recent execution snapshot for each job.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
	path    string
}

// New creates a Store backed by the given file path.
// Existing data is loaded if the file is present.
func New(path string) (*Store, error) {
	s := &Store{
		entries: make(map[string]Entry),
		path:    path,
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("snapshot: load: %w", err)
	}
	return s, nil
}

// Record saves an entry for the given job, overwriting any previous value.
func (s *Store) Record(e Entry) error {
	if e.JobName == "" {
		return fmt.Errorf("snapshot: job name must not be empty")
	}
	s.mu.Lock()
	s.entries[e.JobName] = e
	s.mu.Unlock()
	return s.persist()
}

// Get returns the latest snapshot for a job, or false if none exists.
func (s *Store) Get(jobName string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[jobName]
	return e, ok
}

// All returns a copy of every stored entry keyed by job name.
func (s *Store) All() map[string]Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]Entry, len(s.entries))
	for k, v := range s.entries {
		out[k] = v
	}
	return out
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.entries)
}

func (s *Store) persist() error {
	s.mu.RLock()
	data, err := json.MarshalIndent(s.entries, "", "  ")
	s.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	return os.WriteFile(s.path, data, 0o644)
}
