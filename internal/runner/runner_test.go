package runner

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRun_Success(t *testing.T) {
	ctx := context.Background()
	result := Run(ctx, "echo", "hello")

	if !result.Success() {
		t.Fatalf("expected success, got exit code %d: %v", result.ExitCode, result.Error)
	}
	if !strings.Contains(result.Stdout, "hello") {
		t.Errorf("expected stdout to contain 'hello', got %q", result.Stdout)
	}
	if result.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", result.Duration)
	}
}

func TestRun_Failure(t *testing.T) {
	ctx := context.Background()
	result := Run(ctx, "false")

	if result.Success() {
		t.Fatal("expected failure but got success")
	}
	if result.ExitCode == 0 {
		t.Errorf("expected non-zero exit code")
	}
}

func TestRun_CommandNotFound(t *testing.T) {
	ctx := context.Background()
	result := Run(ctx, "__nonexistent_cmd_cronwrap__")

	if result.Success() {
		t.Fatal("expected failure for missing command")
	}
	if result.Error == nil {
		t.Error("expected non-nil error for missing command")
	}
}

func TestRun_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result := Run(ctx, "sleep", "10")

	if result.Success() {
		t.Fatal("expected failure due to context cancellation")
	}
	if result.ExitCode == 0 {
		t.Error("expected non-zero exit code after cancellation")
	}
}

func TestRun_StderrCapture(t *testing.T) {
	ctx := context.Background()
	result := Run(ctx, "sh", "-c", "echo error-output >&2; exit 1")

	if result.Success() {
		t.Fatal("expected failure")
	}
	if !strings.Contains(result.Stderr, "error-output") {
		t.Errorf("expected stderr to contain 'error-output', got %q", result.Stderr)
	}
}
