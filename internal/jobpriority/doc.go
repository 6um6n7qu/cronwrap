// Package jobpriority implements a thread-safe priority queue for cronwrap
// job scheduling.
//
// Jobs are assigned a Priority level (Low, Normal, or High) when enqueued.
// The queue always returns the highest-priority job first. When two jobs
// share the same priority, they are returned in FIFO order based on
// enqueue time.
//
// # Usage
//
//	q := jobpriority.New()
//
//	_ = q.Push("backup",  "pg_dump mydb",  jobpriority.High)
//	_ = q.Push("report",  "gen-report.sh", jobpriority.Normal)
//	_ = q.Push("cleanup", "rm -rf /tmp/*", jobpriority.Low)
//
//	for q.Len() > 0 {
//		entry := q.Pop()
//		// dispatch entry.Command
//	}
package jobpriority
