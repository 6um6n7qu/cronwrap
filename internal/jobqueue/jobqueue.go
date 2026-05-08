// Package jobqueue provides a bounded, thread-safe queue for pending cron job
// executions. Jobs that cannot be dispatched immediately are held until a
// worker slot becomes available or the queue capacity is reached.
package jobqueue

import (
	"errors"
	"sync"
	"time"
)

// ErrQueueFull is returned when the queue has reached its maximum capacity.
var ErrQueueFull = errors.New("jobqueue: queue is full")

// ErrQueueClosed is returned when an operation is attempted on a closed queue.
var ErrQueueClosed = errors.New("jobqueue: queue is closed")

// Entry represents a single pending job execution.
type Entry struct {
	Name      string
	Command   string
	Args      []string
	EnqueuedAt time.Time
}

// Queue is a bounded channel-backed job queue.
type Queue struct {
	mu     sync.Mutex
	ch     chan Entry
	closed bool
}

// New creates a new Queue with the given capacity.
// capacity must be greater than zero.
func New(capacity int) (*Queue, error) {
	if capacity <= 0 {
		return nil, errors.New("jobqueue: capacity must be greater than zero")
	}
	return &Queue{
		ch: make(chan Entry, capacity),
	}, nil
}

// Enqueue adds an entry to the queue. It returns ErrQueueFull if the queue is
// at capacity and ErrQueueClosed if the queue has been closed.
func (q *Queue) Enqueue(e Entry) error {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return ErrQueueClosed
	}
	q.mu.Unlock()

	if e.EnqueuedAt.IsZero() {
		e.EnqueuedAt = time.Now()
	}

	select {
	case q.ch <- e:
		return nil
	default:
		return ErrQueueFull
	}
}

// Dequeue returns the channel from which entries can be consumed.
// The channel is closed when Close is called.
func (q *Queue) Dequeue() <-chan Entry {
	return q.ch
}

// Len returns the current number of items waiting in the queue.
func (q *Queue) Len() int {
	return len(q.ch)
}

// Close drains and closes the queue. Subsequent Enqueue calls will return
// ErrQueueClosed.
func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if !q.closed {
		q.closed = true
		close(q.ch)
	}
}
