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
- [ ] `GetUpdates` (long polling — `Client` already has `nextOffset`,
      `bufferSize`, `timeout`, `updateParams` fields prepared for it but no
      method)
- [ ] `SetWebhook`, `DeleteWebhook`, `GetWebhookInfo`

### Messaging — send
- [ ] `SendPhoto`, `SendAudio`, `SendDocument`, `SendVideo`, `SendAnimation`,
      `SendVoice`, `SendVideoNote`, `SendPaidMedia`, `SendMediaGroup`
- [ ] `SendLocation`, `SendVenue`, `SendContact`, `SendPoll`, `SendDice`
- [ ] `SendChecklist`, `SendGame`
- [ ] `CopyMessage`, `CopyMessages`, `ForwardMessages`

### Messaging — edit/delete/react
- [ ] `EditMessageText`, `EditMessageCaption`, `EditMessageMedia`,
      `EditMessageReplyMarkup`, `EditMessageLiveLocation`,
      `StopMessageLiveLocation`, `EditMessageChecklist`
- [ ] `StopPoll`
- [ ] `DeleteMessage`, `DeleteMessages`
- [ ] `SetMessageReaction`

### Chat administration
- [ ] `BanChatMember`, `UnbanChatMember`, `RestrictChatMember`,
      `PromoteChatMember`, `SetChatAdministratorCustomTitle`,
      `BanChatSenderChat`, `UnbanChatSenderChat`
- [ ] `SetChatPermissions`, `ExportChatInviteLink`, `CreateChatInviteLink`,
      `EditChatInviteLink`, `RevokeChatInviteLink`,
      `CreateChatSubscriptionInviteLink`, `EditChatSubscriptionInviteLink`,
      `ApproveChatJoinRequest`, `DeclineChatJoinRequest`
- [ ] `SetChatPhoto`, `DeleteChatPhoto`, `SetChatTitle`,
      `SetChatDescription`, `PinChatMessage`, `UnpinChatMessage`,
      `UnpinAllChatMessages`, `LeaveChat`
- [ ] `GetChat`, `GetChatAdministrators`, `GetChatMemberCount`,
      `GetChatMember`, `SetChatStickerSet`, `DeleteChatStickerSet`,
      `GetForumTopicIconStickers`

### Forum topics
- [ ] `CreateForumTopic`, `EditForumTopic`, `CloseForumTopic`,
      `ReopenForumTopic`, `DeleteForumTopic`,
      `UnpinAllForumTopicMessages`, `EditGeneralForumTopic`,
      `CloseGeneralForumTopic`, `ReopenGeneralForumTopic`,
      `HideGeneralForumTopic`, `UnhideGeneralForumTopic`,
      `UnpinAllGeneralForumTopicMessages`

### Bot configuration
- [ ] `SetMyName`, `GetMyName`, `SetMyDescription`, `GetMyDescription`,
      `SetMyShortDescription`, `GetMyShortDescription`
- [ ] `SetMyCommands`, `DeleteMyCommands`, `GetMyCommands`
- [ ] `SetChatMenuButton`, `GetChatMenuButton`,
      `SetMyDefaultAdministratorRights`, `GetMyDefaultAdministratorRights`

### Inline mode
- [ ] `AnswerInlineQuery`, `AnswerWebAppQuery`, `SavePreparedInlineMessage`
- [ ] `AnswerCallbackQuery`

### Stickers
- [ ] `UploadStickerFile`, `CreateNewStickerSet`, `AddStickerToSet`,
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
- [ ] `GetFile`, `GetUserProfilePhotos`, `GetUserChatBoosts`,
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

- [ ] **Typed errors.** `sendRequest` returns
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
- [ ] **`context.Context` on all methods.** None of the current methods
      accept a `Context`, which means cancellation, deadlines, and trace
      propagation are impossible. Take it as the first argument
      consistently when adding new methods, and add `…Ctx` variants for
      the existing ones (or break the signatures — the package is
      pre-`v1`).
