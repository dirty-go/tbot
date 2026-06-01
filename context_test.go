package tbot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestDo_CancelledContextAbortsBeforeRequest(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		fmt.Fprint(w, `{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"x"}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := c.sendRequest(ctx, "/getMe", nil, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if got := hits.Load(); got != 0 {
		t.Fatalf("expected no HTTP request, got %d hits", got)
	}
}

func TestBackoffSleep_HonoursContext(t *testing.T) {
	// A fresh client has no sleepFn override, so backoffSleep uses the
	// context-aware timer path. A cancelled context must abort immediately
	// rather than sleeping the full hour.
	c := NewClient("tok", "https://example.test")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)
	go func() { done <- c.backoffSleep(ctx, time.Hour) }()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("backoffSleep did not honour context cancellation")
	}
}
