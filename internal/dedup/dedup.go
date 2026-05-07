// Package dedup provides job deduplication to prevent overlapping cron executions.
// It uses a shared lock store to ensure only one instance of a named job runs at a time.
package dedup

import (
	"errors"
	"sync"
	"time"
)

// ErrAlreadyRunning is returned when a job with the same name is already active.
var ErrAlreadyRunning = errors.New("dedup: job is already running")

// Entry tracks an active job execution.
type Entry struct {
	StartedAt time.Time
}

// Deduplicator guards against concurrent execution of identically named jobs.
type Deduplicator struct {
	mu      sync.Mutex
	running map[string]Entry
}

// New creates a new Deduplicator.
func New() *Deduplicator {
	return &Deduplicator{
		running: make(map[string]Entry),
	}
}

// Acquire attempts to mark the named job as running.
// Returns ErrAlreadyRunning if the job is already active.
func (d *Deduplicator) Acquire(name string) error {
	if name == "" {
		return errors.New("dedup: job name must not be empty")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.running[name]; ok {
		return ErrAlreadyRunning
	}
	d.running[name] = Entry{StartedAt: time.Now()}
	return nil
}

// Release marks the named job as no longer running.
// It is safe to call Release for a job that is not currently tracked.
func (d *Deduplicator) Release(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.running, name)
}

// IsRunning reports whether the named job is currently active.
func (d *Deduplicator) IsRunning(name string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.running[name]
	return ok
}

// ActiveJobs returns a snapshot of all currently running job names and their start times.
func (d *Deduplicator) ActiveJobs() map[string]Entry {
	d.mu.Lock()
	defer d.mu.Unlock()
	snapshot := make(map[string]Entry, len(d.running))
	for k, v := range d.running {
		snapshot[k] = v
	}
	return snapshot
}
