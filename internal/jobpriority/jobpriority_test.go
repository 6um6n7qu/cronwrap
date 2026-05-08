package jobpriority

import (
	"testing"
)

func TestPush_EmptyNameReturnsError(t *testing.T) {
	q := New()
	if err := q.Push("", "echo hi", Normal); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestPush_EmptyCommandReturnsError(t *testing.T) {
	q := New()
	if err := q.Push("myjob", "", Normal); err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestPop_EmptyQueueReturnsNil(t *testing.T) {
	q := New()
	if e := q.Pop(); e != nil {
		t.Fatalf("expected nil, got %v", e)
	}
}

func TestPop_ReturnsHighestPriorityFirst(t *testing.T) {
	q := New()
	_ = q.Push("low", "echo low", Low)
	_ = q.Push("high", "echo high", High)
	_ = q.Push("normal", "echo normal", Normal)

	e := q.Pop()
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if e.Name != "high" {
		t.Fatalf("expected high priority job first, got %q", e.Name)
	}
}

func TestPop_SamePriorityFIFOOrder(t *testing.T) {
	q := New()
	_ = q.Push("first", "echo first", Normal)
	_ = q.Push("second", "echo second", Normal)

	e := q.Pop()
	if e.Name != "first" {
		t.Fatalf("expected FIFO order, got %q first", e.Name)
	}
}

func TestLen_TracksQueueSize(t *testing.T) {
	q := New()
	if q.Len() != 0 {
		t.Fatalf("expected 0, got %d", q.Len())
	}
	_ = q.Push("a", "echo a", Low)
	_ = q.Push("b", "echo b", High)
	if q.Len() != 2 {
		t.Fatalf("expected 2, got %d", q.Len())
	}
	q.Pop()
	if q.Len() != 1 {
		t.Fatalf("expected 1, got %d", q.Len())
	}
}

func TestPop_DrainAllEntries(t *testing.T) {
	q := New()
	names := []string{"c", "a", "b"}
	priorities := []Priority{Low, High, Normal}
	for i, n := range names {
		_ = q.Push(n, "echo "+n, priorities[i])
	}
	order := []string{}
	for q.Len() > 0 {
		e := q.Pop()
		order = append(order, e.Name)
	}
	if order[0] != "a" || order[1] != "b" || order[2] != "c" {
		t.Fatalf("unexpected drain order: %v", order)
	}
}
