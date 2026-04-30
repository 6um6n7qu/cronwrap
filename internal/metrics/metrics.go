// Package metrics provides job execution metrics collection and reporting
// for cronwrap. It tracks run counts, durations, and failure rates.
package metrics

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// JobMetrics holds aggregated execution statistics for a single job.
type JobMetrics struct {
	JobName      string        `json:"job_name"`
	TotalRuns    int           `json:"total_runs"`
	SuccessRuns  int           `json:"success_runs"`
	FailureRuns  int           `json:"failure_runs"`
	LastRunAt    time.Time     `json:"last_run_at"`
	LastDuration time.Duration `json:"last_duration_ns"`
	AvgDuration  time.Duration `json:"avg_duration_ns"`
	totalNs      int64
}

// Collector accumulates metrics across job executions.
type Collector struct {
	mu   sync.Mutex
	jobs map[string]*JobMetrics
}

// New creates a new Collector instance.
func New() *Collector {
	return &Collector{
		jobs: make(map[string]*JobMetrics),
	}
}

// Record registers the result of a single job execution.
func (c *Collector) Record(jobName string, duration time.Duration, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	m, ok := c.jobs[jobName]
	if !ok {
		m = &JobMetrics{JobName: jobName}
		c.jobs[jobName] = m
	}

	m.TotalRuns++
	m.LastRunAt = time.Now().UTC()
	m.LastDuration = duration
	m.totalNs += duration.Nanoseconds()
	m.AvgDuration = time.Duration(m.totalNs / int64(m.TotalRuns))

	if success {
		m.SuccessRuns++
	} else {
		m.FailureRuns++
	}
}

// Get returns a copy of the metrics for the given job name.
// The second return value is false if no metrics exist for that job.
func (c *Collector) Get(jobName string) (JobMetrics, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	m, ok := c.jobs[jobName]
	if !ok {
		return JobMetrics{}, false
	}
	return *m, true
}

// WriteJSON serialises all collected metrics as a JSON array to the given file path.
func (c *Collector) WriteJSON(path string) error {
	c.mu.Lock()
	slice := make([]JobMetrics, 0, len(c.jobs))
	for _, m := range c.jobs {
		slice = append(slice, *m)
	}
	c.mu.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(slice)
}
