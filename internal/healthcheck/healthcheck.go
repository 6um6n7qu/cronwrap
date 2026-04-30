package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Status represents the health status of the cronwrap process.
type Status struct {
	Healthy     bool              `json:"healthy"`
	LastChecked time.Time         `json:"last_checked"`
	Jobs        map[string]JobStatus `json:"jobs"`
}

// JobStatus holds the last known result for a single job.
type JobStatus struct {
	Name        string    `json:"name"`
	LastRun     time.Time `json:"last_run"`
	LastSuccess time.Time `json:"last_success,omitempty"`
	LastFailure time.Time `json:"last_failure,omitempty"`
	ConsecutiveFails int  `json:"consecutive_fails"`
	Healthy     bool      `json:"healthy"`
}

// Checker maintains job health state and exposes an HTTP handler.
type Checker struct {
	mu   sync.RWMutex
	jobs map[string]*JobStatus
}

// New creates a new Checker instance.
func New() *Checker {
	return &Checker{
		jobs: make(map[string]*JobStatus),
	}
}

// RecordSuccess marks a job run as successful.
func (c *Checker) RecordSuccess(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UTC()
	s := c.getOrCreate(name)
	s.LastRun = now
	s.LastSuccess = now
	s.ConsecutiveFails = 0
	s.Healthy = true
}

// RecordFailure marks a job run as failed.
func (c *Checker) RecordFailure(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UTC()
	s := c.getOrCreate(name)
	s.LastRun = now
	s.LastFailure = now
	s.ConsecutiveFails++
	s.Healthy = false
}

// Handler returns an http.HandlerFunc that serves the health status as JSON.
func (c *Checker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.mu.RLock()
		defer c.mu.RUnlock()

		overallHealthy := true
		jobsCopy := make(map[string]JobStatus, len(c.jobs))
		for k, v := range c.jobs {
			jobsCopy[k] = *v
			if !v.Healthy {
				overallHealthy = false
			}
		}

		status := Status{
			Healthy:     overallHealthy,
			LastChecked: time.Now().UTC(),
			Jobs:        jobsCopy,
		}

		w.Header().Set("Content-Type", "application/json")
		if !overallHealthy {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		_ = json.NewEncoder(w).Encode(status)
	}
}

func (c *Checker) getOrCreate(name string) *JobStatus {
	if _, ok := c.jobs[name]; !ok {
		c.jobs[name] = &JobStatus{Name: name, Healthy: true}
	}
	return c.jobs[name]
}
