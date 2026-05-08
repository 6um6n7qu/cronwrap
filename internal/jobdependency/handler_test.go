package jobdependency

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func buildGraph(t *testing.T) *Graph {
	t.Helper()
	g := New()
	for _, name := range []string{"build", "test", "deploy"} {
		if err := g.Register(name); err != nil {
			t.Fatalf("register %s: %v", name, err)
		}
	}
	_ = g.AddDependency("test", "build")
	_ = g.AddDependency("deploy", "test")
	return g
}

func TestHandler_OrderReturns200(t *testing.T) {
	g := buildGraph(t)
	h := Handler(g)
	req := httptest.NewRequest(http.MethodGet, "/jobdependency", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp orderResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp.Order) != 3 {
		t.Fatalf("expected 3 jobs in order, got %v", resp.Order)
	}
}

func TestHandler_SingleJobFound(t *testing.T) {
	g := buildGraph(t)
	h := Handler(g)
	req := httptest.NewRequest(http.MethodGet, "/jobdependency/test", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp dependencyResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Job != "test" {
		t.Errorf("expected job=test, got %s", resp.Job)
	}
	if len(resp.Dependencies) != 1 || resp.Dependencies[0] != "build" {
		t.Errorf("expected [build], got %v", resp.Dependencies)
	}
}

func TestHandler_SingleJobNotFound(t *testing.T) {
	g := buildGraph(t)
	h := Handler(g)
	req := httptest.NewRequest(http.MethodGet, "/jobdependency/ghost", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	g := buildGraph(t)
	h := Handler(g)
	req := httptest.NewRequest(http.MethodPost, "/jobdependency", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
