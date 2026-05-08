package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestNew_CreatesEmptyStore(t *testing.T) {
	s, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(s.All()); got != 0 {
		t.Errorf("expected 0 entries, got %d", got)
	}
}

func TestRecord_StoresEntry(t *testing.T) {
	s, _ := New(tempPath(t))
	e := Entry{
		JobName:    "backup",
		StartedAt:  time.Now(),
		FinishedAt: time.Now(),
		Duration:   2 * time.Second,
		ExitCode:   0,
		Success:    true,
	}
	if err := s.Record(e); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
	got, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.ExitCode != 0 || !got.Success {
		t.Errorf("unexpected entry: %+v", got)
	}
}

func TestRecord_EmptyNameReturnsError(t *testing.T) {
	s, _ := New(tempPath(t))
	err := s.Record(Entry{})
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestGet_UnknownJobReturnsFalse(t *testing.T) {
	s, _ := New(tempPath(t))
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestPersistAndReload(t *testing.T) {
	path := tempPath(t)
	s1, _ := New(path)
	_ = s1.Record(Entry{JobName: "cleanup", ExitCode: 1, Success: false})

	s2, err := New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	e, ok := s2.Get("cleanup")
	if !ok {
		t.Fatal("expected reloaded entry")
	}
	if e.ExitCode != 1 || e.Success {
		t.Errorf("unexpected reloaded entry: %+v", e)
	}
}

func TestNew_FileNotFoundIsOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	_, err := New(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s, _ := New(tempPath(t))
	_ = s.Record(Entry{JobName: "job1", Success: true})
	all := s.All()
	all["job1"] = Entry{JobName: "job1", Success: false}
	original, _ := s.Get("job1")
	if !original.Success {
		t.Error("All() should return a copy, not a reference")
	}
	_ = os.Remove(s.path)
}
