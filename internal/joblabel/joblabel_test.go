package joblabel_test

import (
	"testing"

	"github.com/yourorg/cronwrap/internal/joblabel"
)

func TestRegister_AddsJob(t *testing.T) {
	s := joblabel.New()
	if err := s.Register("backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels, err := s.Get("backup")
	if err != nil {
		t.Fatalf("expected labels, got error: %v", err)
	}
	if len(labels) != 0 {
		t.Errorf("expected empty labels, got %v", labels)
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	s := joblabel.New()
	if err := s.Register(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestSet_AssignsLabel(t *testing.T) {
	s := joblabel.New()
	_ = s.Register("sync")
	if err := s.Set("sync", "env", "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels, _ := s.Get("sync")
	if labels["env"] != "prod" {
		t.Errorf("expected 'prod', got %q", labels["env"])
	}
}

func TestSet_UnknownJobReturnsError(t *testing.T) {
	s := joblabel.New()
	if err := s.Set("ghost", "k", "v"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestSet_EmptyKeyReturnsError(t *testing.T) {
	s := joblabel.New()
	_ = s.Register("job")
	if err := s.Set("job", "", "val"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestDelete_RemovesLabel(t *testing.T) {
	s := joblabel.New()
	_ = s.Register("job")
	_ = s.Set("job", "team", "ops")
	_ = s.Delete("job", "team")
	labels, _ := s.Get("job")
	if _, ok := labels["team"]; ok {
		t.Error("expected label to be deleted")
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	s := joblabel.New()
	_ = s.Register("job")
	_ = s.Set("job", "region", "us-east")
	labels, _ := s.Get("job")
	labels["region"] = "mutated"
	orig, _ := s.Get("job")
	if orig["region"] != "us-east" {
		t.Errorf("store was mutated via returned map")
	}
}

func TestMatch_ReturnsMatchingJobs(t *testing.T) {
	s := joblabel.New()
	_ = s.Register("job-a")
	_ = s.Register("job-b")
	_ = s.Register("job-c")
	_ = s.Set("job-a", "env", "prod")
	_ = s.Set("job-a", "team", "ops")
	_ = s.Set("job-b", "env", "prod")
	_ = s.Set("job-c", "env", "staging")

	matches := s.Match(joblabel.Labels{"env": "prod"})
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(matches), matches)
	}
}

func TestMatch_EmptySelectorMatchesAll(t *testing.T) {
	s := joblabel.New()
	_ = s.Register("a")
	_ = s.Register("b")
	matches := s.Match(joblabel.Labels{})
	if len(matches) != 2 {
		t.Errorf("expected 2, got %d", len(matches))
	}
}
