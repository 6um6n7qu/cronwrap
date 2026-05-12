package jobnotify

import (
	"testing"
)

func TestRegister_AddsRule(t *testing.T) {
	s := New()
	err := s.Register(Rule{JobName: "backup", On: []Event{EventFailed}, Channels: []string{"email"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	chans, ok := s.ShouldNotify("backup", EventFailed)
	if !ok {
		t.Fatal("expected notification to be triggered")
	}
	if len(chans) != 1 || chans[0] != "email" {
		t.Fatalf("unexpected channels: %v", chans)
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	s := New()
	err := s.Register(Rule{JobName: "", On: []Event{EventFailed}, Channels: []string{"slack"}})
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestRegister_NoEventsReturnsError(t *testing.T) {
	s := New()
	err := s.Register(Rule{JobName: "sync", On: []Event{}, Channels: []string{"slack"}})
	if err == nil {
		t.Fatal("expected error when no events specified")
	}
}

func TestRegister_NoChannelsReturnsError(t *testing.T) {
	s := New()
	err := s.Register(Rule{JobName: "sync", On: []Event{EventStarted}, Channels: []string{}})
	if err == nil {
		t.Fatal("expected error when no channels specified")
	}
}

func TestShouldNotify_UnknownJobReturnsFalse(t *testing.T) {
	s := New()
	_, ok := s.ShouldNotify("ghost", EventFailed)
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestShouldNotify_EventNotInRule(t *testing.T) {
	s := New()
	_ = s.Register(Rule{JobName: "report", On: []Event{EventFailed}, Channels: []string{"pager"}})
	_, ok := s.ShouldNotify("report", EventStarted)
	if ok {
		t.Fatal("expected false when event not in rule")
	}
}

func TestRecord_AppendsHistory(t *testing.T) {
	s := New()
	s.Record("backup", EventFailed, "email")
	s.Record("backup", EventFinished, "slack")
	h := s.History()
	if len(h) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(h))
	}
	if h[0].Event != EventFailed {
		t.Errorf("unexpected first event: %v", h[0].Event)
	}
}

func TestRemove_DeletesRule(t *testing.T) {
	s := New()
	_ = s.Register(Rule{JobName: "cleanup", On: []Event{EventFailed}, Channels: []string{"email"}})
	s.Remove("cleanup")
	_, ok := s.ShouldNotify("cleanup", EventFailed)
	if ok {
		t.Fatal("expected rule to be removed")
	}
}
