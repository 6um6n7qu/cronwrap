package jobevents_test

import (
	"testing"
	"time"

	"github.com/cronwrap/internal/jobevents"
)

func TestNew_InvalidBufSize(t *testing.T) {
	_, err := jobevents.New(0)
	if err == nil {
		t.Fatal("expected error for bufSize=0")
	}
}

func TestSubscribe_EmptyEventTypeReturnsError(t *testing.T) {
	b, _ := jobevents.New(4)
	_, err := b.Subscribe("")
	if err == nil {
		t.Fatal("expected error for empty event type")
	}
}

func TestPublish_EmptyEventTypeReturnsError(t *testing.T) {
	b, _ := jobevents.New(4)
	err := b.Publish(jobevents.Event{Job: "myjob"})
	if err == nil {
		t.Fatal("expected error for empty event type")
	}
}

func TestPublish_EmptyJobNameReturnsError(t *testing.T) {
	b, _ := jobevents.New(4)
	err := b.Publish(jobevents.Event{Type: jobevents.EventStarted})
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestPublish_DeliverToSubscriber(t *testing.T) {
	b, _ := jobevents.New(4)
	ch, err := b.Subscribe(jobevents.EventStarted)
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	evt := jobevents.Event{
		Type:    jobevents.EventStarted,
		Job:     "backup",
		Payload: map[string]any{"attempt": 1},
	}
	if err := b.Publish(evt); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	select {
	case got := <-ch:
		if got.Job != "backup" {
			t.Errorf("expected job=backup, got %s", got.Job)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestPublish_NoSubscribersIsNoop(t *testing.T) {
	b, _ := jobevents.New(4)
	err := b.Publish(jobevents.Event{Type: jobevents.EventFailed, Job: "cleanup"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPublish_DropWhenBufferFull(t *testing.T) {
	b, _ := jobevents.New(1)
	ch, _ := b.Subscribe(jobevents.EventFinished)

	// Fill the buffer.
	_ = b.Publish(jobevents.Event{Type: jobevents.EventFinished, Job: "job1"})
	// This should be dropped silently.
	_ = b.Publish(jobevents.Event{Type: jobevents.EventFinished, Job: "job2"})

	got := <-ch
	if got.Job != "job1" {
		t.Errorf("expected job1, got %s", got.Job)
	}
}

func TestClose_ClosesChannels(t *testing.T) {
	b, _ := jobevents.New(4)
	ch, _ := b.Subscribe(jobevents.EventStarted)
	b.Close()

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for close")
	}
}
