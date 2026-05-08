// Package jobpriority provides a priority queue for scheduling cron jobs
// with different urgency levels. Jobs with higher priority are dequeued first.
package jobpriority

import (
	"container/heap"
	"errors"
	"sync"
	"time"
)

// Priority represents the urgency level of a job.
type Priority int

const (
	Low    Priority = 1
	Normal Priority = 5
	High   Priority = 10
)

// Entry holds a job scheduled for execution with an associated priority.
type Entry struct {
	Name      string
	Command   string
	Priority  Priority
	EnqueuedAt time.Time
	index     int
}

// Queue is a thread-safe priority queue for job entries.
type Queue struct {
	mu   sync.Mutex
	items priorityHeap
}

// New creates an empty Queue.
func New() *Queue {
	q := &Queue{}
	heap.Init(&q.items)
	return q
}

// Push adds a job entry to the queue.
func (q *Queue) Push(name, command string, p Priority) error {
	if name == "" {
		return errors.New("jobpriority: name must not be empty")
	}
	if command == "" {
		return errors.New("jobpriority: command must not be empty")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	e := &Entry{
		Name:       name,
		Command:    command,
		Priority:   p,
		EnqueuedAt: time.Now(),
	}
	heap.Push(&q.items, e)
	return nil
}

// Pop removes and returns the highest-priority entry.
// Returns nil if the queue is empty.
func (q *Queue) Pop() *Entry {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.items.Len() == 0 {
		return nil
	}
	return heap.Pop(&q.items).(*Entry)
}

// Len returns the number of entries currently in the queue.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.items.Len()
}

// priorityHeap implements heap.Interface.
type priorityHeap []*Entry

func (h priorityHeap) Len() int { return len(h) }
func (h priorityHeap) Less(i, j int) bool {
	if h[i].Priority != h[j].Priority {
		return h[i].Priority > h[j].Priority
	}
	return h[i].EnqueuedAt.Before(h[j].EnqueuedAt)
}
func (h priorityHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *priorityHeap) Push(x any) {
	e := x.(*Entry)
	e.index = len(*h)
	*h = append(*h, e)
}
func (h *priorityHeap) Pop() any {
	old := *h
	n := len(old)
	e := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	return e
}
