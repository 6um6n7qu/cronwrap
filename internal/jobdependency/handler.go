package jobdependency

import (
	"encoding/json"
	"net/http"
	"strings"
)

type dependencyResponse struct {
	Job          string   `json:"job"`
	Dependencies []string `json:"dependencies"`
}

type orderResponse struct {
	Order []string `json:"order"`
}

// Handler returns an http.HandlerFunc that exposes the dependency graph over HTTP.
//
//	GET /jobdependency          -> topological order of all jobs
//	GET /jobdependency/{name}   -> direct dependencies of the named job
func Handler(g *Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Trim leading slash and split on remaining slashes.
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		// Last segment is the job name when provided.
		name := ""
		if len(segments) > 1 {
			name = segments[len(segments)-1]
		}

		w.Header().Set("Content-Type", "application/json")

		if name == "" {
			order, err := g.Order()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.NewEncoder(w).Encode(orderResponse{Order: order})
			return
		}

		deps, err := g.Dependencies(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(dependencyResponse{Job: name, Dependencies: deps})
	}
}
