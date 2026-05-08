package jobqueue

import (
	"testing"
	"time"
)

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero capacity, got nil")
	}
	_, err = New(-5)
	if err == nil {
		t.Fatal("expected error for negative capacity, got nil")
	}
}

func TestEnqueue_AddsEntry(t *testing.T) {
	q, _ := New(4)
	err := q.Enqueue(Entry{Name: "backup", Command: "tar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Len() != 1 {
		t.Fatalf("expected len 1, got %d", q.Len())
	}
}

func TestEnqueue_SetsEnqueuedAt(t *testing.T) {
	q, _ := New(4)
	before := time.Now()
	_ = q.Enqueue(Entry{Name: "ping", Command: "ping"})
	after := time.Now()

	e := <-q.Dequeue()
	if e.EnqueuedAt.Before(before) || e.EnqueuedAt.After(after) {
		t.Errorf("EnqueuedAt %v not within expected range", e.EnqueuedAt)
	}
}

func TestEnqueue_ReturnsErrQueueFull(t *testing.T) {
	q, _ := New(2)
	_ = q.Enqueue(Entry{Name: "a"})
	_ = q.Enqueue(Entry{Name: "b"})

	err := q.Enqueue(Entry{Name: "c"})
	if err != ErrQueueFull {
		t.Fatalf("expected ErrQueueFull, got %v", err)
	}
}

func TestDequeue_ReceivesEntries(t *testing.T) {
	q, _ := New(4)
	_ = q.Enqueue(Entry{Name: "job1", Command: "echo", Args: []string{"hello"}})

	select {
	case e := <-q.Dequeue():
		if e.Name != "job1" {
			t.Errorf("expected job1, got %s", e.Name)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for dequeue")
	}
}

func TestClose_PreventsEnqueue(t *testing.T) {
	q, _ := New(4)
	q.Close()

	err := q.Enqueue(Entry{Name: "late"})
	if err != ErrQueueClosed {
		t.Fatalf("expected ErrQueueClosed, got %v", err)
	}
}

func TestClose_ClosesDequeueChannel(t *testing.T) {
	q, _ := New(4)
	_ = q.Enqueue(Entry{Name: "x"})
	q.Close()

	count := 0
	for range q.Dequeue() {
		count++
	}
	if count != 1 {
		t.Errorf("expected 1 drained entry, got %d", count)
	}
}

func TestClose_IdempotentDoubleClose(t *testing.T) {
	q, _ := New(4)
	q.Close()
	q.Close() // must not panic
}
