# 2026-05-07 — Telegram Bot API 9.6 type refresh

Brought the type definitions in line with Telegram Bot API 9.6 (the package
previously tracked roughly Bot API 5.x). Method coverage on `Client` is
unchanged; this is purely a types/structs update plus a few correctness fixes
that surface as Go API changes.

## Highlights

- **Full Bot API 9.6 type parity.** Every type referenced by the API
  reference now exists in the package, including `Update` and `WebhookInfo`
  (which were missing entirely), the new `MessageOrigin` / `ReplyParameters`
  / `LinkPreviewOptions` model, the reaction system, forum topics, business
  connections, chat boosts, giveaways, paid media, gifts, Star transactions,
  inline-mode types, and the bot-config types.
- **Domain-grouped layout.** `types.go` now contains only small media leaves
  and carries a header comment indexing the new files. The previous monolith
  is split into 18 `types_*.go` files (chat, message, service, business,
  reaction, forum, boost, giveaway, keyboard, sticker, payment, gift, stars,
  passport, inline, input, bot, update). `User`, `ChatPhoto`,
  `ChatPermissions`, `LoginURL`, the unexported `replyKeyboardRemove` /
  `forceReply` and the unexported `responseParameters` were extracted from
  `client.go` / `api.go` into their domain files; the latter two are now
  exported as `ReplyKeyboardRemove`, `ForceReply`, `ResponseParameters`.

## Bug fixes (Go API changes)

- `ReplyKeyboardMarkup.Keyboard` is now `[][]KeyboardButton` (rows of
  buttons) — previously a flat `[]KeyboardButton`, which would have produced
  a single-row keyboard at best and not matched the API.
- `PollAnswer.PollID` is now `string`. The API has always sent a string; the
  previous `int` would never have unmarshaled.
- `KeyboardButton.RequestPoll` corrects the previous typo `RequestPool`. The
  old field is kept as a deprecated alias.
- `KeyboardButtonPollType` is the correct spelling. `KeyboardButtonPoolType`
  is a deprecated type alias.
- All chat / user / sender / forward / migrate identifiers are now `int64`.
  The API documents that these may exceed 32-bit precision.
- All `FileSize` fields are now `int64` for the same reason.
- `Contact.UserID` is now `int64`.
- `MessageEntity.URL`, `MessageEntity.User`, `MessageEntity.Language`,
  `MessageEntity.CustomEmojiID` and the new `UnixTime` /  `DateTimeFormat`
  fields are all wired up; ditto for new fields on `Sticker`, `Video`,
  `Audio`, `Document`, `Animation`, `Location`, `Venue`, `Poll`, `User`.

## Deprecation policy

Where Telegram renamed a field/type (Bot API 5.x → 9.6), the old Go
field/type is kept and annotated with `// Deprecated:` pointing at the
replacement. Callers compile unchanged. Specifically:

- `Thumb` is kept on every media type alongside the new `Thumbnail`
  (renamed in Bot API 6.6).
- `VoiceChat*` exist as type aliases for `VideoChat*` (renamed in Bot
  API 6.0).
- `Message.ForwardFrom`, `ForwardFromChat`, `ForwardSignature`,
  `ForwardSenderName`, `ForwardDate` are kept alongside the new
  `ForwardOrigin` (Bot API 7.0).
- `Chat.AllMembersAreAdministrators` is kept though the API no longer
  returns it.
- `ChatPermissions.CanSendMediaMessages` is kept alongside the granular
  `CanSendAudios` / `CanSendDocuments` / `CanSendPhotos` / `CanSendVideos` /
  `CanSendVideoNotes` / `CanSendVoiceNotes` (Bot API 6.5 split).
- `Poll.CorrectOptionID` is kept alongside the now-array `CorrectOptionIDs`.

## Union types

Telegram's union types (`MessageOrigin`, `ChatMember`, `ReactionType`,
`BackgroundType`, `BackgroundFill`, `PaidMedia`, `TransactionPartner`,
`RevenueWithdrawalState`, `BotCommandScope`, `MenuButton`,
`InlineQueryResult`, `InputMessageContent`, `InputMedia`, `InputPaidMedia`,
`PassportElementError`, `OwnedGift`, `ChatBoostSource`) are modeled as a
single tagged struct with a `Type` (or `Source` / `Status`) discriminator and
a flattened union of variant fields. Companion `const` values list the valid
discriminator strings. Switch on the discriminator to read variant-specific
fields.

This deliberately differs from the per-variant struct approach because Go's
`encoding/json` does not natively dispatch unions, and a typed-union mock
would either need custom `UnmarshalJSON` everywhere or force callers to do
their own JSON decoding. The flat-struct approach reads cleanly with
`encoding/json` out of the box.

## Build / test

`go build ./...` and `go test ./...` both pass. The existing test suite is
the placeholder shipped in `client_test.go`; it covers the unchanged
constructor and `Me` shape only. No regressions.

## Out of scope

Method coverage on `Client` remains at the original six methods. See
`TODO.md` for the follow-up plan.
