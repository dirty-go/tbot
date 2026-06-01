package tbot

import (
	"encoding/json"
	"strconv"
)

// ChatID identifies a target chat. The Bot API accepts either a numeric
// identifier or an @-prefixed channel username, so this type carries
// whichever the caller supplied and renders it correctly for both
// form-encoded request fields (via String) and JSON request bodies (via
// MarshalJSON).
//
// Construct one with Int64 or Username; the zero value is invalid and reports
// true from IsZero.
type ChatID struct {
	id       int64
	username string
}

// Int64 builds a ChatID from a numeric chat identifier.
func Int64(id int64) ChatID { return ChatID{id: id} }

// Username builds a ChatID from a public chat's @username. The leading "@" is
// preserved as given; callers that omit it get whatever Telegram accepts for
// that field.
func Username(username string) ChatID { return ChatID{username: username} }

// IsZero reports whether c was never assigned a chat (the zero value).
func (c ChatID) IsZero() bool { return c.id == 0 && c.username == "" }

// String renders the chat identifier for a form-encoded request field: the
// username when set, otherwise the decimal numeric id.
func (c ChatID) String() string {
	if c.username != "" {
		return c.username
	}
	return strconv.FormatInt(c.id, 10)
}

// MarshalJSON emits a JSON string for a username and a JSON number for a
// numeric id, matching the chat_id shape the Bot API documents.
func (c ChatID) MarshalJSON() ([]byte, error) {
	if c.username != "" {
		return json.Marshal(c.username)
	}
	return json.Marshal(c.id)
}
