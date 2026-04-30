package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cronwrap-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
job:
  name: backup
  command: /usr/bin/rsync
  args: ["-av", "/src", "/dst"]
  timeout: 30m
notifier:
  smtp_host: smtp.example.com
  smtp_port: 587
  from: alerts@example.com
  recipients:
    - ops@example.com
logger:
  level: info
  output: stdout
`)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Job.Name != "backup" {
		t.Errorf("expected job name 'backup', got %q", cfg.Job.Name)
	}
	if cfg.Job.Timeout != 30*time.Minute {
		t.Errorf("expected timeout 30m, got %v", cfg.Job.Timeout)
	}
	if len(cfg.Notifier.Recipients) != 1 {
		t.Errorf("expected 1 recipient, got %d", len(cfg.Notifier.Recipients))
	}
}

func TestLoad_MissingJobName(t *testing.T) {
	path := writeTemp(t, `
job:
  command: /bin/true
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing job name")
	}
}

func TestLoad_MissingCommand(t *testing.T) {
	path := writeTemp(t, `
job:
  name: myjob
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing command")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
