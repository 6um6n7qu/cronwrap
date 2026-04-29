// Package notifier implements failure alerting for cronwrap.
//
// When a cron job exits with a non-zero status code, the notifier
// composes and delivers an email alert containing the job name,
// command, exit code, duration, and captured stdout/stderr output.
//
// Basic usage:
//
//	n := notifier.New(notifier.Config{
//		SMTPHost: "smtp.example.com",
//		SMTPPort: 587,
//		From:     "cronwrap@example.com",
//		To:       []string{"ops@example.com"},
//		JobName:  "nightly-backup",
//	})
//
//	err := n.Notify(notifier.Alert{
//		JobName:    "nightly-backup",
//		Command:    "/usr/local/bin/backup.sh",
//		ExitCode:   2,
//		Stderr:     "fatal: cannot open archive",
//		StartedAt:  startTime,
//		FinishedAt: endTime,
//	})
package notifier
