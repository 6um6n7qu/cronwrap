package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/circuitbreaker"
)

func TestNew_DefaultsApplied(t *testing.T) {
	cb := circuitbreaker.New(0, 0)
	if cb == nil {
		t.Fatal("expected non-nil CircuitBreaker")
	}
	if cb.CurrentState() != circuitbreaker.StateClosed {
		t.Errorf("expected initial state Closed, got %v", cb.CurrentState())
	}
}

func TestAllow_ClosedStatePermitsCall(t *testing.T) {
	cb := circuitbreaker.New(3, time.Second)
	if err := cb.Allow(); err != nil {
		t.Errorf("expected nil error in closed state, got %v", err)
	}
}

func TestRecordFailure_OpensCircuitAtThreshold(t *testing.T) {
	cb := circuitbreaker.New(3, time.Minute)

	cb.RecordFailure()
	cb.RecordFailure()
	if cb.CurrentState() != circuitbreaker.StateClosed {
		t.Error("expected circuit to remain closed before threshold")
	}

	cb.RecordFailure()
	if cb.CurrentState() != circuitbreaker.StateOpen {
		t.Errorf("expected circuit to be open after threshold, got %v", cb.CurrentState())
	}
}

func TestAllow_OpenStateRejectsCall(t *testing.T) {
	cb := circuitbreaker.New(1, time.Minute)
	cb.RecordFailure()

	if err := cb.Allow(); err != circuitbreaker.ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestAllow_TransitionsToHalfOpenAfterTimeout(t *testing.T) {
	cb := circuitbreaker.New(1, 50*time.Millisecond)
	cb.RecordFailure()

	time.Sleep(60 * time.Millisecond)

	if err := cb.Allow(); err != nil {
		t.Errorf("expected nil error after reset timeout, got %v", err)
	}
	if cb.CurrentState() != circuitbreaker.StateHalfOpen {
		t.Errorf("expected HalfOpen state, got %v", cb.CurrentState())
	}
}

func TestRecordSuccess_ClosesCircuit(t *testing.T) {
	cb := circuitbreaker.New(1, 50*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(60 * time.Millisecond)
	_ = cb.Allow() // transitions to half-open

	cb.RecordSuccess()
	if cb.CurrentState() != circuitbreaker.StateClosed {
		t.Errorf("expected Closed after success, got %v", cb.CurrentState())
	}
	if cb.Failures() != 0 {
		t.Errorf("expected 0 failures after success, got %d", cb.Failures())
	}
}

func TestFailures_TracksCount(t *testing.T) {
	cb := circuitbreaker.New(5, time.Minute)
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.Failures() != 2 {
		t.Errorf("expected 2 failures, got %d", cb.Failures())
	}
}
