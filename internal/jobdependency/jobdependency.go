package jobdependency

import (
	"errors"
	"fmt"
	"sync"
)

// ErrCyclicDependency is returned when adding a job would create a cycle.
var ErrCyclicDependency = errors.New("cyclic dependency detected")

// ErrUnknownJob is returned when referencing a job that has not been registered.
var ErrUnknownJob = errors.New("unknown job")

// Graph tracks directed dependencies between named cron jobs.
// A job may only run after all of its declared dependencies have completed.
type Graph struct {
	mu   sync.RWMutex
	edges map[string][]string // job -> dependencies
}

// New returns an empty dependency Graph.
func New() *Graph {
	return &Graph{edges: make(map[string][]string)}
}

// Register adds a job to the graph with no dependencies.
// Registering an already-known job is a no-op.
func (g *Graph) Register(name string) error {
	if name == "" {
		return errors.New("job name must not be empty")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.edges[name]; !ok {
		g.edges[name] = []string{}
	}
	return nil
}

// AddDependency declares that job depends on dep.
// Both jobs must be registered first. Returns ErrCyclicDependency if
// the relationship would introduce a cycle.
func (g *Graph) AddDependency(job, dep string) error {
	if job == "" || dep == "" {
		return errors.New("job and dep names must not be empty")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.edges[job]; !ok {
		return fmt.Errorf("%w: %s", ErrUnknownJob, job)
	}
	if _, ok := g.edges[dep]; !ok {
		return fmt.Errorf("%w: %s", ErrUnknownJob, dep)
	}
	// Temporarily add the edge and check for cycles.
	g.edges[job] = append(g.edges[job], dep)
	if g.hasCycle() {
		// Roll back.
		deps := g.edges[job]
		g.edges[job] = deps[:len(deps)-1]
		return ErrCyclicDependency
	}
	return nil
}

// Dependencies returns the direct dependencies of job.
func (g *Graph) Dependencies(job string) ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	deps, ok := g.edges[job]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownJob, job)
	}
	out := make([]string, len(deps))
	copy(out, deps)
	return out, nil
}

// Order returns a topological ordering of all registered jobs.
// Jobs with no dependencies appear first.
func (g *Graph) Order() ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.topoSort()
}

// hasCycle performs a DFS cycle check. Caller must hold the write lock.
func (g *Graph) hasCycle() bool {
	visited := make(map[string]bool)
	onStack := make(map[string]bool)
	var dfs func(n string) bool
	dfs = func(n string) bool {
		visited[n] = true
		onStack[n] = true
		for _, dep := range g.edges[n] {
			if !visited[dep] && dfs(dep) {
				return true
			} else if onStack[dep] {
				return true
			}
		}
		onStack[n] = false
		return false
	}
	for n := range g.edges {
		if !visited[n] && dfs(n) {
			return true
		}
	}
	return false
}

// topoSort returns a stable topological ordering. Caller must hold the read lock.
func (g *Graph) topoSort() ([]string, error) {
	visited := make(map[string]bool)
	var order []string
	var visit func(n string)
	visit = func(n string) {
		if visited[n] {
			return
		}
		visited[n] = true
		for _, dep := range g.edges[n] {
			visit(dep)
		}
		order = append(order, n)
	}
	for n := range g.edges {
		visit(n)
	}
	// Reverse so dependencies come first.
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return order, nil
}
