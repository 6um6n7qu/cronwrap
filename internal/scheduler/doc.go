// Package scheduler provides a thin wrapper around a cron scheduler,
// allowing cronwrap to register named jobs with validated cron expressions
// and manage their lifecycle via a context.
//
// Usage:
//
//	s := scheduler.New()
//	err := s.Add(scheduler.Job{
//		Name:     "cleanup",
//		Schedule: "0 2 * * *",
//		Command:  "/usr/local/bin/cleanup",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	s.Start(ctx)
package scheduler
