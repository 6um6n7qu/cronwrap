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

func newTestServer(t *testing.T, statusCode int, capture *webhook.Payload) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if capture != nil {
			_ = json.NewDecoder(r.Body).Decode(capture)
		}
		w.WriteHeader(statusCode)
	}))
}

func TestNew_EmptyURLReturnsError(t *testing.T) {
	_, err := webhook.New("", "", nil)
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestSend_PostsPayload(t *testing.T) {
	var received webhook.Payload
	srv := newTestServer(t, http.StatusOK, &received)
	defer srv.Close()

	c, err := webhook.New(srv.URL, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := webhook.Payload{JobName: "backup", Event: "finished", ExitCode: 0}
	if err := c.Send(context.Background(), p); err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if received.JobName != "backup" {
		t.Errorf("expected job_name=backup, got %q", received.JobName)
	}
	if received.Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}
}

func TestSend_AddsSecretHeader(t *testing.T) {
	var gotSecret string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSecret = r.Header.Get("X-Webhook-Secret")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c, _ := webhook.New(srv.URL, "mysecret", nil)
	_ = c.Send(context.Background(), webhook.Payload{JobName: "j", Event: "started"})
	if gotSecret != "mysecret" {
		t.Errorf("expected secret header mysecret, got %q", gotSecret)
	}
}

func TestSend_Non2xxReturnsError(t *testing.T) {
	srv := newTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := webhook.New(srv.URL, "", nil)
	err := c.Send(context.Background(), webhook.Payload{JobName: "j", Event: "failed"})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestSend_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c, _ := webhook.New(srv.URL, "", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := c.Send(ctx, webhook.Payload{JobName: "j", Event: "started"}); err == nil {
		t.Fatal("expected context cancellation error")
	}
}
