package jobnotify

import (
	"errors"
	"sync"
	"time"
)

// Event represents a notification event type.
type Event string

const (
	EventStarted  Event = "started"
	EventFinished Event = "finished"
	EventFailed   Event = "failed"
)

// Rule defines when a notification should be fired for a job.
type Rule struct {
	JobName  string
	On       []Event
	Channels []string
}

// Entry records a dispatched notification.
type Entry struct {
	JobName   string
	Event     Event
	Channel   string
	Dispached time.Time
}

// Store manages per-job notification rules and dispatch history.
type Store struct {
	mu      sync.RWMutex
	rules   map[string]*Rule
	history []Entry
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		rules: make(map[string]*Rule),
	}
}

// Register adds or replaces a notification rule for a job.
func (s *Store) Register(r Rule) error {
	if r.JobName == "" {
		return errors.New("jobnotify: job name must not be empty")
	}
	if len(r.On) == 0 {
		return errors.New("jobnotify: at least one event is required")
	}
	if len(r.Channels) == 0 {
		return errors.New("jobnotify: at least one channel is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := r
	s.rules[r.JobName] = &copy
	return nil
}

// ShouldNotify reports whether a notification should be sent for the given
// job and event, and returns the matching channels.
func (s *Store) ShouldNotify(jobName string, ev Event) ([]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[jobName]
	if !ok {
		return nil, false
	}
	for _, e := range r.On {
		if e == ev {
			return r.Channels, true
		}
	}
	return nil, false
}

// Record appends a dispatched notification entry to the history.
func (s *Store) Record(jobName string, ev Event, channel string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = append(s.history, Entry{
		JobName:   jobName,
		Event:     ev,
		Channel:   channel,
		Dispached: time.Now().UTC(),
	})
}

// History returns a copy of all dispatched notification entries.
func (s *Store) History() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, len(s.history))
	copy(out, s.history)
	return out
}

// Remove deletes the notification rule for a job.
func (s *Store) Remove(jobName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.rules, jobName)
}
