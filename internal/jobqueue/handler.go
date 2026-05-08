package jobqueue

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler returns an HTTP handler that exposes job queue state.
// GET /queue        — returns all entries currently in the queue
// GET /queue/{name} — returns the queue entry for a specific job
func Handler(q *Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Trim leading slash and split on remaining slashes.
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		// If the last segment is non-empty it is treated as a job name filter.
		name := ""
		if len(parts) > 1 {
			name = parts[len(parts)-1]
		}

		w.Header().Set("Content-Type", "application/json")

		if name != "" {
			entry, ok := q.Peek(name)
			if !ok {
				http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(entry)
			return
		}

		entries := q.All()
		if entries == nil {
			entries = []Entry{}
		}
		json.NewEncoder(w).Encode(entries)
	}
}
