package jobhistory

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler returns an http.HandlerFunc that exposes job history over HTTP.
//
// Routes:
//
//	GET /jobhistory          — returns all entries
//	GET /jobhistory/{name}   — returns entries for a specific job
func Handler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Trim leading slash and the "jobhistory" prefix to extract optional name.
		path := strings.TrimPrefix(r.URL.Path, "/")
		parts := strings.SplitN(path, "/", 2)

		var entries []Entry
		if len(parts) == 2 && parts[1] != "" {
			name := parts[1]
			entries = s.ForJob(name)
			if len(entries) == 0 {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
		} else {
			entries = s.All()
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(entries); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	}
}
