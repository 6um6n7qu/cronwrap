package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/example/cronwrap/internal/audit"
)

func TestRecord_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.Record(audit.Entry{
		Event:   audit.EventJobStarted,
		JobName: "backup",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var entry audit.Entry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if entry.JobName != "backup" {
		t.Errorf("expected job_name=backup, got %q", entry.JobName)
	}
	if entry.Event != audit.EventJobStarted {
		t.Errorf("expected event=%s, got %s", audit.EventJobStarted, entry.Event)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_EmptyJobNameReturnsError(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.Record(audit.Entry{Event: audit.EventJobStarted})
	if err == nil {
		t.Fatal("expected error for empty job_name, got nil")
	}
}

func TestJobFinished_IncludesExitCodeAndDuration(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.JobFinished("cleanup", 0, 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.ExitCode == nil || *entry.ExitCode != 0 {
		t.Errorf("expected exit_code=0, got %v", entry.ExitCode)
	}
	if entry.Duration == nil || *entry.Duration != 2.0 {
		t.Errorf("expected duration_seconds=2.0, got %v", entry.Duration)
	}
}

func TestJobFailed_WritesFailedEvent(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.JobFailed("sync", 1, "connection refused")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.Event != audit.EventJobFailed {
		t.Errorf("expected event=%s, got %s", audit.EventJobFailed, entry.Event)
	}
	if entry.Message != "connection refused" {
		t.Errorf("unexpected message: %q", entry.Message)
	}
}

func TestNew_NilWriterDefaultsToStdout(t *testing.T) {
	// Should not panic when nil writer is passed.
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil Logger")
	}
}
