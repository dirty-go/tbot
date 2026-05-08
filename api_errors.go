package tbot

import "fmt"

// APIError describes a non-success response from the Telegram Bot API.
// It is returned whenever the bot API replies with ok=false or an HTTP
// error status. RetryAfter is populated from the response body's
// parameters.retry_after (or the Retry-After header on a 429) when
// present, so callers can implement higher-level back-pressure on top of
// the client's built-in retries.
type APIError struct {
	StatusCode  int
	ErrorCode   int
	Description string
	RetryAfter  int
}

func (e *APIError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("telegram api: %d %s (retry after %ds)", e.ErrorCode, e.Description, e.RetryAfter)
	}
	if e.ErrorCode != 0 {
		return fmt.Sprintf("telegram api: %d %s", e.ErrorCode, e.Description)
	}
	return fmt.Sprintf("telegram api: http %d %s", e.StatusCode, e.Description)
}
