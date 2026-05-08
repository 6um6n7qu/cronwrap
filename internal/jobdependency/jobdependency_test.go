package jobdependency

import (
	"errors"
	"testing"
)

func TestRegister_AddsJob(t *testing.T) {
	g := New()
	if err := g.Register("job-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps, err := g.Dependencies("job-a")
	if err != nil {
		t.Fatalf("expected job-a to exist: %v", err)
	}
	if len(deps) != 0 {
		t.Fatalf("expected no deps, got %v", deps)
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	g := New()
	if err := g.Register(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAddDependency_LinksTwoJobs(t *testing.T) {
	g := New()
	_ = g.Register("a")
	_ = g.Register("b")
	if err := g.AddDependency("a", "b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps, _ := g.Dependencies("a")
	if len(deps) != 1 || deps[0] != "b" {
		t.Fatalf("expected [b], got %v", deps)
	}
}

func TestAddDependency_UnknownJobReturnsError(t *testing.T) {
	g := New()
	_ = g.Register("a")
	err := g.AddDependency("a", "ghost")
	if !errors.Is(err, ErrUnknownJob) {
		t.Fatalf("expected ErrUnknownJob, got %v", err)
	}
}

func TestAddDependency_DetectsCycle(t *testing.T) {
	g := New()
	_ = g.Register("a")
	_ = g.Register("b")
	_ = g.AddDependency("a", "b")
	err := g.AddDependency("b", "a")
	if !errors.Is(err, ErrCyclicDependency) {
		t.Fatalf("expected ErrCyclicDependency, got %v", err)
	}
}

func TestAddDependency_SelfCycleDetected(t *testing.T) {
	g := New()
	_ = g.Register("a")
	err := g.AddDependency("a", "a")
	if !errors.Is(err, ErrCyclicDependency) {
		t.Fatalf("expected ErrCyclicDependency, got %v", err)
	}
}

func TestOrder_TopologicalOrdering(t *testing.T) {
	g := New()
	_ = g.Register("build")
	_ = g.Register("test")
	_ = g.Register("deploy")
	_ = g.AddDependency("test", "build")
	_ = g.AddDependency("deploy", "test")

	order, err := g.Order()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pos := func(name string) int {
		for i, n := range order {
			if n == name {
				return i
			}
		}
		return -1
	}
	if pos("build") >= pos("test") {
		t.Errorf("expected build before test in %v", order)
	}
	if pos("test") >= pos("deploy") {
		t.Errorf("expected test before deploy in %v", order)
	}
}

func TestDependencies_UnknownJobReturnsError(t *testing.T) {
	g := New()
	_, err := g.Dependencies("ghost")
	if !errors.Is(err, ErrUnknownJob) {
		t.Fatalf("expected ErrUnknownJob, got %v", err)
	}
}
