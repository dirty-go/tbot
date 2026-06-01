package tbot

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("tok", "")
	if c.baseURL != apiBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, apiBaseURL)
	}
	if c.url != "https://api.telegram.org/bottok%s" {
		t.Errorf("url template = %q", c.url)
	}
	if c.maxRetries != defaultMaxRetries {
		t.Errorf("maxRetries = %d, want %d", c.maxRetries, defaultMaxRetries)
	}
	if c.rateLimiter == nil {
		t.Error("expected a default rate limiter")
	}
	if c.logger == nil {
		t.Error("expected a non-nil default logger")
	}
	if c.httpClient == nil {
		t.Error("expected a default http client")
	}
}

func TestNewClient_WithBaseURLOverridesTemplate(t *testing.T) {
	c := NewClient("tok", "https://example.test")
	if c.url != "https://example.test/bottok%s" {
		t.Errorf("url template = %q", c.url)
	}
	c2 := NewClient("tok", "", WithBaseURL("https://opt.test"))
	if c2.url != "https://opt.test/bottok%s" {
		t.Errorf("WithBaseURL template = %q", c2.url)
	}
}

func TestClient_Me(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/getMe") {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		fmt.Fprint(w, `{"ok":true,"result":{"id":7,"is_bot":true,"first_name":"botto","username":"botto_bot"}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	me, err := c.Me(context.Background())
	if err != nil {
		t.Fatalf("Me() error: %v", err)
	}
	if me.ID != 7 || !me.IsBot || me.Username != "botto_bot" {
		t.Fatalf("unexpected user: %+v", me)
	}
}
