package jobstatus

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler returns an http.HandlerFunc that exposes job status over HTTP.
//
// GET /jobstatus        — returns all job status entries as a JSON array.
// GET /jobstatus/{name} — returns the status entry for a single job.
func Handler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(r.URL.Path, "/")
		w.Header().Set("Content-Type", "application/json")

		if name == "" {
			entries := s.All()
			if entries == nil {
				entries = []Entry{}
			}
			_ = json.NewEncoder(w).Encode(entries)
			return
		}

		e, ok := s.Get(name)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(e)
	}
}
