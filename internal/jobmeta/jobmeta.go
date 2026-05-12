package jobmeta

import (
	"errors"
	"sync"
	"time"
)

// Entry holds arbitrary metadata key-value pairs attached to a job.
type Entry struct {
	Name      string            `json:"name"`
	Meta      map[string]string `json:"meta"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Store manages per-job metadata.
type Store struct {
	mu   sync.RWMutex
	jobs map[string]*Entry
}

// New returns an initialised Store.
func New() *Store {
	return &Store{jobs: make(map[string]*Entry)}
}

// Register adds a job to the store. Returns an error if name is empty or
// the job already exists.
func (s *Store) Register(name string) error {
	if name == "" {
		return errors.New("jobmeta: name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[name]; ok {
		return errors.New("jobmeta: job already registered: " + name)
	}
	s.jobs[name] = &Entry{Name: name, Meta: make(map[string]string)}
	return nil
}

// Set assigns a metadata value for key on the named job.
func (s *Store) Set(name, key, value string) error {
	if name == "" {
		return errors.New("jobmeta: name must not be empty")
	}
	if key == "" {
		return errors.New("jobmeta: key must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.jobs[name]
	if !ok {
		return errors.New("jobmeta: unknown job: " + name)
	}
	e.Meta[key] = value
	e.UpdatedAt = time.Now()
	return nil
}

// Delete removes a single metadata key from a job.
func (s *Store) Delete(name, key string) error {
	if name == "" {
		return errors.New("jobmeta: name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.jobs[name]
	if !ok {
		return errors.New("jobmeta: unknown job: " + name)
	}
	delete(e.Meta, key)
	e.UpdatedAt = time.Now()
	return nil
}

// Get returns the Entry for a job and whether it was found.
func (s *Store) Get(name string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.jobs[name]
	if !ok {
		return Entry{}, false
	}
	// Return a shallow copy so callers cannot mutate internal state.
	copy := *e
	copy.Meta = make(map[string]string, len(e.Meta))
	for k, v := range e.Meta {
		copy.Meta[k] = v
	}
	return copy, true
}

// All returns a snapshot of every registered job's metadata.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.jobs))
	for _, e := range s.jobs {
		copy := *e
		copy.Meta = make(map[string]string, len(e.Meta))
		for k, v := range e.Meta {
			copy.Meta[k] = v
		}
		out = append(out, copy)
	}
	return out
}
