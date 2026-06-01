# TODO

A package-level audit produced after the Bot API 9.6 type refresh
(2026-05-07). Items are ordered roughly by user impact: things in *Method
coverage* are blocking real-world use of the library; items lower down are
ergonomics, correctness margin, and housekeeping.

## Method coverage (the big gap)

`Client` currently exposes 6 methods (`Me`, `SendMessage`, `ForwardMessage`,
`SendSticker`, `SendStickerFile`, `SendChatAction`). The Bot API defines
~168. Until the rest are wrapped, callers can't build a real bot with this
package — they fall back to using `sendRequest` directly. Group by domain:

### Updates / webhooks
- [x] `GetUpdates` (long polling — `Client` already has `nextOffset`,
      `bufferSize`, `timeout`, `updateParams` fields prepared for it but no
      method)
- [x] `SetWebhook`, `DeleteWebhook`, `GetWebhookInfo`

### Messaging — send
- [x] `SendPhoto`, `SendAudio`, `SendDocument`, `SendVideo`, `SendAnimation`,
      `SendVoice`, `SendVideoNote`, `SendMediaGroup` (file_id/URL/upload via the
      new `InputFile`). `SendPaidMedia` still TODO.
- [x] `SendLocation`, `SendVenue`, `SendContact`, `SendPoll`, `SendDice`
- [ ] `SendChecklist`, `SendGame`
- [x] `CopyMessage`, `CopyMessages`, `ForwardMessages`

### Messaging — edit/delete/react
- [x] `EditMessageText`, `EditMessageCaption`, `EditMessageMedia`,
      `EditMessageReplyMarkup`, `EditMessageLiveLocation`,
      `StopMessageLiveLocation` (edit-of-inline returns `(nil, nil)`).
      `EditMessageChecklist` still TODO.
- [x] `StopPoll`
- [x] `DeleteMessage`, `DeleteMessages`
- [x] `SetMessageReaction`

### Chat administration
- [x] `BanChatMember`, `UnbanChatMember`, `RestrictChatMember`,
      `PromoteChatMember`, `SetChatAdministratorCustomTitle`,
      `BanChatSenderChat`, `UnbanChatSenderChat`
- [x] `SetChatPermissions`, `ExportChatInviteLink`, `CreateChatInviteLink`,
      `EditChatInviteLink`, `RevokeChatInviteLink`,
      `CreateChatSubscriptionInviteLink`, `EditChatSubscriptionInviteLink`,
      `ApproveChatJoinRequest`, `DeclineChatJoinRequest`
- [x] `SetChatPhoto`, `DeleteChatPhoto`, `SetChatTitle`,
      `SetChatDescription`, `PinChatMessage`, `UnpinChatMessage`,
      `UnpinAllChatMessages`, `LeaveChat`
- [x] `GetChat`, `GetChatAdministrators`, `GetChatMemberCount`,
      `GetChatMember`, `SetChatStickerSet`, `DeleteChatStickerSet`,
      `GetForumTopicIconStickers`

### Forum topics
- [x] `CreateForumTopic`, `EditForumTopic`, `CloseForumTopic`,
      `ReopenForumTopic`, `DeleteForumTopic`,
      `UnpinAllForumTopicMessages`, `EditGeneralForumTopic`,
      `CloseGeneralForumTopic`, `ReopenGeneralForumTopic`,
      `HideGeneralForumTopic`, `UnhideGeneralForumTopic`,
      `UnpinAllGeneralForumTopicMessages`

### Bot configuration
- [x] `SetMyName`, `GetMyName`, `SetMyDescription`, `GetMyDescription`,
      `SetMyShortDescription`, `GetMyShortDescription`
- [x] `SetMyCommands`, `DeleteMyCommands`, `GetMyCommands`
- [x] `SetChatMenuButton`, `GetChatMenuButton`,
      `SetMyDefaultAdministratorRights`, `GetMyDefaultAdministratorRights`

### Inline mode
- [x] `AnswerInlineQuery`, `AnswerWebAppQuery`, `SavePreparedInlineMessage`
- [x] `AnswerCallbackQuery`

