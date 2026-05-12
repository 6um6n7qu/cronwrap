package jobtag_test

import (
	"testing"

	"github.com/yourorg/cronwrap/internal/jobtag"
)

func TestRegister_AddsJob(t *testing.T) {
	s := jobtag.New()
	if err := s.Register("backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, err := s.Tags("backup")
	if err != nil {
		t.Fatalf("expected job to exist: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("expected no tags, got %v", tags)
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	s := jobtag.New()
	if err := s.Register(""); err != jobtag.ErrEmptyName {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}
}

func TestAdd_AssignsTag(t *testing.T) {
	s := jobtag.New()
	_ = s.Register("sync")
	if err := s.Add("sync", "critical"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, _ := s.Tags("sync")
	if len(tags) != 1 || tags[0] != "critical" {
		t.Fatalf("expected [critical], got %v", tags)
	}
}

func TestAdd_UnknownJobReturnsError(t *testing.T) {
	s := jobtag.New()
	if err := s.Add("ghost", "nightly"); err != jobtag.ErrJobNotFound {
		t.Fatalf("expected ErrJobNotFound, got %v", err)
	}
}

func TestAdd_EmptyTagReturnsError(t *testing.T) {
	s := jobtag.New()
	_ = s.Register("job")
	if err := s.Add("job", ""); err != jobtag.ErrEmptyTag {
		t.Fatalf("expected ErrEmptyTag, got %v", err)
	}
}

func TestRemove_DetachesTag(t *testing.T) {
	s := jobtag.New()
	_ = s.Register("report")
	_ = s.Add("report", "weekly")
	if err := s.Remove("report", "weekly"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, _ := s.Tags("report")
	if len(tags) != 0 {
		t.Fatalf("expected no tags after removal, got %v", tags)
	}
}

func TestRemove_UnknownJobReturnsError(t *testing.T) {
	s := jobtag.New()
	if err := s.Remove("missing", "tag"); err != jobtag.ErrJobNotFound {
		t.Fatalf("expected ErrJobNotFound, got %v", err)
	}
}

func TestJobsWithTag_ReturnsMatchingJobs(t *testing.T) {
	s := jobtag.New()
	_ = s.Register("a")
	_ = s.Register("b")
	_ = s.Register("c")
	_ = s.Add("a", "nightly")
	_ = s.Add("c", "nightly")
	jobs := s.JobsWithTag("nightly")
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d: %v", len(jobs), jobs)
	}
}

func TestJobsWithTag_NoMatch(t *testing.T) {
	s := jobtag.New()
	_ = s.Register("x")
	if jobs := s.JobsWithTag("unknown"); len(jobs) != 0 {
		t.Fatalf("expected empty result, got %v", jobs)
	}
}
