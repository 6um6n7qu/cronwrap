package jobmeta

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func buildStore(t *testing.T) *Store {
	t.Helper()
	s := New()
	if err := s.Register("deploy"); err != nil {
		t.Fatal(err)
	}
	_ = s.Set("deploy", "env", "staging")
	return s
}

func TestHandler_AllReturns200(t *testing.T) {
	s := buildStore(t)
	req := httptest.NewRequest(http.MethodGet, "/jobmeta", nil)
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestHandler_SingleJobFound(t *testing.T) {
	s := buildStore(t)
	req := httptest.NewRequest(http.MethodGet, "/jobmeta/deploy", nil)
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var e Entry
	if err := json.NewDecoder(rr.Body).Decode(&e); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if e.Meta["env"] != "staging" {
		t.Errorf("expected staging, got %s", e.Meta["env"])
	}
}

func TestHandler_SingleJobNotFound(t *testing.T) {
	s := buildStore(t)
	req := httptest.NewRequest(http.MethodGet, "/jobmeta/ghost", nil)
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestHandler_PostSetsMeta(t *testing.T) {
	s := buildStore(t)
	body := `{"key":"region","value":"us-east-1"}`
	req := httptest.NewRequest(http.MethodPost, "/jobmeta/deploy", strings.NewReader(body))
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
	e, _ := s.Get("deploy")
	if e.Meta["region"] != "us-east-1" {
		t.Errorf("expected us-east-1, got %s", e.Meta["region"])
	}
}

func TestHandler_PostInvalidBodyReturns400(t *testing.T) {
	s := buildStore(t)
	req := httptest.NewRequest(http.MethodPost, "/jobmeta/deploy", bytes.NewBufferString("not-json"))
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHandler_DeleteRemovesKey(t *testing.T) {
	s := buildStore(t)
	req := httptest.NewRequest(http.MethodDelete, "/jobmeta/deploy/env", nil)
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
	e, _ := s.Get("deploy")
	if _, ok := e.Meta["env"]; ok {
		t.Error("expected key to be removed")
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	s := buildStore(t)
	req := httptest.NewRequest(http.MethodPut, "/jobmeta", nil)
	rr := httptest.NewRecorder()
	Handler(s)(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
