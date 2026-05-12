package jobmeta

import (
	"testing"
)

func TestRegister_AddsJob(t *testing.T) {
	s := New()
	if err := s.Register("backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected job to be registered")
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	s := New()
	if err := s.Register(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_DuplicateReturnsError(t *testing.T) {
	s := New()
	_ = s.Register("backup")
	if err := s.Register("backup"); err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestSet_AssignsMeta(t *testing.T) {
	s := New()
	_ = s.Register("sync")
	if err := s.Set("sync", "owner", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := s.Get("sync")
	if e.Meta["owner"] != "alice" {
		t.Errorf("expected alice, got %s", e.Meta["owner"])
	}
}

func TestSet_UnknownJobReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("ghost", "k", "v"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestSet_EmptyKeyReturnsError(t *testing.T) {
	s := New()
	_ = s.Register("sync")
	if err := s.Set("sync", "", "v"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestDelete_RemovesKey(t *testing.T) {
	s := New()
	_ = s.Register("sync")
	_ = s.Set("sync", "env", "prod")
	if err := s.Delete("sync", "env"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := s.Get("sync")
	if _, ok := e.Meta["env"]; ok {
		t.Error("expected key to be deleted")
	}
}

func TestDelete_UnknownJobReturnsError(t *testing.T) {
	s := New()
	if err := s.Delete("ghost", "k"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Register("job")
	_ = s.Set("job", "x", "1")
	e, _ := s.Get("job")
	e.Meta["x"] = "mutated"
	e2, _ := s.Get("job")
	if e2.Meta["x"] != "1" {
		t.Error("Get should return an independent copy")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s := New()
	_ = s.Register("a")
	_ = s.Register("b")
	if len(s.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(s.All()))
	}
}
