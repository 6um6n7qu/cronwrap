// Package joblabel provides key-value label management for cron jobs,
// allowing jobs to be tagged and filtered by arbitrary metadata.
package joblabel

import (
	"errors"
	"fmt"
	"sync"
)

// ErrJobNotFound is returned when operating on an unregistered job.
var ErrJobNotFound = errors.New("joblabel: job not found")

// ErrEmptyName is returned when an empty job name is provided.
var ErrEmptyName = errors.New("joblabel: job name must not be empty")

// ErrEmptyKey is returned when an empty label key is provided.
var ErrEmptyKey = errors.New("joblabel: label key must not be empty")

// Labels is a map of string key-value pairs attached to a job.
type Labels map[string]string

// Store manages labels for registered jobs.
type Store struct {
	mu   sync.RWMutex
	jobs map[string]Labels
}

// New creates an empty label Store.
func New() *Store {
	return &Store{jobs: make(map[string]Labels)}
}

// Register adds a job to the store with no initial labels.
func (s *Store) Register(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[name]; !ok {
		s.jobs[name] = make(Labels)
	}
	return nil
}

// Set assigns a label key-value pair to the named job.
func (s *Store) Set(name, key, value string) error {
	if name == "" {
		return ErrEmptyName
	}
	if key == "" {
		return ErrEmptyKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	labels, ok := s.jobs[name]
	if !ok {
		return fmt.Errorf("%w: %s", ErrJobNotFound, name)
	}
	labels[key] = value
	return nil
}

// Get returns all labels for the named job.
func (s *Store) Get(name string) (Labels, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	labels, ok := s.jobs[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrJobNotFound, name)
	}
	copy := make(Labels, len(labels))
	for k, v := range labels {
		copy[k] = v
	}
	return copy, nil
}

// Delete removes a label key from the named job.
func (s *Store) Delete(name, key string) error {
	if name == "" {
		return ErrEmptyName
	}
	if key == "" {
		return ErrEmptyKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	labels, ok := s.jobs[name]
	if !ok {
		return fmt.Errorf("%w: %s", ErrJobNotFound, name)
	}
	delete(labels, key)
	return nil
}

// Match returns the names of all jobs that have every label in the selector.
func (s *Store) Match(selector Labels) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var matched []string
	for name, labels := range s.jobs {
		if matchesAll(labels, selector) {
			matched = append(matched, name)
		}
	}
	return matched
}

func matchesAll(labels, selector Labels) bool {
	for k, v := range selector {
		if labels[k] != v {
			return false
		}
	}
	return true
}
