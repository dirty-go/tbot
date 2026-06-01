// Package tbot is a zero-dependency Telegram Bot API client for Go.
//
// tbot wraps the [Telegram Bot API] (currently Bot API 9.6) using only the
// Go standard library. Every method takes a [context.Context] as its first
// argument; errors are returned as typed values that callers can inspect with
// [errors.As].
//
// # Quick start
//
// Create a client and send a message in three lines:
//
//	c := tbot.NewClient("123456:ABC-DEF...", "")
//	chatID := tbot.Int64(-1001234567890)
//	msg, err := c.SendMessage(ctx, chatID, "Hello, world!")
//
// # Client
//
// [NewClient] accepts the bot token, an optional base-URL override (pass ""
// for the public api.telegram.org), and zero or more [ClientOptions]:
//
//   - [WithBaseURL] — point at a local Bot API server or a test fake
//   - [WithRateLimit] — token-bucket throttle (default: 30 req/s, burst 30)
//   - [WithMaxRetries] — retry budget for 429 / 5xx (default: 3)
//   - [WithPollTimeout] — long-poll timeout sent to getUpdates
//   - [WithHTTPClient] — supply a pre-configured *http.Client
//   - [WithHTTPTimeout] — override the timeout on the default HTTP client
//   - [WithLogger] — swap the default no-op logger (see [Logger])
//
// [Client] methods are safe for concurrent use by multiple goroutines. The
// single exception is [Client.GetUpdates], which mutates the long-poll offset
// and must be driven from one goroutine at a time.
//
// # ChatID
//
// [ChatID] identifies a target chat. The Bot API accepts either a numeric
// identifier or a @-prefixed username:
//
//	chatID := tbot.Int64(-1001234567890)   // numeric
//	chatID := tbot.Username("@mychannel")  // public username
//	if chatID.IsZero() { /* unset */ }
//
// # InputFile
//
// [InputFile] is the value accepted by all methods that send files. The
// request layer selects form-encoding for references and multipart delivery
// for uploads automatically; callers never assemble multipart parts:
//
//	tbot.FileID("AgACAgIAAxkB...")            // file already on Telegram
//	tbot.FileURL("https://example.com/f.png") // Telegram fetches from URL
//	tbot.FilePath("/tmp/report.pdf")           // upload from local filesystem
//	tbot.FileReader("data.csv", r)             // upload from io.Reader
//	tbot.FileBytes("img.png", pngBytes)        // upload from memory
//
// # Errors
//
// Every client method returns an error that wraps [*APIError] whenever the
// Bot API replies with ok=false or an HTTP error status. Use [errors.As] to
// inspect the structured fields:
//
//	var apiErr *tbot.APIError
//	if errors.As(err, &apiErr) {
//	    fmt.Printf("code=%d desc=%q retryAfter=%d\n",
//	        apiErr.ErrorCode, apiErr.Description, apiErr.RetryAfter)
//
//	    // Handle group-to-supergroup migration.
//	    if apiErr.Parameters != nil && apiErr.Parameters.MigrateToChatID != 0 {
//	        newID := tbot.Int64(apiErr.Parameters.MigrateToChatID)
//	        // retry against newID
//	    }
//	}
//
// RetryAfter is populated from parameters.retry_after (or the Retry-After
// header on HTTP 429). The client already retries on 429 and 5xx within its
// retry budget, so callers typically see [*APIError] only after retries are
// exhausted.
//
// # Functional options
//
// Methods with many optional Bot API parameters accept a variadic
// [...SendOption] tail. Pre-built options live in options.go and compose
// freely:
//
//	msg, err := c.SendMessage(ctx, chatID, "Hello <b>world</b>",
//	    tbot.OptParseModeHTML,
//	    tbot.OptDisableNotification,
//	    tbot.OptProtectContent,
//	    tbot.OptReplyParameters(&tbot.ReplyParameters{MessageID: 42}),
//	    tbot.OptInlineKeyboardMarkup(&tbot.InlineKeyboardMarkup{...}),
//	)
//
// Commonly used options include:
//
//   - [OptParseModeHTML] / [OptParseModeMarkdown] — text formatting mode
//   - [OptDisableNotification] — send silently
//   - [OptProtectContent] — forbid forwarding and saving
//   - [OptAllowPaidBroadcast] — paid broadcast (up to 1000 msg/s)
//   - [OptMessageThreadID] — target a forum topic
//   - [OptBusinessConnectionID] — send via a connected business account
//   - [OptMessageEffectID] — attach a message effect (private chats only)
//   - [OptReplyParameters] — reply with optional cross-chat quote
//   - [OptLinkPreviewOptions] — configure link preview generation
//   - [OptCaption] — set a media caption
//   - [OptInlineKeyboardMarkup] / [OptReplyKeyboardMarkup] — attach markup
//
// # SendMessageOptions struct
//
// For callers who prefer struct literals over building a []SendOption slice,
// [SendMessageOptions] collects all common options into one struct. Call
// [SendMessageOptions.Apply] to convert it to the []SendOption slice accepted
// by every send method:
//
//	opts := tbot.SendMessageOptions{
//	    ParseMode:           "HTML",
//	    DisableNotification: true,
//	    ReplyParameters:     &tbot.ReplyParameters{MessageID: 5},
//	}
//	msg, err := c.SendMessage(ctx, chatID, "text", opts.Apply()...)
//
// The two styles interoperate freely:
//
//	extra := append(opts.Apply(), tbot.OptProtectContent)
//	msg, err := c.SendPhoto(ctx, chatID, photo, extra...)
//
// # Union-type helpers
//
// Several Telegram types are discriminated unions carrying a string Type (or
// Source / Status) field. Predicate helpers eliminate raw string comparisons:
//
// [MessageOrigin] (forwarded-message origin):
//
//	if origin.IsUser()       { /* sender_user */ }
//	if origin.IsHiddenUser() { /* sender_user_name */ }
//	if origin.IsChat()       { /* sender_chat */ }
//	if origin.IsChannel()    { /* chat + message_id */ }
//
// [ChatMember] (membership status):
//
//	if member.IsCreator()       { /* owner */ }
//	if member.IsAdministrator() { /* admin */ }
//	if member.IsMemberStatus()  { /* plain member */ }
//	if member.IsRestricted()    { /* restricted */ }
//	if member.HasLeft()         { /* left */ }
//	if member.IsBanned()        { /* kicked */ }
//
// [ChatBoostSource] (boost origin):
//
//	if src.IsPremium()  { /* Premium subscription */ }
//	if src.IsGiftCode() { /* redeemed gift code */ }
//	if src.IsGiveaway() { /* giveaway */ }
//
// [ReactionType] (message reaction):
//
//	if r.IsEmoji()       { /* standard emoji */ }
//	if r.IsCustomEmoji() { /* custom emoji */ }
//	if r.IsPaid()        { /* paid (star) reaction */ }
//
// Union types sent back to the API implement MarshalJSON, which emits only
// the fields valid for the selected variant.
//
// # Logger
//
// [Logger] is an interface with leveled printf-style and print-style methods
// (Debugf/Infof/Printf/Warnf/Errorf and Debug/Info/Print/Warn/Error). The
// default implementation is a silent no-op. [BasicLogger] delegates to the
// standard library's log package. Swap in any compatible logger with
// [WithLogger]:
//
//	c := tbot.NewClient(token, "", tbot.WithLogger(tbot.BasicLogger{}))
//
// # Long polling and webhooks
//
// See examples/longpoll/main.go for a minimal long-poll bot driven by
// [Client.GetUpdates], and examples/webhook/main.go for a webhook-based bot
// using [Client.SetWebhook].
//
// # Deprecation policy
//
// The package is pre-v1; all APIs are subject to change. Older Bot API fields
// are retained for source compatibility and carry a Deprecated doc comment:
//
//   - Animation.Thumb, Document.Thumb, Video.Thumb, VideoNote.Thumb,
//     Sticker.Thumb — use Thumbnail (renamed in Bot API 6.6; json:"-")
//   - Animation.MimeType, Video.MimeType, Voice.MimeType — use MIMEType
//   - Message.ForwardFrom / ForwardFromChat / ForwardDate — use
//     Message.ForwardOrigin and switch on its Type (replaced in Bot API 7.0)
//   - Message.VoiceChatScheduled / VoiceChatStarted / VoiceChatEnded /
//     VoiceChatParticipantsInvited — use the VideoChat* equivalents
//     (renamed in Bot API 6.0)
//   - Poll.CorrectOptionID (scalar) — use Poll.CorrectOptionIDs (array)
//   - ChatPermissions.CanSendMediaMessages — use the granular per-media-type
//     fields (replaced in Bot API 6.5)
//   - OptReplyToMessageID / OptDisableWebPagePreview — use
//     OptReplyParameters / OptLinkPreviewOptions
//
// [Telegram Bot API]: https://core.telegram.org/bots/api
package tbot
