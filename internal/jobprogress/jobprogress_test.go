package jobprogress

import (
	"testing"
)

func TestRegister_AddsEntry(t *testing.T) {
	s := New()
	if err := s.Register("backup", 100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Total != 100 {
		t.Errorf("expected total 100, got %d", e.Total)
	}
	if e.Done != 0 {
		t.Errorf("expected done 0, got %d", e.Done)
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	s := New()
	if err := s.Register("", 10); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_ZeroTotalReturnsError(t *testing.T) {
	s := New()
	if err := s.Register("job", 0); err == nil {
		t.Fatal("expected error for zero total")
	}
}

func TestUpdate_SetsProgressAndPercent(t *testing.T) {
	s := New()
	_ = s.Register("export", 200)
	if err := s.Update("export", 50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := s.Get("export")
	if e.Done != 50 {
		t.Errorf("expected done 50, got %d", e.Done)
	}
	if e.Percent != 25.0 {
		t.Errorf("expected percent 25.0, got %f", e.Percent)
	}
}

func TestUpdate_ClampsToTotal(t *testing.T) {
	s := New()
	_ = s.Register("sync", 10)
	_ = s.Update("sync", 999)
	e, _ := s.Get("sync")
	if e.Done != 10 {
		t.Errorf("expected done clamped to 10, got %d", e.Done)
	}
	if e.Percent != 100.0 {
		t.Errorf("expected percent 100.0, got %f", e.Percent)
	}
}

func TestUpdate_UnknownJobReturnsError(t *testing.T) {
	s := New()
	if err := s.Update("ghost", 5); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGet_UnknownJobReturnsFalse(t *testing.T) {
	s := New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected ok=false for unknown job")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s := New()
	_ = s.Register("a", 10)
	_ = s.Register("b", 20)
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
