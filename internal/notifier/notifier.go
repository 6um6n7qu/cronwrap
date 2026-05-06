// Package notifier provides alerting functionality for cronwrap.
// It sends failure notifications when cron jobs exit with non-zero status.
package notifier

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

// Config holds the configuration for the notifier.
type Config struct {
	// SMTPHost is the SMTP server hostname.
	SMTPHost string
	// SMTPPort is the SMTP server port.
	SMTPPort int
	// From is the sender email address.
	From string
	// To is a list of recipient email addresses.
	To []string
	// JobName is the name of the cron job for display in alerts.
	JobName string
}

// Alert represents a failure alert to be sent.
type Alert struct {
	JobName   string
	Command   string
	ExitCode  int
	Stdout    string
	Stderr    string
	StartedAt time.Time
	FinishedAt time.Time
}

// Notifier sends failure alerts via email.
type Notifier struct {
	cfg  Config
	send func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// New creates a new Notifier with the given configuration.
func New(cfg Config) *Notifier {
	return &Notifier{
		cfg:  cfg,
		send: smtp.SendMail,
	}
}

// Notify sends a failure alert email for the given Alert.
func (n *Notifier) Notify(alert Alert) error {
	if len(n.cfg.To) == 0 {
		return fmt.Errorf("notifier: no recipients configured")
	}
	if n.cfg.SMTPHost == "" {
		return fmt.Errorf("notifier: SMTP host is not configured")
	}
	if n.cfg.From == "" {
		return fmt.Errorf("notifier: sender address is not configured")
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s\r\n", n.cfg.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(n.cfg.To, ", ")))
	buf.WriteString(fmt.Sprintf("Subject: [cronwrap] Job failed: %s\r\n", alert.JobName))
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(fmt.Sprintf("Job '%s' failed.\n\n", alert.JobName))
	buf.WriteString(fmt.Sprintf("Command:   %s\n", alert.Command))
	buf.WriteString(fmt.Sprintf("Exit Code: %d\n", alert.ExitCode))
	buf.WriteString(fmt.Sprintf("Started:   %s\n", alert.StartedAt.Format(time.RFC3339)))
	buf.WriteString(fmt.Sprintf("Finished:  %s\n", alert.FinishedAt.Format(time.RFC3339)))
	buf.WriteString(fmt.Sprintf("Duration:  %s\n", alert.FinishedAt.Sub(alert.StartedAt)))

	if alert.Stdout != "" {
		buf.WriteString(fmt.Sprintf("\n--- stdout ---\n%s\n", alert.Stdout))
	}
	if alert.Stderr != "" {
		buf.WriteString(fmt.Sprintf("\n--- stderr ---\n%s\n", alert.Stderr))
	}

	addr := fmt.Sprintf("%s:%d", n.cfg.SMTPHost, n.cfg.SMTPPort)
	if err := n.send(addr, nil, n.cfg.From, n.cfg.To, buf.Bytes()); err != nil {
		return fmt.Errorf("notifier: failed to send alert for job %q: %w", alert.JobName, err)
	}
	return nil
}
