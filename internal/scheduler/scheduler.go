package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// Job represents a scheduled cron job entry.
type Job struct {
	Name     string
	Schedule string
	Command  string
	Args     []string
}

// Scheduler wraps a cron instance and manages job lifecycle.
type Scheduler struct {
	c    *cron.Cron
	jobs []Job
}

// New creates a new Scheduler with second-level precision disabled.
func New() *Scheduler {
	return &Scheduler{
		c: cron.New(cron.WithLocation(time.UTC)),
	}
}

// Add registers a job with the scheduler. Returns an error if the
// cron expression is invalid.
func (s *Scheduler) Add(j Job) error {
	if j.Name == "" {
		return fmt.Errorf("scheduler: job name must not be empty")
	}
	if j.Schedule == "" {
		return fmt.Errorf("scheduler: schedule for job %q must not be empty", j.Name)
	}
	if j.Command == "" {
		return fmt.Errorf("scheduler: command for job %q must not be empty", j.Name)
	}

	_, err := s.c.AddFunc(j.Schedule, func() {})
	if err != nil {
		return fmt.Errorf("scheduler: invalid schedule %q for job %q: %w", j.Schedule, j.Name, err)
	}

	s.jobs = append(s.jobs, j)
	return nil
}

// Jobs returns a copy of the registered jobs slice.
func (s *Scheduler) Jobs() []Job {
	out := make([]Job, len(s.jobs))
	copy(out, s.jobs)
	return out
}

// Start begins the scheduler in the background. It is a no-op if no
// jobs have been registered.
func (s *Scheduler) Start(ctx context.Context) {
	if len(s.jobs) == 0 {
		return
	}
	s.c.Start()
	go func() {
		<-ctx.Done()
		s.c.Stop()
	}()
}