### Stickers
- [x] `UploadStickerFile`, `CreateNewStickerSet`, `AddStickerToSet`,
      `SetStickerPositionInSet`, `DeleteStickerFromSet`, `ReplaceStickerInSet`,
      `SetStickerEmojiList`, `SetStickerKeywords`, `SetStickerMaskPosition`,
      `SetStickerSetTitle`, `SetStickerSetThumbnail`,
      `SetCustomEmojiStickerSetThumbnail`, `DeleteStickerSet`,
      `GetStickerSet`, `GetCustomEmojiStickers`

### Payments / Stars / gifts
- [ ] `SendInvoice`, `CreateInvoiceLink`, `AnswerShippingQuery`,
      `AnswerPreCheckoutQuery`, `RefundStarPayment`,
      `EditUserStarSubscription`
- [ ] `GetStarTransactions`, `GetMyStarBalance`
- [ ] `GetAvailableGifts`, `SendGift`, `GiftPremiumSubscription`,
      `VerifyUserGift`, `ConvertGiftToStars`, `UpgradeGift`,
      `TransferGift`, `GetReceivedGifts`, `SaveGift`

### Misc
- [x] `GetFile`
- [x] `GetUserProfilePhotos`
- [ ] `GetUserChatBoosts`,
      `GetBusinessConnection`, `SetUserEmojiStatus`,
      `ReadBusinessMessage`, `DeleteBusinessMessages`,
      `SetBusinessAccountName`, `SetBusinessAccountUsername`,
      `SetBusinessAccountBio`, `SetBusinessAccountProfilePhoto`,
      `RemoveBusinessAccountProfilePhoto`,
      `SetBusinessAccountGiftSettings`, `GetBusinessAccountStarBalance`,
      `GetBusinessAccountGifts`, `TransferBusinessAccountStars`
- [ ] Game scores: `SetGameScore`, `GetGameHighScores`

A practical first cut is the *Send / edit / delete* family plus
`GetUpdates` / `SetWebhook` and `AnswerCallbackQuery` — that lets people
write the most common bot.

## Ergonomics

- [x] **Typed errors.** `sendRequest` returns
      `fmt.Errorf("%d : %s", code, desc)` — callers can't programmatically
      detect flood-control (`retry_after`) or migration
      (`migrate_to_chat_id`). Introduce
      ```go
      type APIError struct {
          Code        int
          Description string
          Parameters  *ResponseParameters
      }
      ```
      so callers can use `errors.As`.
- [x] **`context.Context` on all methods.** Every method now takes
      `ctx context.Context` as its first argument (breaking; the package is
      pre-`v1`). `ctx` is plumbed through `sendRequest` / `sendMultipart` /
      `do` via `http.NewRequestWithContext`, and `backoffSleep` honours
      cancellation between and during retry waits.
- [x] **`ChatID` parameter type.** Added `ChatID` with `Int64(int64)` /
      `Username(string)` constructors, `String()` (form encoding),
      `MarshalJSON` (number vs string), and `IsZero()`. Used by every method
      that takes a chat target.
- [x] **`InputFile` is a string alias.** Replaced with a sealed `InputFile`
      interface plus `FileID` / `FileURL` / `FilePath` / `FileReader` /
      `FileBytes` constructors. The request layer (`addFile` → `send`)
      routes references to the form body and uploads to multipart parts,
      so callers no longer thread `inputFile{field,name}` themselves.
- [~] **Union-type helpers.** Done for the unions touched by this pass:
      `ReactionType.IsEmoji/IsCustomEmoji/IsPaid` (+ `EmojiReaction` /
      `CustomEmojiReaction` constructors). `MessageOrigin.IsUser`,
      `ChatMember.IsAdministrator`, etc. still TODO (admin/forward domains).
- [x] **`MaybeInaccessibleMessage`.** Added `IsInaccessible()`
      (`return m.Date == 0`) with a doc note.
- [~] **Functional options on send methods.** Exported `SendOption` and added
      typed setters incl. `OptReplyParameters`, `OptLinkPreviewOptions`,
      `OptBusinessConnectionID`, `OptMessageEffectID`, `OptAllowPaidBroadcast`,
      `OptMessageThreadID`, `OptCaption`, media/location setters, and edit
      targeting (`OptChatID` / `OptMessageID` / `OptInlineMessageID`). A
      consolidated `SendMessageOptions` struct is still open.
