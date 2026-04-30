package scheduler

import (
	"fmt"
	"strings"
)

// ParseJobsFromConfig converts a slice of raw config maps into typed Job
// values, returning all validation errors collected during parsing.
func ParseJobsFromConfig(raw []map[string]string) ([]Job, error) {
	var (
		jobs []Job
		errs []string
	)

	for i, entry := range raw {
		j := Job{
			Name:     strings.TrimSpace(entry["name"]),
			Schedule: strings.TrimSpace(entry["schedule"]),
			Command:  strings.TrimSpace(entry["command"]),
		}

		if args, ok := entry["args"]; ok && strings.TrimSpace(args) != "" {
			for _, a := range strings.Fields(args) {
				j.Args = append(j.Args, a)
			}
		}

		if j.Name == "" {
			errs = append(errs, fmt.Sprintf("entry %d: missing name", i))
			continue
		}
		if j.Schedule == "" {
			errs = append(errs, fmt.Sprintf("entry %d (%s): missing schedule", i, j.Name))
			continue
		}
		if j.Command == "" {
			errs = append(errs, fmt.Sprintf("entry %d (%s): missing command", i, j.Name))
			continue
		}

		jobs = append(jobs, j)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("scheduler: parse errors:\n  %s", strings.Join(errs, "\n  "))
	}
	return jobs, nil
}
