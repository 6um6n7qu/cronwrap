// Package metrics provides lightweight job execution tracking for cronwrap.
//
// It records per-job success/failure counts and duration statistics,
// persisting results to a JSON file for external monitoring tools.
//
// Basic usage:
//
//	m := metrics.New("/var/lib/cronwrap/metrics.json")
//	m.Record("backup-job", true, 4*time.Second)
//	if err := m.WriteJSON(); err != nil {
//		log.Fatal(err)
//	}
//
// Retrieve stats for a specific job:
//
//	stats, err := m.Get("backup-job")
//	if err != nil {
//		log.Println("unknown job")
//	}
//	fmt.Printf("success rate: %d/%d\n", stats.Successes, stats.Successes+stats.Failures)
package metrics
