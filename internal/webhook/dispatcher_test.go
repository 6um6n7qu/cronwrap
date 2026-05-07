package webhook_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/cronwrap/internal/webhook"
)

func captureServer(t *testing.T) (*httptest.Server, *[]webhook.Payload) {
	t.Helper()
	var payloads []webhook.Payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p webhook.Payload
		_ = json.NewDecoder(r.Body).Decode(&p)
		payloads = append(payloads, p)
		w.WriteHeader(http.StatusOK)
	}))
	return srv, &payloads
}

func TestDispatcher_JobStarted(t *testing.T) {
	srv, payloads := captureServer(t)
	defer srv.Close()

	c, _ := webhook.New(srv.URL, "", nil)
	d := webhook.NewDispatcher(c)
	d.JobStarted(context.Background(), "cleanup")

	if len(*payloads) != 1 || (*payloads)[0].Event != "started" {
		t.Fatalf("expected started event, got %+v", *payloads)
	}
}

func TestDispatcher_JobFinished(t *testing.T) {
	srv, payloads := captureServer(t)
	defer srv.Close()

	c, _ := webhook.New(srv.URL, "", nil)
	d := webhook.NewDispatcher(c)
	d.JobFinished(context.Background(), "backup", 0, 2*time.Second)

	if len(*payloads) != 1 {
		t.Fatalf("expected 1 payload, got %d", len(*payloads))
	}
	p := (*payloads)[0]
	if p.Event != "finished" {
		t.Errorf("expected event=finished, got %q", p.Event)
	}
	if p.ExitCode != 0 {
		t.Errorf("expected exit_code=0, got %d", p.ExitCode)
	}
	if p.Duration != 2000 {
		t.Errorf("expected duration_ms=2000, got %d", p.Duration)
	}
}

func TestDispatcher_JobFailed(t *testing.T) {
	srv, payloads := captureServer(t)
	defer srv.Close()

	c, _ := webhook.New(srv.URL, "", nil)
	d := webhook.NewDispatcher(c)
	d.JobFailed(context.Background(), "report", 1, "exit status 1")

	if len(*payloads) != 1 {
		t.Fatalf("expected 1 payload, got %d", len(*payloads))
	}
	p := (*payloads)[0]
	if p.Event != "failed" || p.Error == "" {
		t.Errorf("unexpected payload: %+v", p)
	}
}

func TestDispatcher_ErrorIsSilent(t *testing.T) {
	// Point at a server that returns 500; Dispatcher should not panic or return.
	srv := newTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := webhook.New(srv.URL, "", nil)
	d := webhook.NewDispatcher(c)
	// Should complete without panic.
	d.JobStarted(context.Background(), "myjob")
}
