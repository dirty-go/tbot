package tbot

import (
	"fmt"
	"net/http"
	"time"
)

type ClientOptions func(*Client)

// WithBaseURL overrides the Telegram Bot API base URL. Useful for
// pointing at a local bot API server or a test fake.
func WithBaseURL(baseURL string) ClientOptions {
	return func(c *Client) {
		if baseURL == "" {
			return
		}
		c.baseURL = baseURL
		c.url = fmt.Sprintf("%s/bot%s", baseURL, c.token) + "%s"
	}
}

// WithRateLimit installs a token-bucket rate limiter that admits up to
// rps requests per second with the given burst. Pass rps <= 0 or
// burst <= 0 to disable throttling entirely.
func WithRateLimit(rps float64, burst int) ClientOptions {
	return func(c *Client) {
		c.rateLimiter = newRateLimiter(rps, burst)
	}
}

// WithMaxRetries sets the maximum number of times a request will be
// retried after a 429 Too Many Requests, a 5xx response, or a transport
// error. Pass 0 to disable retries.
func WithMaxRetries(n int) ClientOptions {
	return func(c *Client) {
		if n < 0 {
			n = 0
		}
		c.maxRetries = n
	}
}

// WithLogger swaps the default no-op logger.
func WithLogger(l Logger) ClientOptions {
	return func(c *Client) {
		if l == nil {
			return
		}
		c.logger = l
	}
}

// WithHTTPClient lets callers supply a pre-configured *http.Client (for
// example to inject a custom Transport or timeout).
func WithHTTPClient(hc *http.Client) ClientOptions {
	return func(c *Client) {
		if hc == nil {
			return
		}
		c.httpClient = hc
	}
}

// WithHTTPTimeout overrides only the request timeout on the default
// http.Client. Ignored if WithHTTPClient was supplied.
func WithHTTPTimeout(d time.Duration) ClientOptions {
	return func(c *Client) {
		if d <= 0 || c.httpClient == nil {
			return
		}
		c.httpClient.Timeout = d
	}
}
