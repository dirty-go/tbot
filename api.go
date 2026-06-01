package tbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type apiResponse struct {
	OK          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result"`
	Description string              `json:"description"`
	ErrorCode   int                 `json:"error_code"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

var netTransport = &http.Transport{
	TLSHandshakeTimeout:   10 * time.Second,
	MaxIdleConns:          10,
	IdleConnTimeout:       30 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

const (
	backoffBase = 500 * time.Millisecond
	backoffMax  = 30 * time.Second
)

// fileField pairs a multipart form field name with the file to upload there.
type fileField struct {
	field string
	file  InputFile
}

// endpoint renders the absolute URL for a Bot API method (e.g. "/sendMessage").
func (c *Client) endpoint(method string) string {
	return fmt.Sprintf(c.url, method)
}

// sendRequest performs a form-encoded Bot API call and decodes the result into
// response (which may be nil to discard it).
func (c *Client) sendRequest(ctx context.Context, method string, request url.Values, response any) error {
	var body string
	if request != nil {
		body = request.Encode()
	}
	bodyFn := func() io.Reader {
		if body == "" {
			return nil
		}
		return strings.NewReader(body)
	}
	return c.do(ctx, c.endpoint(method), "application/x-www-form-urlencoded", bodyFn, response)
}

// sendMultipart performs a multipart/form-data Bot API call, streaming each
// file in files into the body. The whole body is buffered once so it can be
// replayed on retry; uploads sourced from one-shot readers are consumed here,
// before the first HTTP attempt.
func (c *Client) sendMultipart(ctx context.Context, method string, request url.Values, files []fileField, response any) error {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for key, values := range request {
		for _, value := range values {
			if err := mw.WriteField(key, value); err != nil {
				return err
			}
		}
	}
	for _, f := range files {
		if err := writeUpload(mw, f); err != nil {
			return err
		}
	}
	if err := mw.Close(); err != nil {
		return err
	}

	contentType := mw.FormDataContentType()
	payload := buf.Bytes()
	bodyFn := func() io.Reader { return bytes.NewReader(payload) }
	return c.do(ctx, c.endpoint(method), contentType, bodyFn, response)
}

func writeUpload(mw *multipart.Writer, f fileField) error {
	rc, name, err := f.file.open()
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()
	fw, err := mw.CreateFormFile(f.field, name)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, rc)
	return err
}

// do executes a single Telegram Bot API call with rate limiting and retry on
// 429 / 5xx / transport errors. bodyFn must produce a fresh io.Reader on each
// call so the body can be replayed on retry. The call is bound to ctx: a
// cancelled context aborts before the next attempt and interrupts back-off
// sleeps.
func (c *Client) do(ctx context.Context, endpoint, contentType string, bodyFn func() io.Reader, out any) error {
	maxAttempts := max(c.maxRetries+1, 1)

	var lastErr error
	for attempt := range maxAttempts {
		if err := ctx.Err(); err != nil {
			return err
		}
		c.rateLimiter.wait()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bodyFn())
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < maxAttempts-1 {
				d := backoffDuration(attempt)
				c.logger.Warnf("telegram transport error %v: retrying in %v (attempt %d/%d)", err, d, attempt+1, maxAttempts)
				if serr := c.backoffSleep(ctx, d); serr != nil {
					return serr
				}
				continue
			}
			return err
		}

		body, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			if attempt < maxAttempts-1 {
				d := backoffDuration(attempt)
				c.logger.Warnf("telegram body read error %v: retrying in %v (attempt %d/%d)", readErr, d, attempt+1, maxAttempts)
				if serr := c.backoffSleep(ctx, d); serr != nil {
					return serr
				}
				continue
			}
			return readErr
		}

		var apiResp apiResponse
		parsed := json.Unmarshal(body, &apiResp) == nil

		// Explicit HTTP 429 — Telegram's documented flood control path.
		if resp.StatusCode == http.StatusTooManyRequests {
			ra := retryAfterFrom(resp, &apiResp, parsed)
			apiErr := &APIError{
				StatusCode:  resp.StatusCode,
				ErrorCode:   apiRespCode(parsed, &apiResp, resp.StatusCode),
				Description: apiRespDesc(parsed, &apiResp, resp.Status, body),
				RetryAfter:  ra,
			}
			if attempt < maxAttempts-1 {
				d := waitFor(ra, attempt)
				c.logger.Warnf("telegram 429: retrying in %v (attempt %d/%d)", d, attempt+1, maxAttempts)
				lastErr = apiErr
				if serr := c.backoffSleep(ctx, d); serr != nil {
					return serr
				}
				continue
			}
			return apiErr
		}

		// 5xx — transient server-side issue, back off and retry.
		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			apiErr := &APIError{
				StatusCode:  resp.StatusCode,
				ErrorCode:   apiRespCode(parsed, &apiResp, resp.StatusCode),
				Description: apiRespDesc(parsed, &apiResp, resp.Status, body),
			}
			if attempt < maxAttempts-1 {
				d := backoffDuration(attempt)
				c.logger.Warnf("telegram %d: retrying in %v (attempt %d/%d)", resp.StatusCode, d, attempt+1, maxAttempts)
				lastErr = apiErr
				if serr := c.backoffSleep(ctx, d); serr != nil {
					return serr
				}
				continue
			}
			return apiErr
		}

		// Non-2xx outside 429/5xx — return without retry.
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return &APIError{
				StatusCode:  resp.StatusCode,
				ErrorCode:   apiRespCode(parsed, &apiResp, resp.StatusCode),
				Description: apiRespDesc(parsed, &apiResp, resp.Status, body),
			}
		}

		// 2xx but body wasn't valid JSON.
		if !parsed {
			return fmt.Errorf("tbot: unable to decode response: %s", truncate(body, 256))
		}

		// Telegram occasionally returns 200 with ok=false and
		// error_code=429 instead of an HTTP 429. Treat it the same way.
		if !apiResp.OK {
			ra := 0
			if apiResp.Parameters != nil {
				ra = apiResp.Parameters.RetryAfter
			}
			apiErr := &APIError{
				StatusCode:  resp.StatusCode,
				ErrorCode:   apiResp.ErrorCode,
				Description: apiResp.Description,
				Parameters:  apiResp.Parameters,
				RetryAfter:  ra,
			}
			if apiResp.ErrorCode == http.StatusTooManyRequests && attempt < maxAttempts-1 {
				d := waitFor(ra, attempt)
				c.logger.Warnf("telegram 200/429: retrying in %v (attempt %d/%d)", d, attempt+1, maxAttempts)
				lastErr = apiErr
				if serr := c.backoffSleep(ctx, d); serr != nil {
					return serr
				}
				continue
			}
			return apiErr
		}

		if out == nil || len(apiResp.Result) == 0 {
			return nil
		}
		return json.Unmarshal(apiResp.Result, out)
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("tbot: request failed after %d attempts", maxAttempts)
}

// backoffSleep waits for d before the next retry. It honours ctx so a
// cancelled or expired context aborts the wait immediately. The c.sleepFn hook
// (used by tests) takes precedence when set.
func (c *Client) backoffSleep(ctx context.Context, d time.Duration) error {
	if c.sleepFn != nil {
		c.sleepFn(d)
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// backoffDuration returns an exponential backoff delay capped at
// backoffMax: 500ms, 1s, 2s, 4s, …
func backoffDuration(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	shift := min(uint(attempt), 16)
	d := backoffBase << shift
	if d > backoffMax || d <= 0 {
		return backoffMax
	}
	return d
}

// waitFor converts a retry_after (seconds) into a duration, falling
// back to exponential backoff when the server didn't tell us how long
// to wait.
func waitFor(retryAfter, attempt int) time.Duration {
	if retryAfter > 0 {
		return time.Duration(retryAfter) * time.Second
	}
	return backoffDuration(attempt)
}

// retryAfterFrom extracts the retry_after hint from (in priority
// order) the parsed parameters, the Retry-After header.
func retryAfterFrom(resp *http.Response, apiResp *apiResponse, parsed bool) int {
	if parsed && apiResp.Parameters != nil && apiResp.Parameters.RetryAfter > 0 {
		return apiResp.Parameters.RetryAfter
	}
	if h := resp.Header.Get("Retry-After"); h != "" {
		if v, err := strconv.Atoi(strings.TrimSpace(h)); err == nil && v > 0 {
			return v
		}
	}
	return 0
}

func apiRespCode(parsed bool, apiResp *apiResponse, statusCode int) int {
	if parsed && apiResp.ErrorCode != 0 {
		return apiResp.ErrorCode
	}
	return statusCode
}

func apiRespDesc(parsed bool, apiResp *apiResponse, status string, body []byte) string {
	if parsed && apiResp.Description != "" {
		return apiResp.Description
	}
	if len(body) > 0 {
		return fmt.Sprintf("%s: %s", status, truncate(body, 256))
	}
	return status
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "…"
}
