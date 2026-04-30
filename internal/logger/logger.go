// Package logger provides structured JSON logging for cronwrap job execution.
package logger

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Level represents a log severity level.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelError Level = "ERROR"
	LevelDebug Level = "DEBUG"
)

// Entry represents a single structured log entry.
type Entry struct {
	Timestamp string `json:"timestamp"`
	Level     Level  `json:"level"`
	Job       string `json:"job"`
	Message   string `json:"message"`
	Duration  string `json:"duration,omitempty"`
	ExitCode  *int   `json:"exit_code,omitempty"`
}

// Logger writes structured JSON log entries.
type Logger struct {
	out io.Writer
	job string
}

// New creates a new Logger for the given job name.
// If out is nil, os.Stdout is used.
func New(job string, out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout
	}
	return &Logger{out: out, job: job}
}

// Info logs an informational message.
func (l *Logger) Info(msg string) {
	l.write(LevelInfo, msg, nil, "")
}

// Error logs an error message.
func (l *Logger) Error(msg string) {
	l.write(LevelError, msg, nil, "")
}

// JobStarted logs that a job has started.
func (l *Logger) JobStarted() {
	l.Info("job started")
}

// JobFinished logs that a job completed with the given exit code and duration.
func (l *Logger) JobFinished(exitCode int, duration time.Duration) {
	l.write(LevelInfo, "job finished", &exitCode, duration.String())
}

// JobFailed logs that a job failed with the given exit code and duration.
func (l *Logger) JobFailed(exitCode int, duration time.Duration) {
	l.write(LevelError, "job failed", &exitCode, duration.String())
}

func (l *Logger) write(level Level, msg string, exitCode *int, duration string) {
	entry := Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Job:       l.job,
		Message:   msg,
		Duration:  duration,
		ExitCode:  exitCode,
	}
	data, _ := json.Marshal(entry)
	l.out.Write(append(data, '\n'))
}
