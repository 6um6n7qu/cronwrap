// Package jobtag provides tagging support for cron jobs, allowing jobs to be
// grouped and filtered by arbitrary string tags.
package jobtag

import (
	"errors"
	"sync"
)

// ErrJobNotFound is returned when an operation targets an unregistered job.
var ErrJobNotFound = errors.New("jobtag: job not found")

// ErrEmptyName is returned when an empty job name is provided.
var ErrEmptyName = errors.New("jobtag: job name must not be empty")

// ErrEmptyTag is returned when an empty tag value is provided.
var ErrEmptyTag = errors.New("jobtag: tag must not be empty")

// Store holds tag associations for registered jobs.
type Store struct {
	mu   sync.RWMutex
	jobs map[string]map[string]struct{}
}

// New creates and returns an empty Store.
func New() *Store {
	return &Store{
		jobs: make(map[string]map[string]struct{}),
	}
}

// Register adds a job to the store with no tags.
func (s *Store) Register(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[name]; !ok {
		s.jobs[name] = make(map[string]struct{})
	}
	return nil
}

// Add attaches a tag to the named job.
func (s *Store) Add(name, tag string) error {
	if name == "" {
		return ErrEmptyName
	}
	if tag == "" {
		return ErrEmptyTag
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	tags, ok := s.jobs[name]
	if !ok {
		return ErrJobNotFound
	}
	tags[tag] = struct{}{}
	return nil
}

// Remove detaches a tag from the named job. It is a no-op if the tag is absent.
func (s *Store) Remove(name, tag string) error {
	if name == "" {
		return ErrEmptyName
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[name]; !ok {
		return ErrJobNotFound
	}
	delete(s.jobs[name], tag)
	return nil
}

// Tags returns a sorted slice of tags for the named job.
func (s *Store) Tags(name string) ([]string, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	tags, ok := s.jobs[name]
	if !ok {
		return nil, ErrJobNotFound
	}
	out := make([]string, 0, len(tags))
	for t := range tags {
		out = append(out, t)
	}
	return out, nil
}

// JobsWithTag returns the names of all jobs that carry the given tag.
func (s *Store) JobsWithTag(tag string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []string
	for name, tags := range s.jobs {
		if _, ok := tags[tag]; ok {
			out = append(out, name)
		}
	}
	return out
}
