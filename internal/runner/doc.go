// Package runner provides a thin wrapper around os/exec for executing
// external commands within cronwrap.
//
// It captures stdout, stderr, exit code, and wall-clock duration into a
// single Result value that the rest of the application can inspect for
// structured logging and failure alerting.
//
// Basic usage:
//
//	ctx := context.Background()
//	result := runner.Run(ctx, "my-job", "--flag", "value")
//	if !result.Success() {
//		// handle failure
//	}
package runner

import "strings"
