package metrics_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/metrics"
)

func TestRecord_SuccessIncrements(t *testing.T) {
	c := metrics.New()
	c.Record("backup", 2*time.Second, true)

	m, ok := c.Get("backup")
	if !ok {
		t.Fatal("expected metrics to exist for 'backup'")
	}
	if m.TotalRuns != 1 {
		t.Errorf("expected TotalRuns=1, got %d", m.TotalRuns)
	}
	if m.SuccessRuns != 1 {
		t.Errorf("expected SuccessRuns=1, got %d", m.SuccessRuns)
	}
	if m.FailureRuns != 0 {
		t.Errorf("expected FailureRuns=0, got %d", m.FailureRuns)
	}
}

func TestRecord_FailureIncrements(t *testing.T) {
	c := metrics.New()
	c.Record("cleanup", 500*time.Millisecond, false)

	m, _ := c.Get("cleanup")
	if m.FailureRuns != 1 {
		t.Errorf("expected FailureRuns=1, got %d", m.FailureRuns)
	}
}

func TestRecord_AverageDuration(t *testing.T) {
	c := metrics.New()
	c.Record("job", 2*time.Second, true)
	c.Record("job", 4*time.Second, true)

	m, _ := c.Get("job")
	if m.TotalRuns != 2 {
		t.Errorf("expected TotalRuns=2, got %d", m.TotalRuns)
	}
	want := 3 * time.Second
	if m.AvgDuration != want {
		t.Errorf("expected AvgDuration=%v, got %v", want, m.AvgDuration)
	}
}

func TestGet_UnknownJob(t *testing.T) {
	c := metrics.New()
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected ok=false for unknown job")
	}
}

func TestWriteJSON_CreatesValidFile(t *testing.T) {
	c := metrics.New()
	c.Record("export", time.Second, true)
	c.Record("export", 3*time.Second, false)

	tmp, err := os.CreateTemp(t.TempDir(), "metrics-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmp.Close()

	if err := c.WriteJSON(tmp.Name()); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("failed to read metrics file: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
	if result[0]["job_name"] != "export" {
		t.Errorf("unexpected job_name: %v", result[0]["job_name"])
	}
}
