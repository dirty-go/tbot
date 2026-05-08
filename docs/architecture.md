# Architecture

**tbot** is a single-package (`package tbot`) Telegram Bot API client library with **zero external dependencies** — stdlib only. This constraint is intentional and must be preserved.

## Layout

| File(s) | Role |
| --- | --- |
| `client.go` | `Client` struct, `NewClient`, public API methods (`Me`, `SendMessage`, `SendSticker*`, `ForwardMessage`, `SendChatAction`). Also defines `sendOption` and the per-send `Opt*` package vars (`OptParseModeHTML`, `OptReplyToMessageID`, `OptInlineKeyboardMarkup`, …). |
| `api.go` | HTTP transport: `sendRequest` / `sendRequestWithFiles` (multipart via `io.Pipe`), shared `netTransport`, `apiResponse` envelope. |
| `client_options.go` | `ClientOptions` functional-option type + `WithBaseURL`. **Note:** currently defined but not wired into `NewClient` — `NewClient(token, baseURL)` takes baseURL positionally. |
| `tbot.go` | Package-level `apiBaseURL = "https://api.telegram.org"`. |
| `logger.go` | `Logger` interface, `nopLogger` (default), `BasicLogger` (stdlib `log` adapter). |
| `types.go` + `types_*.go` | Telegram domain types, split by domain (chat, message, service, keyboard, sticker, reaction, forum, business, boost, giveaway, payment, gift, stars, passport, inline, input, bot, update, user). The canonical map is in the package doc at the top of `types.go`. |

## Key patterns

**Functional options for send calls** — `sendOption` is `func(url.Values)`. Adding a new Telegram parameter means adding a new `Opt*` var or `func(...) sendOption` constructor; method signatures don't change.

**HTTP request flow** — public methods build `url.Values`, apply `opts...`, call `sendRequest` which POSTs to `{baseURL}/bot{token}/{method}`, parses the `apiResponse` envelope (`{"ok": bool, "result": …}`), and unmarshals `Result` into the caller's response type.

**File uploads** — `sendRequestWithFiles` writes multipart data through `io.Pipe` concurrently with the HTTP request, avoiding full in-memory buffering.

**Testing against a fake server** — pass the test `httptest.Server.URL` as the `baseURL` arg to `NewClient`; an empty string falls back to `apiBaseURL`.

**Type back-compat invariant** — when Telegram renames a field (e.g. `thumb` → `thumbnail`, `mime_type` casing), the old Go field is kept alongside the new one with a `// Deprecated:` comment and `json:"-"` so older callers still compile. Don't remove these when refreshing types to a newer Bot API version; only add the new field. Track which Bot API version `types.go`'s package doc references when bumping (currently 9.6).
