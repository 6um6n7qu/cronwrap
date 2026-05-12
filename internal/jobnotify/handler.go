package jobnotify

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler returns an http.HandlerFunc that exposes notification rules and
// dispatch history over HTTP.
//
// GET  /jobnotify          — list all dispatched history entries
// GET  /jobnotify/{job}    — list history entries for a specific job
func Handler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Trim leading slash and split path segments.
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		// segments[0] == "jobnotify"
		jobName := ""
		if len(segments) >= 2 {
			jobName = segments[1]
		}

		all := s.History()

		if jobName == "" {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(all); err != nil {
				http.Error(w, "encoding error", http.StatusInternalServerError)
			}
			return
		}

		var filtered []Entry
		for _, e := range all {
			if e.JobName == jobName {
				filtered = append(filtered, e)
			}
		}
		if filtered == nil {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(filtered); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	}
}
