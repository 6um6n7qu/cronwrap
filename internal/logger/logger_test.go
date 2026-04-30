package logger_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/logger"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	l := logger.New("myjob", nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestInfo_WritesJSONEntry(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New("test-job", &buf)
	l.Info("hello world")

	var entry logger.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}
	if entry.Level != logger.LevelInfo {
		t.Errorf("expected level INFO, got %s", entry.Level)
	}
	if entry.Message != "hello world" {
		t.Errorf("expected message 'hello world', got %s", entry.Message)
	}
	if entry.Job != "test-job" {
		t.Errorf("expected job 'test-job', got %s", entry.Job)
	}
}

func TestJobFinished_IncludesExitCodeAndDuration(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New("finish-job", &buf)
	l.JobFinished(0, 2*time.Second)

	var entry logger.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}
	if entry.ExitCode == nil || *entry.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %v", entry.ExitCode)
	}
	if entry.Duration != "2s" {
		t.Errorf("expected duration '2s', got %s", entry.Duration)
	}
	if entry.Level != logger.LevelInfo {
		t.Errorf("expected level INFO, got %s", entry.Level)
	}
}

func TestJobFailed_WritesErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New("fail-job", &buf)
	l.JobFailed(1, 500*time.Millisecond)

	var entry logger.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}
	if entry.Level != logger.LevelError {
		t.Errorf("expected level ERROR, got %s", entry.Level)
	}
	if entry.ExitCode == nil || *entry.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %v", entry.ExitCode)
	}
}

func TestJobStarted_WritesInfoMessage(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New("start-job", &buf)
	l.JobStarted()

	var entry logger.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}
	if entry.Message != "job started" {
		t.Errorf("expected 'job started', got %s", entry.Message)
	}
}
