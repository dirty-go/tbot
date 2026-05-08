package tbot

// ForumTopic represents a forum topic.
type ForumTopic struct {
	MessageThreadID   int    `json:"message_thread_id"`
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
	IsNameImplicit    bool   `json:"is_name_implicit,omitempty"`
}

// ForumTopicCreated is the service message: forum topic created.
type ForumTopicCreated struct {
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
	IsNameImplicit    bool   `json:"is_name_implicit,omitempty"`
}

// ForumTopicEdited is the service message: forum topic edited. Both fields
// are optional; only the changed ones are present.
type ForumTopicEdited struct {
	Name              string `json:"name,omitempty"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicClosed is a placeholder service message — no fields.
type ForumTopicClosed struct{}

// ForumTopicReopened is a placeholder service message — no fields.
type ForumTopicReopened struct{}

// GeneralForumTopicHidden is a placeholder service message — no fields.
type GeneralForumTopicHidden struct{}

// GeneralForumTopicUnhidden is a placeholder service message — no fields.
type GeneralForumTopicUnhidden struct{}
