package dedup_test

import (
	"sync"
	"testing"

	"github.com/example/cronwrap/internal/dedup"
)

func TestAcquire_SucceedsForNewJob(t *testing.T) {
	d := dedup.New()
	if err := d.Acquire("backup"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAcquire_ReturnsErrorForDuplicate(t *testing.T) {
	d := dedup.New()
	_ = d.Acquire("backup")
	err := d.Acquire("backup")
	if err != dedup.ErrAlreadyRunning {
		t.Fatalf("expected ErrAlreadyRunning, got %v", err)
	}
}

func TestAcquire_EmptyNameReturnsError(t *testing.T) {
	d := dedup.New()
	if err := d.Acquire(""); err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestRelease_AllowsReacquire(t *testing.T) {
	d := dedup.New()
	_ = d.Acquire("sync")
	d.Release("sync")
	if err := d.Acquire("sync"); err != nil {
		t.Fatalf("expected reacquire to succeed, got %v", err)
	}
}

func TestRelease_NoopForUnknownJob(t *testing.T) {
	d := dedup.New()
	// Should not panic or error.
	d.Release("nonexistent")
}

func TestIsRunning_ReflectsState(t *testing.T) {
	d := dedup.New()
	if d.IsRunning("report") {
		t.Fatal("expected false before acquire")
	}
	_ = d.Acquire("report")
	if !d.IsRunning("report") {
		t.Fatal("expected true after acquire")
	}
	d.Release("report")
	if d.IsRunning("report") {
		t.Fatal("expected false after release")
	}
}

func TestActiveJobs_ReturnsSnapshot(t *testing.T) {
	d := dedup.New()
	_ = d.Acquire("job-a")
	_ = d.Acquire("job-b")
	active := d.ActiveJobs()
	if len(active) != 2 {
		t.Fatalf("expected 2 active jobs, got %d", len(active))
	}
	if _, ok := active["job-a"]; !ok {
		t.Error("expected job-a in active jobs")
	}
}

func TestConcurrent_SafeUnderContention(t *testing.T) {
	d := dedup.New()
	const workers = 50
	var wg sync.WaitGroup
	acquired := make(chan struct{}, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := d.Acquire("contested"); err == nil {
				acquired <- struct{}{}
				d.Release("contested")
			}
		}()
	}
	wg.Wait()
	close(acquired)

	count := 0
	for range acquired {
		count++
	}
	if count == 0 {
		t.Fatal("expected at least one successful acquire")
	}
}
