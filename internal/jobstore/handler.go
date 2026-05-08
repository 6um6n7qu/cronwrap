package jobstore

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler returns an http.Handler that exposes job store state over HTTP.
//
// GET /jobs        — returns all entries as a JSON array.
// GET /jobs/{name} — returns a single entry or 404.
func Handler(s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Trim leading slash and optional "jobs" prefix.
		path := strings.TrimPrefix(r.URL.Path, "/")
		path = strings.TrimPrefix(path, "jobs")
		path = strings.TrimPrefix(path, "/")

		w.Header().Set("Content-Type", "application/json")

		if path == "" {
			_ = json.NewEncoder(w).Encode(s.All())
			return
		}

		entry, ok := s.Get(path)
		if !ok {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(entry)
	})
}
