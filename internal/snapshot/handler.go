package snapshot

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.Handler that exposes the snapshot store
// as a JSON endpoint.
//
// GET /snapshots        — returns all entries
// GET /snapshots/{job}  — returns a single entry (404 if not found)
func Handler(s *Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/snapshots", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(s.All())
	})

	mux.HandleFunc("/snapshots/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		jobName := r.URL.Path[len("/snapshots/"):]
		if jobName == "" {
			http.Error(w, "job name required", http.StatusBadRequest)
			return
		}
		e, ok := s.Get(jobName)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(e)
	})

	return mux
}
