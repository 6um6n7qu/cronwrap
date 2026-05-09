package joblabel

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler returns an http.Handler exposing label read and write endpoints.
//
// Routes:
//
//	GET  /labels/{job}           — retrieve all labels for a job
//	POST /labels/{job}           — set a label ({"key":"k","value":"v"})
//	DELETE /labels/{job}/{key}   — remove a single label
func Handler(s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip leading slash and split path segments.
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		// segments[0] == "labels", segments[1] == job name, segments[2] == key (optional)
		if len(segments) < 2 || segments[1] == "" {
			http.Error(w, "job name required", http.StatusBadRequest)
			return
		}
		name := segments[1]

		switch r.Method {
		case http.MethodGet:
			labels, err := s.Get(name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(labels)

		case http.MethodPost:
			var body struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON body", http.StatusBadRequest)
				return
			}
			if err := s.Set(name, body.Key, body.Value); err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			if len(segments) < 3 || segments[2] == "" {
				http.Error(w, "label key required", http.StatusBadRequest)
				return
			}
			key := segments[2]
			if err := s.Delete(name, key); err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
