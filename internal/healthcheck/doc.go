// Package healthcheck provides a lightweight health-check mechanism for
// cronwrap jobs.
//
// A Checker instance tracks the last known result (success or failure) for
// each named job. It exposes an http.HandlerFunc that returns a JSON payload
// describing the overall health of the process and the per-job status.
//
// HTTP response codes:
//   - 200 OK                  — all tracked jobs are healthy (or none exist)
//   - 503 Service Unavailable — at least one job has a recorded failure
//
// JSON response format:
//
//	{
//	  "healthy": true,
//	  "jobs": {
//	    "my-job": {
//	      "healthy": true,
//	      "last_run": "2024-01-15T10:30:00Z",
//	      "message": ""
//	    }
//	  }
//	}
//
// Typical usage:
//
//	hc := healthcheck.New()
//	http.Handle("/healthz", hc.Handler())
//	// later, after running a job:
//	hc.RecordSuccess("my-job")
//	// or, on failure:
//	hc.RecordFailure("my-job", "exit status 1")
package healthcheck
