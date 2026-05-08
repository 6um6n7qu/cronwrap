// Package jobqueue implements a lightweight, bounded FIFO queue for pending
// cron job executions within cronwrap.
//
// When the scheduler fires a job trigger faster than the runner can consume
// them — for example during a burst of missed executions after a restart —
// the queue absorbs the backlog up to a configurable capacity. Entries that
// arrive when the queue is full are rejected with ErrQueueFull so the caller
// can decide whether to log, alert, or drop the execution.
//
// # Typical usage
//
//	q, err := jobqueue.New(64)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Producer (scheduler goroutine)
//	if err := q.Enqueue(jobqueue.Entry{Name: "report", Command: "./gen-report.sh"}); err != nil {
//		log.Printf("warn: dropping job: %v", err)
//	}
//
//	// Consumer (worker goroutine)
//	for entry := range q.Dequeue() {
//		runner.Run(ctx, entry.Command, entry.Args...)
//	}
//
// Call Close when the process is shutting down to drain the channel and
// prevent further enqueues.
package jobqueue
