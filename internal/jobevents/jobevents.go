// Package jobevents provides a simple pub/sub event bus for job lifecycle
// events. Subscribers register for named events and receive payloads
// asynchronously via buffered channels.
package jobevents

import (
	"errors"
	"sync"
)

// EventType represents a named job lifecycle event.
type EventType string

const (
	EventStarted  EventType = "job.started"
	EventFinished EventType = "job.finished"
	EventFailed   EventType = "job.failed"
)

// Event carries the event type and associated metadata.
type Event struct {
	Type    EventType
	Job     string
	Payload map[string]any
}

// Bus is a thread-safe event bus.
type Bus struct {
	mu          sync.RWMutex
	subscribers map[EventType][]chan Event
	bufSize     int
}

// New creates a new Bus. bufSize controls the channel buffer per subscriber.
func New(bufSize int) (*Bus, error) {
	if bufSize < 1 {
		return nil, errors.New("jobevents: bufSize must be at least 1")
	}
	return &Bus{
		subscribers: make(map[EventType][]chan Event),
		bufSize:     bufSize,
	}, nil
}

// Subscribe registers a new subscriber for the given event type and returns
// a receive-only channel that will receive matching events.
func (b *Bus) Subscribe(et EventType) (<-chan Event, error) {
	if et == "" {
		return nil, errors.New("jobevents: event type must not be empty")
	}
	ch := make(chan Event, b.bufSize)
	b.mu.Lock()
	b.subscribers[et] = append(b.subscribers[et], ch)
	b.mu.Unlock()
	return ch, nil
}

// Publish sends an event to all subscribers registered for its type.
// Sends are non-blocking; if a subscriber's buffer is full the event is dropped.
func (b *Bus) Publish(e Event) error {
	if e.Type == "" {
		return errors.New("jobevents: event type must not be empty")
	}
	if e.Job == "" {
		return errors.New("jobevents: job name must not be empty")
	}
	b.mu.RLock()
	chans := b.subscribers[e.Type]
	b.mu.RUnlock()
	for _, ch := range chans {
		select {
		case ch <- e:
		default:
		}
	}
	return nil
}

// Close closes all subscriber channels, signalling that no more events will
// be sent.
func (b *Bus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for et, chans := range b.subscribers {
		for _, ch := range chans {
			close(ch)
		}
		delete(b.subscribers, et)
	}
}
