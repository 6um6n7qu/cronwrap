package jobtag

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler returns an http.HandlerFunc that exposes job tag data over HTTP.
//
// Routes:
//
//	GET  /jobtags            — list all jobs and their tags
//	GET  /jobtags/{name}     — list tags for a single job
//	GET  /jobtags?tag={tag}  — list jobs carrying a specific tag
func Handler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// ?tag=<value> — jobs carrying the tag
		if tag := r.URL.Query().Get("tag"); tag != "" {
			jobs := s.JobsWithTag(tag)
			if jobs == nil {
				jobs = []string{}
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"tag": tag, "jobs": jobs})
			return
		}

		// /jobtags/{name}
		name := strings.TrimPrefix(r.URL.Path, "/jobtags/")
		if name != "" && name != r.URL.Path {
			tags, err := s.Tags(name)
			if err == ErrJobNotFound {
				http.Error(w, "job not found", http.StatusNotFound)
				return
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if tags == nil {
				tags = []string{}
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"job": name, "tags": tags})
			return
		}

		// list all
		s.mu.RLock()
		result := make(map[string][]string, len(s.jobs))
		for job, tagSet := range s.jobs {
			list := make([]string, 0, len(tagSet))
			for t := range tagSet {
				list = append(list, t)
			}
			result[job] = list
		}
		s.mu.RUnlock()
		_ = json.NewEncoder(w).Encode(result)
	}
}
