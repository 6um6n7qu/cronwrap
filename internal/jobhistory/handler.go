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
//	GET /jobhistory        → all recorded entries (newest first)
//	GET /jobhistory/{name} → entries for a single named job
func Handler(h *History) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Trim leading slash and optional prefix so both
		// "/jobhistory" and "/jobhistory/my-job" work.
		path := strings.TrimPrefix(r.URL.Path, "/")
		parts := strings.SplitN(path, "/", 2)

		if len(parts) == 2 && parts[1] != "" {
			name := parts[1]
			entries, ok := h.Get(name)
			if !ok {
				http.Error(w, "job not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(entries)
			return
		}

		all := h.All()
		if all == nil {
			all = []Entry{}
		}
		json.NewEncoder(w).Encode(all)
	}
}
