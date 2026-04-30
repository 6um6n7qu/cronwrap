// Package logger provides structured JSON logging for cronwrap cron job execution.
//
// Each log entry is written as a single line of JSON to the configured writer
// (defaulting to os.Stdout), making output easy to ingest by log aggregators
// such as Datadog, Splunk, or CloudWatch Logs.
//
// Example usage:
//
//	l := logger.New("backup-job", nil)
//	l.JobStarted()
//	// ... run job ...
//	l.JobFinished(0, elapsed)
//
// Output format:
//
//	{"timestamp":"2024-01-15T10:00:00Z","level":"INFO","job":"backup-job","message":"job started"}
//	{"timestamp":"2024-01-15T10:00:05Z","level":"INFO","job":"backup-job","message":"job finished","duration":"5s","exit_code":0}
package logger
