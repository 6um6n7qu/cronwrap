package scheduler

import (
	"testing"

	"github.com/your-org/cronwrap/internal/config"
)

func TestParseJobsFromConfig_ReturnsAllJobs(t *testing.T) {
	cfg := &config.Config{
		Jobs: []config.Job{
			{Name: "alpha", Command: "echo alpha", Schedule: "@every 1m"},
			{Name: "beta", Command: "echo beta", Schedule: "@every 5m"},
		},
	}

	jobs := ParseJobsFromConfig(cfg)
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}
}

func TestParseJobsFromConfig_EmptyJobs(t *testing.T) {
	cfg := &config.Config{Jobs: []config.Job{}}
	jobs := ParseJobsFromConfig(cfg)
	if len(jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestParseJobsFromConfig_PreservesFields(t *testing.T) {
	cfg := &config.Config{
		Jobs: []config.Job{
			{
				Name:     "daily-report",
				Command:  "/usr/local/bin/report.sh",
				Schedule: "0 9 * * *",
				Timeout:  "30s",
			},
		},
	}

	jobs := ParseJobsFromConfig(cfg)
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}

	j := jobs[0]
	if j.Name != "daily-report" {
		t.Errorf("expected name 'daily-report', got %q", j.Name)
	}
	if j.Command != "/usr/local/bin/report.sh" {
		t.Errorf("expected command '/usr/local/bin/report.sh', got %q", j.Command)
	}
	if j.Schedule != "0 9 * * *" {
		t.Errorf("expected schedule '0 9 * * *', got %q", j.Schedule)
	}
	if j.Timeout != "30s" {
		t.Errorf("expected timeout '30s', got %q", j.Timeout)
	}
}

func TestParseJobsFromConfig_NilConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ParseJobsFromConfig panicked on nil config: %v", r)
		}
	}()

	cfg := &config.Config{}
	jobs := ParseJobsFromConfig(cfg)
	if jobs == nil {
		t.Error("expected non-nil slice, got nil")
	}
}
