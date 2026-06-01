# tbot — Telegram Bot API Client

[![GitHub Actions](https://github.com/dirty-go/tbot/workflows/Test/badge.svg)](https://github.com/dirty-go/tbot/actions)

![logo](https://raw.githubusercontent.com/dirty-go/tbot/main/logo.png)

## Features

- **Zero** external dependencies — stdlib only
- Type-safe API client with functional options
- Middleware support
- Go modules compatible
- External logger support

## Installation

```bash
go get github.com/dirty-go/tbot
```

Requires Go 1.24+.

---

## Core Types

### ChatID

`ChatID` identifies a target chat. The Bot API accepts either a numeric identifier or a `@`-prefixed username, so `ChatID` carries whichever form you supply and renders it correctly for both form-encoded request fields and JSON bodies.

```go
// Numeric chat identifier.
chatID := tbot.Int64(-1001234567890)

// Public channel by username.
chatID := tbot.Username("@mychannel")

// IsZero detects unset values.
if chatID.IsZero() { /* ... */ }
```

### InputFile

`InputFile` is the interface accepted by all methods that send files. The request layer automatically selects form-encoded delivery for references and multipart delivery for uploads — callers never assemble multipart parts themselves.

| Constructor | Description |
|---|---|
| `FileID(id string)` | Reference a file already on Telegram's servers by `file_id`. |
| `FileURL(url string)` | Let Telegram fetch the file from an HTTPS URL. |
| `FilePath(path string)` | Upload a file from the local filesystem. |
| `FileReader(name string, r io.Reader)` | Upload from an arbitrary `io.Reader`. |
| `FileBytes(name string, data []byte)` | Upload from an in-memory byte slice. |

```go
// Send a photo already on Telegram's servers.
c.SendPhoto(ctx, chatID, tbot.FileID("AgACAgIAAxkB..."))

// Upload a local file.
c.SendDocument(ctx, chatID, tbot.FilePath("/tmp/report.pdf"))

// Upload from memory.
c.SendAudio(ctx, chatID, tbot.FileBytes("song.mp3", audioData))
```

### APIError

When the Bot API returns `ok: false` or an HTTP error status, every client method returns an `*APIError`. Use `errors.As` to inspect it:

```go
msg, err := c.SendMessage(ctx, chatID, text)
if err != nil {
    var apiErr *tbot.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("code=%d desc=%q retryAfter=%d\n",
            apiErr.ErrorCode, apiErr.Description, apiErr.RetryAfter)

        // React to a group-to-supergroup migration.
        if apiErr.Parameters != nil && apiErr.Parameters.MigrateToChatID != 0 {
            newChatID := tbot.Int64(apiErr.Parameters.MigrateToChatID)
            // retry against newChatID
        }
    }
}
```

`RetryAfter` is populated from `parameters.retry_after` (or the `Retry-After` header on HTTP 429) when present. The client already retries on 429 and 5xx, so callers typically only see `*APIError` after the retry budget is exhausted.

---

## Union-Type Helpers

Several Telegram types are discriminated unions. Each carries a `Type` (or `Source` / `Status`) field with a string discriminator. Predicate helpers are provided so you do not have to compare string constants yourself.

### MessageOrigin

```go
origin := msg.ForwardOrigin // *tbot.MessageOrigin
if origin.IsUser()       { /* sender_user */ }
if origin.IsHiddenUser() { /* sender_user_name */ }
if origin.IsChat()       { /* sender_chat */ }
if origin.IsChannel()    { /* chat + message_id */ }
```

### ChatMember

```go
member, _ := c.GetChatMember(ctx, chatID, userID)
if member.IsCreator()       { /* owner */ }
if member.IsAdministrator() { /* admin */ }
if member.IsMemberStatus()  { /* plain member */ }
if member.IsRestricted()    { /* restricted */ }
if member.HasLeft()         { /* left */ }
if member.IsBanned()        { /* kicked */ }
```

### ChatBoostSource

```go
boost.Source.IsPremium()  // Premium subscription
boost.Source.IsGiftCode() // Redeemed gift code
boost.Source.IsGiveaway() // Giveaway
```

### ReactionType

```go
reaction.IsEmoji()       // standard emoji
reaction.IsCustomEmoji() // custom emoji
reaction.IsPaid()        // paid (star) reaction
```

All union types that are sent back to the API implement `MarshalJSON`, which emits only the fields valid for the selected variant — cross-variant fields are silently dropped.

---

## Functional Options

Methods that accept many optional Bot API parameters use the `SendOption` functional-option pattern. Pre-built options live in `options.go`; they compose freely.

```go
msg, err := c.SendMessage(ctx, chatID, "Hello <b>world</b>",
    tbot.OptParseModeHTML,
    tbot.OptDisableNotification,
    tbot.OptProtectContent,
    tbot.OptReplyParameters(&tbot.ReplyParameters{MessageID: 42}),
    tbot.OptInlineKeyboardMarkup(&tbot.InlineKeyboardMarkup{
        InlineKeyboard: [][]tbot.InlineKeyboardButton{{
            {Text: "Go!", URL: "https://example.com"},
        }},
    }),
)
```

### Common Options

| Option | Effect |
|---|---|
| `OptParseModeHTML` | Format text as HTML |
| `OptParseModeMarkdown` | Format text as MarkdownV2 |
| `OptDisableNotification` | Send silently |
| `OptProtectContent` | Forbid forwarding/saving |
| `OptAllowPaidBroadcast` | Paid broadcast (up to 1000/s) |
| `OptMessageThreadID(id)` | Target a forum topic |
| `OptBusinessConnectionID(id)` | Send via business account |
| `OptMessageEffectID(id)` | Attach a message effect |
| `OptReplyParameters(p)` | Reply with optional cross-chat quote |
| `OptLinkPreviewOptions(o)` | Configure link preview |
| `OptCaption(text)` | Set media caption |
| `OptThumbnail(ref)` | Set thumbnail reference |
| `OptDuration(s)` | Set media duration (seconds) |
| `OptWidth(px)` / `OptHeight(px)` | Set media dimensions |
| `OptInlineKeyboardMarkup(m)` | Attach inline keyboard |
| `OptReplyKeyboardMarkup(m)` | Attach reply keyboard |
| `OptReplyKeyboardRemove` | Remove custom keyboard |
| `OptForceReply` | Force reply interface |
| `OptChatID(id)` / `OptMessageID(id)` | Target a message to edit |
| `OptInlineMessageID(id)` | Target an inline message to edit |

### SendMessageOptions Struct

For callers who prefer struct literals over building a `[]SendOption` slice, `SendMessageOptions` collects all common options into a single struct. Call `Apply()` to convert it to the `[]SendOption` slice accepted by every send method:

```go
opts := tbot.SendMessageOptions{
    ParseMode:           "HTML",
    DisableNotification: true,
    ReplyParameters:     &tbot.ReplyParameters{MessageID: 5},
    ReplyMarkup:         tbot.OptInlineKeyboardMarkup(kb),
}

msg, err := c.SendMessage(ctx, chatID, "text", opts.Apply()...)
```

Both styles are fully supported and interoperate:

```go
// Mix struct options with raw OptXxx calls.
extra := append(opts.Apply(), tbot.OptProtectContent)
msg, err := c.SendPhoto(ctx, chatID, photo, extra...)
```

---

## Deprecation Policy

**Pre-v1 — all APIs are subject to change.**

Older Bot API fields are retained for source compatibility:

- `Animation.Thumb`, `Document.Thumb`, `Video.Thumb`, `VideoNote.Thumb`, `Sticker.Thumb` — use `Thumbnail` instead (renamed in Bot API 6.6). The deprecated fields have `json:"-"` and are never marshaled.
- `Animation.MimeType`, `Video.MimeType`, `Voice.MimeType` — use `MIMEType`.
- `Message.ForwardFrom`, `ForwardFromChat`, `ForwardDate`, etc. — use `Message.ForwardOrigin` and switch on its `Type` (replaced in Bot API 7.0).
- `Message.VoiceChatScheduled`, `VoiceChatStarted`, `VoiceChatEnded`, `VoiceChatParticipantsInvited` — use the `VideoChat*` equivalents (renamed in Bot API 6.0).
- `Poll.CorrectOptionID` (scalar) — use `Poll.CorrectOptionIDs` (array).
- `ChatPermissions.CanSendMediaMessages` — use the granular `CanSendAudios`, `CanSendDocuments`, etc. (replaced in Bot API 6.5).
- `OptReplyToMessageID` / `OptDisableWebPagePreview` — use `OptReplyParameters` / `OptLinkPreviewOptions`.

All deprecated identifiers carry a `// Deprecated:` doc comment pointing to the replacement.

---

## Integration Tests

Integration tests (those that call the real Bot API) are gated behind the `integration` build tag and two environment variables:

```bash
export TBOT_TEST_TOKEN="123456:ABC-DEF..."
export TBOT_TEST_CHAT_ID="-1001234567890"   # optional; required for SendMessage test

go test -tags=integration ./...
```

If `TBOT_TEST_TOKEN` is not set the tests are skipped automatically. `TBOT_TEST_CHAT_ID` is only required for `TestIntegration_SendMessage`.
