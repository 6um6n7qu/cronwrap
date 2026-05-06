// Package audit provides a structured audit trail for cron job executions,
// recording job lifecycle events to a persistent append-only log file.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// EventType classifies the kind of audit event being recorded.
type EventType string

const (
	EventJobStarted  EventType = "job.started"
	EventJobFinished EventType = "job.finished"
	EventJobFailed   EventType = "job.failed"
	EventJobSkipped  EventType = "job.skipped"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	JobName   string    `json:"job_name"`
	ExitCode  *int      `json:"exit_code,omitempty"`
	Duration  *float64  `json:"duration_seconds,omitempty"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes audit entries to an underlying writer.
type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

// New creates an audit Logger that writes to the given writer.
// Pass os.Stdout for console output or an *os.File for persistent storage.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w}
}

// Record writes a single audit entry. It is safe for concurrent use.
func (l *Logger) Record(e Entry) error {
	if e.JobName == "" {
		return fmt.Errorf("audit: job_name must not be empty")
	}
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	b = append(b, '\n')
	_, err = l.out.Write(b)
	return err
}

// JobStarted is a convenience wrapper that records a job.started event.
func (l *Logger) JobStarted(jobName string) error {
	return l.Record(Entry{Event: EventJobStarted, JobName: jobName})
}

// JobFinished records a job.finished event with exit code and elapsed duration.
func (l *Logger) JobFinished(jobName string, exitCode int, d time.Duration) error {
	secs := d.Seconds()
	return l.Record(Entry{
		Event:    EventJobFinished,
		JobName:  jobName,
		ExitCode: &exitCode,
		Duration: &secs,
	})
}

// JobFailed records a job.failed event with an optional descriptive message.
func (l *Logger) JobFailed(jobName string, exitCode int, msg string) error {
	return l.Record(Entry{
		Event:   EventJobFailed,
		JobName: jobName,
		ExitCode: func() *int { c := exitCode; return &c }(),
		Message: msg,
	})
}
