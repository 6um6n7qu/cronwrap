package healthcheck

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordSuccess_MarksHealthy(t *testing.T) {
	c := New()
	c.RecordSuccess("backup")

	c.mu.RLock()
	defer c.mu.RUnlock()

	s, ok := c.jobs["backup"]
	if !ok {
		t.Fatal("expected job 'backup' to exist")
	}
	if !s.Healthy {
		t.Error("expected job to be healthy after success")
	}
	if s.ConsecutiveFails != 0 {
		t.Errorf("expected 0 consecutive fails, got %d", s.ConsecutiveFails)
	}
}

func TestRecordFailure_MarksUnhealthy(t *testing.T) {
	c := New()
	c.RecordFailure("cleanup")
	c.RecordFailure("cleanup")

	c.mu.RLock()
	defer c.mu.RUnlock()

	s := c.jobs["cleanup"]
	if s.Healthy {
		t.Error("expected job to be unhealthy after failure")
	}
	if s.ConsecutiveFails != 2 {
		t.Errorf("expected 2 consecutive fails, got %d", s.ConsecutiveFails)
	}
}

func TestRecordSuccess_ResetsConsecutiveFails(t *testing.T) {
	c := New()
	c.RecordFailure("sync")
	c.RecordFailure("sync")
	c.RecordSuccess("sync")

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.jobs["sync"].ConsecutiveFails != 0 {
		t.Error("expected consecutive fails to reset after success")
	}
}

func TestHandler_HealthyReturns200(t *testing.T) {
	c := New()
	c.RecordSuccess("job1")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	c.Handler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var status Status
	if err := json.NewDecoder(rec.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !status.Healthy {
		t.Error("expected overall status to be healthy")
	}
}

func TestHandler_UnhealthyReturns503(t *testing.T) {
	c := New()
	c.RecordFailure("job1")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	c.Handler()(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rec.Code)
	}

	var status Status
	if err := json.NewDecoder(rec.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if status.Healthy {
		t.Error("expected overall status to be unhealthy")
	}
}

func TestHandler_EmptyJobsReturns200(t *testing.T) {
	c := New()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	c.Handler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 with no jobs, got %d", rec.Code)
	}
}
