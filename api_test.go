package tbot

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// newTestClient builds a *Client wired to the given test server. Rate
// limiting and real sleeps are disabled so tests run instantly. Each
// caller gets a counter recording the cumulative sleep that the retry
// loop *would* have performed.
func newTestClient(t *testing.T, srv *httptest.Server) (*Client, *atomic.Int64) {
	t.Helper()
	var slept atomic.Int64
	c := NewClient("test-token", srv.URL,
		WithRateLimit(0, 0),
		WithMaxRetries(3),
	)
	c.sleepFn = func(d time.Duration) {
		slept.Add(int64(d))
	}
	return c, &slept
}

func TestSendRequest_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/getMe") {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		fmt.Fprint(w, `{"ok":true,"result":{"id":42,"is_bot":true,"first_name":"bot"}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	var u User
	if err := c.sendRequest(context.Background(), "/getMe", nil, &u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID != 42 || !u.IsBot || u.FirstName != "bot" {
		t.Fatalf("unexpected user: %+v", u)
	}
}

func TestSendRequest_RetriesOn429ThenSucceeds(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := hits.Add(1)
		if n <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, `{"ok":false,"error_code":429,"description":"Too Many Requests: retry after 2","parameters":{"retry_after":2}}`)
			return
		}
		fmt.Fprint(w, `{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"x"}}`)
	}))
	defer srv.Close()

	c, slept := newTestClient(t, srv)
	var u User
	if err := c.sendRequest(context.Background(), "/getMe", nil, &u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := hits.Load(); got != 3 {
		t.Fatalf("expected 3 hits, got %d", got)
	}
	// two retries × 2-second retry_after each.
	if got := time.Duration(slept.Load()); got < 4*time.Second {
		t.Fatalf("expected slept >= 4s, got %v", got)
	}
}

func TestSendRequest_HonoursRetryAfterHeader(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hits.Add(1) == 1 {
			w.Header().Set("Retry-After", "5")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, `not even json`)
			return
		}
		fmt.Fprint(w, `{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"x"}}`)
	}))
	defer srv.Close()

	c, slept := newTestClient(t, srv)
	var u User
	if err := c.sendRequest(context.Background(), "/getMe", nil, &u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := time.Duration(slept.Load()); got < 5*time.Second {
		t.Fatalf("expected at least one 5s sleep, got %v", got)
	}
}

func TestSendRequest_429Exhausted_ReturnsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, `{"ok":false,"error_code":429,"description":"Too Many Requests: retry after 1","parameters":{"retry_after":1}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	err := c.sendRequest(context.Background(), "/getMe", nil, nil)
	if err == nil {
		t.Fatal("expected error after exhausted retries")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests || apiErr.RetryAfter != 1 {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}
}

func TestSendRequest_200OkFalse429IsRetried(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := hits.Add(1)
		if n == 1 {
			fmt.Fprint(w, `{"ok":false,"error_code":429,"description":"flood","parameters":{"retry_after":1}}`)
			return
		}
		fmt.Fprint(w, `{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"x"}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	var u User
	if err := c.sendRequest(context.Background(), "/getMe", nil, &u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hits.Load() != 2 {
		t.Fatalf("expected 2 hits, got %d", hits.Load())
	}
}

func TestSendRequest_5xxRetriedThenFailsAsAPIError(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, `bad gateway`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	err := c.sendRequest(context.Background(), "/getMe", nil, nil)
	if err == nil {
		t.Fatal("expected error after exhausted retries")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
	if got := hits.Load(); got != 4 { // 1 initial + 3 retries
		t.Fatalf("expected 4 hits, got %d", got)
	}
}

func TestSendRequest_4xxNotRetried(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"ok":false,"error_code":400,"description":"Bad Request"}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	err := c.sendRequest(context.Background(), "/getMe", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if hits.Load() != 1 {
		t.Fatalf("expected single hit, got %d", hits.Load())
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.ErrorCode != 400 {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendRequest_FormBodyEchoed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Errorf("unexpected content type %q", ct)
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "chat_id=42") {
			t.Errorf("expected chat_id in form body, got %q", body)
		}
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"date":0,"chat":{"id":42,"type":"private"}}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	v := url.Values{}
	v.Set("chat_id", "42")
	v.Set("text", "hi")
	var msg Message
	if err := c.sendRequest(context.Background(), "/sendMessage", v, &msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRateLimiter_ThrottlesConcurrentRequests(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"x"}}`)
	}))
	defer srv.Close()

	c := NewClient("test-token", srv.URL,
		WithRateLimit(1000, 2), // burst 2 → after 2 calls everyone sleeps
		WithMaxRetries(0),
	)

	var slept atomic.Int64
	c.sleepFn = func(d time.Duration) { slept.Add(int64(d)) }
	// also stub the limiter's own sleep so the test doesn't actually pause
	if c.rateLimiter != nil {
		c.rateLimiter.sleep = func(d time.Duration) { slept.Add(int64(d)) }
	}

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var u User
			_ = c.sendRequest(context.Background(), "/getMe", nil, &u)
		}()
	}
	wg.Wait()
	// After consuming the initial 2-token burst, the remaining 8
	// requests should each wait — each token at 1000/s costs ~1ms.
	if slept.Load() == 0 {
		t.Fatal("expected limiter to enforce a wait once the burst was consumed")
	}
}

func TestBackoffDuration_IsCappedAndMonotonic(t *testing.T) {
	prev := time.Duration(0)
	for i := range 10 {
		d := backoffDuration(i)
		if d > backoffMax {
			t.Fatalf("attempt %d: %v exceeds cap %v", i, d, backoffMax)
		}
		if d < prev && d != backoffMax {
			t.Fatalf("attempt %d: %v < prev %v (should be monotonic until cap)", i, d, prev)
		}
		prev = d
	}
	if backoffDuration(1000) != backoffMax {
		t.Fatalf("very large attempt should saturate at %v", backoffMax)
	}
	if backoffDuration(-1) != backoffBase {
		t.Fatalf("negative attempt should not panic and should fall back to base")
	}
}

func TestAPIError_Error(t *testing.T) {
	cases := []struct {
		err  *APIError
		want string
	}{
		{&APIError{StatusCode: 429, ErrorCode: 429, Description: "flood", RetryAfter: 5}, "telegram api: 429 flood (retry after 5s)"},
		{&APIError{StatusCode: 400, ErrorCode: 400, Description: "bad"}, "telegram api: 400 bad"},
		{&APIError{StatusCode: 502, Description: "Bad Gateway"}, "telegram api: http 502 Bad Gateway"},
	}
	for _, tc := range cases {
		if got := tc.err.Error(); got != tc.want {
			t.Errorf("Error() = %q, want %q", got, tc.want)
		}
	}
}
