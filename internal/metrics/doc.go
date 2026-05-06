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
//
// Load previously persisted metrics from disk:
//
//	m := metrics.New("/var/lib/cronwrap/metrics.json")
//	if err := m.ReadJSON(); err != nil {
//		log.Printf("could not load metrics: %v", err)
//	}
//
// A typical workflow combining load, record, and persist:
//
//	m := metrics.New("/var/lib/cronwrap/metrics.json")
//	_ = m.ReadJSON() // ignore error if file does not yet exist
//	m.Record("backup-job", success, duration)
//	if err := m.WriteJSON(); err != nil {
//		log.Printf("could not save metrics: %v", err)
//	}
package metrics
