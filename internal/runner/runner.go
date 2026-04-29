package runner

import (
	"context"
	"os/exec"
	"time"
)

// Result holds the outcome of a command execution.
type Result struct {
	Command  string
	Args     []string
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
	Error    error
	Started  time.Time
	Finished time.Time
}

// Run executes the given command with args, respecting the provided context.
// It captures stdout, stderr, duration, and exit code.
func Run(ctx context.Context, command string, args ...string) Result {
	result := Result{
		Command: command,
		Args:    args,
		Started: time.Now(),
	}

	cmd := exec.CommandContext(ctx, command, args...)

	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	result.Finished = time.Now()
	result.Duration = result.Finished.Sub(result.Started)
	result.Stdout = stdoutBuf.String()
	result.Stderr = stderrBuf.String()

	if err != nil {
		result.Error = err
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}

	return result
}

// Success returns true when the command exited with code 0.
func (r Result) Success() bool {
	return r.ExitCode == 0 && r.Error == nil
}
