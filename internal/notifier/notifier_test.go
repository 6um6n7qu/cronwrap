package notifier_test

import (
	"net/smtp"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/notifier"
)

func TestNotify_SendsEmail(t *testing.T) {
	var capturedAddr string
	var capturedFrom string
	var capturedTo []string
	var capturedMsg []byte

	n := notifier.New(notifier.Config{
		SMTPHost: "localhost",
		SMTPPort: 25,
		From:     "cronwrap@example.com",
		To:       []string{"ops@example.com"},
		JobName:  "backup",
	})
	n.SetSendFunc(func(addr string, _ smtp.Auth, from string, to []string, msg []byte) error {
		capturedAddr = addr
		capturedFrom = from
		capturedTo = to
		capturedMsg = msg
		return nil
	})

	now := time.Now()
	alert := notifier.Alert{
		JobName:    "backup",
		Command:    "/usr/bin/backup.sh",
		ExitCode:   1,
		Stdout:     "backing up...",
		Stderr:     "error: disk full",
		StartedAt:  now,
		FinishedAt: now.Add(5 * time.Second),
	}

	if err := n.Notify(alert); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if capturedAddr != "localhost:25" {
		t.Errorf("expected addr localhost:25, got %s", capturedAddr)
	}
	if capturedFrom != "cronwrap@example.com" {
		t.Errorf("unexpected from: %s", capturedFrom)
	}
	if len(capturedTo) != 1 || capturedTo[0] != "ops@example.com" {
		t.Errorf("unexpected to: %v", capturedTo)
	}

	body := string(capturedMsg)
	for _, want := range []string{
		"Subject: [cronwrap] Job failed: backup",
		"Exit Code: 1",
		"/usr/bin/backup.sh",
		"error: disk full",
		"backing up...",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("expected body to contain %q", want)
		}
	}
}

func TestNotify_NoRecipients(t *testing.T) {
	n := notifier.New(notifier.Config{
		SMTPHost: "localhost",
		SMTPPort: 25,
		From:     "cronwrap@example.com",
		To:       []string{},
	})

	err := n.Notify(notifier.Alert{JobName: "test"})
	if err == nil {
		t.Fatal("expected error for empty recipients, got nil")
	}
}