- [x] **`SendMessage` thread_id parameter.** Dropped the snake-case string
      `thread_id` positional; thread targeting is now `OptMessageThreadID(int)`.

## Correctness margin

- [x] **Race in `sendRequestWithFiles`.** `req` and `resp` are written
      inside a goroutine and read after the `<-done` join, which is
      synchronized — but `err` is also captured in the closure, written
      both inside the goroutine and at file-open time on the main path,
      and the multipart writer can fail without the goroutine knowing.
      The control flow is fragile. Refactor to a single goroutine that
      owns the pipe writer and returns its error via a channel; the main
      goroutine performs the HTTP request synchronously after the writer
      goroutine closes its end.
- [x] **Logger interface is exposed but barely wired.** `c.logger` is
      consulted in exactly one error path
      (`sendRequestWithFiles` → `Error(err)`) and is left zero-valued —
      so a caller who never calls a `WithLogger` option will hit a nil
      dereference if that path fires. Default the logger to `nopLogger{}`
      in `NewClient`, and either log other request lifecycle events or
      remove the unused interface methods from `Logger`.
- [x] **`fmt.Errorf(apiResp.Description)`** in
      `sendRequestWithFiles`'s success path treats user-supplied data as a
      format string. Use `errors.New(apiResp.Description)` or
      `fmt.Errorf("%s", apiResp.Description)`.
- [x] **Deprecated `Thumb` JSON tag.** The legacy `Thumb` fields kept on
      media types still carry `json:"thumb,omitempty"`. The fields are
      receive-only in practice, so this is harmless today, but a caller
      who manually populates `Thumb` would emit both `thumb` and
      `thumbnail` on the wire. When/if these are removed, replace the
      tag with `json:"-"` first as a safer interim.
- [x] **Custom `UnmarshalJSON` for unions** (optional). The flat-struct
      approach reads fine for unmarshalling, but writers who construct a
      union manually can leave variant fields populated for the wrong
      `Type`. Adding a `MarshalJSON` that whitelists fields per variant
      would prevent that.

## Tests

- [x] **Real unit tests.** The empty no-op `TestClient_*` stubs are gone.
      `httptest.Server`-backed tests now cover `Me`, the webhook trio,
      `AnswerCallbackQuery`, the location/venue/contact/poll/dice family,
      `SendPhoto` (reference + multipart upload), `SendMediaGroup`, the
      edit/delete/react/copy family, `StopPoll`, `ChatID`, `InputFile`, and
      the `APIError` / retry / context-cancellation paths. Original goal:
  - drive `NewClient(token, server.URL)` (the existing `WithBaseURL`
    option supports this)
  - cover golden-path JSON round-trips for the most-used types (the
    JSON shape is the contract this package promises)
  - cover the `apiResponse{ok:false}` error path and the new
    `APIError` once introduced
- [ ] **JSON round-trip tests** for every union type, exercising each
      `Type` discriminator value, to lock in marshal/unmarshal symmetry.
- [ ] **Integration test against the live API**, gated by an environment
      variable (`TBOT_TEST_TOKEN`), as a smoke test for protocol drift
      when Telegram updates the API.

## Documentation

- [ ] Update `README.md` to advertise the new types, the union-type
      conventions, and the deprecation policy.
- [ ] Add `examples/` showing webhook + long-poll setups, the keyboard
      types, and a small payment flow.
- [ ] Generate godoc-rendered package overview from the headers in
      `types.go`.

## Future cleanup

- [ ] **Pre-`v1` breaking change pass.** When the package cuts a `v1`,
      drop the deprecated `Thumb` fields, `VoiceChat*` aliases,
      `Forward*` legacy fields, `RequestPool` / `KeyboardButtonPoolType`,
      `ChatPermissions.CanSendMediaMessages`, `Poll.CorrectOptionID`,
      `Chat.AllMembersAreAdministrators`. They exist solely to keep
      pre-9.x callers compiling through this transition.