- [ ] **`ChatID` parameter type.** Telegram accepts either an integer or a
      `@username`. The current `SendMessage` takes `chatID string` so the
      caller does the conversion. Introduce a small `ChatID` type with
      `Int64(int64)` / `Username(string)` constructors and `MarshalJSON`
      to serialize correctly. Apply consistently to every new method.
- [ ] **`InputFile` is a string alias.** That works for `file_id` and URL
      cases but doesn't help with multipart uploads. Replace with a
      proper interface or struct that the request layer recognizes and
      attaches to the multipart body, so callers don't have to thread
      `inputFile{field, name}` through `sendRequestWithFiles` themselves.
- [ ] **Union-type helpers.** Add small predicate methods on the tagged
      structs — e.g. `MessageOrigin.IsUser()`, `ChatMember.IsAdministrator()`,
      `ReactionType.IsEmoji()` — so callers can avoid bare string
      comparisons. Optionally add typed accessors that return the
      variant-specific subset.
- [ ] **`MaybeInaccessibleMessage`.** Currently embeds `Message`. Add a
      `IsInaccessible()` method (`return m.Date == 0`) and a doc note.
- [ ] **Functional options on send methods.** The current `sendOption`
      pattern is fine for simple flags; for newer features (`reply_parameters`,
      `link_preview_options`, `business_connection_id`, `effect_id`,
      `message_effect_id`, `allow_paid_broadcast`) consider exposing them
      via typed setters or a `SendMessageOptions` struct so they're
      discoverable in IDEs.
- [ ] **`SendMessage` thread_id parameter.** The signature is
      `SendMessage(chatID string, thread_id string, text string, …)`.
      Snake-case parameter names are non-idiomatic, and a string thread ID
      is unusual. Either rename to `threadID int` or fold both into an
      options struct.

## Correctness margin

- [ ] **Race in `sendRequestWithFiles`.** `req` and `resp` are written
      inside a goroutine and read after the `<-done` join, which is
      synchronized — but `err` is also captured in the closure, written
      both inside the goroutine and at file-open time on the main path,
      and the multipart writer can fail without the goroutine knowing.
      The control flow is fragile. Refactor to a single goroutine that
      owns the pipe writer and returns its error via a channel; the main
      goroutine performs the HTTP request synchronously after the writer
      goroutine closes its end.
- [ ] **Logger interface is exposed but barely wired.** `c.logger` is
      consulted in exactly one error path
      (`sendRequestWithFiles` → `Error(err)`) and is left zero-valued —
      so a caller who never calls a `WithLogger` option will hit a nil
      dereference if that path fires. Default the logger to `nopLogger{}`
      in `NewClient`, and either log other request lifecycle events or
      remove the unused interface methods from `Logger`.
- [ ] **`fmt.Errorf(apiResp.Description)`** in
      `sendRequestWithFiles`'s success path treats user-supplied data as a
      format string. Use `errors.New(apiResp.Description)` or
      `fmt.Errorf("%s", apiResp.Description)`.
- [ ] **Deprecated `Thumb` JSON tag.** The legacy `Thumb` fields kept on
      media types still carry `json:"thumb,omitempty"`. The fields are
      receive-only in practice, so this is harmless today, but a caller
      who manually populates `Thumb` would emit both `thumb` and
      `thumbnail` on the wire. When/if these are removed, replace the
      tag with `json:"-"` first as a safer interim.
- [ ] **Custom `UnmarshalJSON` for unions** (optional). The flat-struct
      approach reads fine for unmarshalling, but writers who construct a
      union manually can leave variant fields populated for the wrong
      `Type`. Adding a `MarshalJSON` that whitelists fields per variant
      would prevent that.

## Tests

- [ ] **Real unit tests.** `client_test.go` defines two
      `TestClient_*` functions whose `tests` slice is empty, so
      they're effectively no-ops. Add `httptest.Server`-backed tests
      that:
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
