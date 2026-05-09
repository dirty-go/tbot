package tbot

import "encoding/json"

// ReactionType describes a reaction. The variant is given by Type:
//   - "emoji"        — Emoji is set
//   - "custom_emoji" — CustomEmojiID is set
//   - "paid"         — neither, the reaction is paid
type ReactionType struct {
	Type          string `json:"type"`
	Emoji         string `json:"emoji,omitempty"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

// MarshalJSON emits only fields valid for the current reaction variant.
func (r ReactionType) MarshalJSON() ([]byte, error) {
	type payload struct {
		Type          string `json:"type"`
		Emoji         string `json:"emoji,omitempty"`
		CustomEmojiID string `json:"custom_emoji_id,omitempty"`
	}
	out := payload{Type: r.Type}
	switch r.Type {
	case ReactionTypeEmoji:
		out.Emoji = r.Emoji
	case ReactionTypeCustomEmoji:
		out.CustomEmojiID = r.CustomEmojiID
	}
	return json.Marshal(out)
}

// ReactionType Type values.
const (
	ReactionTypeEmoji       = "emoji"
	ReactionTypeCustomEmoji = "custom_emoji"
	ReactionTypePaid        = "paid"
)

// ReactionCount is the count of one reaction on a message.
type ReactionCount struct {
	Type       ReactionType `json:"type"`
	TotalCount int          `json:"total_count"`
}

// MessageReactionUpdated is sent when a reaction to a message changes.
type MessageReactionUpdated struct {
	Chat        Chat           `json:"chat"`
	MessageID   int            `json:"message_id"`
	User        *User          `json:"user,omitempty"`
	ActorChat   *Chat          `json:"actor_chat,omitempty"`
	Date        int64          `json:"date"`
	OldReaction []ReactionType `json:"old_reaction"`
	NewReaction []ReactionType `json:"new_reaction"`
}

// MessageReactionCountUpdated is sent when reactions on an anonymously-reacted
// message change.
type MessageReactionCountUpdated struct {
	Chat      Chat            `json:"chat"`
	MessageID int             `json:"message_id"`
	Date      int64           `json:"date"`
	Reactions []ReactionCount `json:"reactions"`
}
