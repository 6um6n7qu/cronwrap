package jobmeta

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler returns an http.HandlerFunc that exposes job metadata over HTTP.
//
// Routes:
//
//	GET  /jobmeta          – list all entries
//	GET  /jobmeta/{name}   – get a single entry
//	POST /jobmeta/{name}   – set a key/value pair (body: {"key":"…","value":"…"})
//	DELETE /jobmeta/{name}/{key} – remove a key
func Handler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Strip leading slash and split path segments.
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		// parts[0] == "jobmeta"

		switch {
		case r.Method == http.MethodGet && len(parts) == 1:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(s.All())

		case r.Method == http.MethodGet && len(parts) == 2:
			name := parts[1]
			e, ok := s.Get(name)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(e)

		case r.Method == http.MethodPost && len(parts) == 2:
			name := parts[1]
			var body struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return
			}
			if err := s.Set(name, body.Key, body.Value); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		case r.Method == http.MethodDelete && len(parts) == 3:
			name, key := parts[1], parts[2]
			if err := s.Delete(name, key); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
