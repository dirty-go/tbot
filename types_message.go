package tbot

import "encoding/json"

// Message represents a message sent in a chat. It is the central type returned
// by sendMessage, sendPhoto, etc., and embedded in many update payloads.
//
// Most fields are optional; only message_id, date and chat are guaranteed.
// Service-message fields (e.g. NewChatMembers, ForumTopicCreated, Giveaway,
// VideoChat*) are mutually exclusive with content fields.
//
// The integer chat / sender / message thread / forward identifiers are
// `int64` because Telegram explicitly notes they may exceed 32-bit precision.
type Message struct {
	MessageID                     int                            `json:"message_id"`
	MessageThreadID               int                            `json:"message_thread_id,omitempty"`
	DirectMessagesTopic           *DirectMessagesTopic           `json:"direct_messages_topic,omitempty"`
	From                          *User                          `json:"from,omitempty"`
	SenderChat                    *Chat                          `json:"sender_chat,omitempty"`
	SenderBoostCount              int                            `json:"sender_boost_count,omitempty"`
	SenderBusinessBot             *User                          `json:"sender_business_bot,omitempty"`
	SenderTag                     string                         `json:"sender_tag,omitempty"`
	Date                          int64                          `json:"date"`
	BusinessConnectionID          string                         `json:"business_connection_id,omitempty"`
	Chat                          Chat                           `json:"chat"`
	ForwardOrigin                 *MessageOrigin                 `json:"forward_origin,omitempty"`
	IsTopicMessage                bool                           `json:"is_topic_message,omitempty"`
	IsAutomaticForward            bool                           `json:"is_automatic_forward,omitempty"`
	ReplyToMessage                *Message                       `json:"reply_to_message,omitempty"`
	ExternalReply                 *ExternalReplyInfo             `json:"external_reply,omitempty"`
	Quote                         *TextQuote                     `json:"quote,omitempty"`
	ReplyToStory                  *Story                         `json:"reply_to_story,omitempty"`
	ReplyToChecklistTaskID        int                            `json:"reply_to_checklist_task_id,omitempty"`
	ReplyToPollOptionID           string                         `json:"reply_to_poll_option_id,omitempty"`
	ViaBot                        *User                          `json:"via_bot,omitempty"`
	EditDate                      int64                          `json:"edit_date,omitempty"`
	HasProtectedContent           bool                           `json:"has_protected_content,omitempty"`
	IsFromOffline                 bool                           `json:"is_from_offline,omitempty"`
	IsPaidPost                    bool                           `json:"is_paid_post,omitempty"`
	MediaGroupID                  string                         `json:"media_group_id,omitempty"`
	AuthorSignature               string                         `json:"author_signature,omitempty"`
	PaidStarCount                 int                            `json:"paid_star_count,omitempty"`
	Text                          string                         `json:"text,omitempty"`
	Entities                      []MessageEntity                `json:"entities,omitempty"`
	LinkPreviewOptions            *LinkPreviewOptions            `json:"link_preview_options,omitempty"`
	SuggestedPostInfo             *SuggestedPostInfo             `json:"suggested_post_info,omitempty"`
	EffectID                      string                         `json:"effect_id,omitempty"`
	Animation                     *Animation                     `json:"animation,omitempty"`
	Audio                         *Audio                         `json:"audio,omitempty"`
	Document                      *Document                      `json:"document,omitempty"`
	PaidMedia                     *PaidMediaInfo                 `json:"paid_media,omitempty"`
	Photo                         []PhotoSize                    `json:"photo,omitempty"`
	Sticker                       *Sticker                       `json:"sticker,omitempty"`
	Story                         *Story                         `json:"story,omitempty"`
	Video                         *Video                         `json:"video,omitempty"`
	VideoNote                     *VideoNote                     `json:"video_note,omitempty"`
	Voice                         *Voice                         `json:"voice,omitempty"`
	Caption                       string                         `json:"caption,omitempty"`
	CaptionEntities               []MessageEntity                `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia         bool                           `json:"show_caption_above_media,omitempty"`
	HasMediaSpoiler               bool                           `json:"has_media_spoiler,omitempty"`
	Checklist                     *Checklist                     `json:"checklist,omitempty"`
	Contact                       *Contact                       `json:"contact,omitempty"`
	Dice                          *Dice                          `json:"dice,omitempty"`
	Game                          *Game                          `json:"game,omitempty"`
	Poll                          *Poll                          `json:"poll,omitempty"`
	Venue                         *Venue                         `json:"venue,omitempty"`
	Location                      *Location                      `json:"location,omitempty"`
	NewChatMembers                []User                         `json:"new_chat_members,omitempty"`
	LeftChatMember                *User                          `json:"left_chat_member,omitempty"`
	ChatOwnerLeft                 *ChatOwnerLeft                 `json:"chat_owner_left,omitempty"`
	ChatOwnerChanged              *ChatOwnerChanged              `json:"chat_owner_changed,omitempty"`
	NewChatTitle                  string                         `json:"new_chat_title,omitempty"`
	NewChatPhoto                  []PhotoSize                    `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto               bool                           `json:"delete_chat_photo,omitempty"`
	GroupChatCreated              bool                           `json:"group_chat_created,omitempty"`
	SupergroupChatCreated         bool                           `json:"supergroup_chat_created,omitempty"`
	ChannelChatCreated            bool                           `json:"channel_chat_created,omitempty"`
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	MigrateToChatID               int64                          `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID             int64                          `json:"migrate_from_chat_id,omitempty"`
	PinnedMessage                 *MaybeInaccessibleMessage      `json:"pinned_message,omitempty"`
	Invoice                       *Invoice                       `json:"invoice,omitempty"`
	SuccessfulPayment             *SuccessfulPayment             `json:"successful_payment,omitempty"`
	RefundedPayment               *RefundedPayment               `json:"refunded_payment,omitempty"`
	UsersShared                   *UsersShared                   `json:"users_shared,omitempty"`
	ChatShared                    *ChatShared                    `json:"chat_shared,omitempty"`
	Gift                          *GiftInfo                      `json:"gift,omitempty"`
	UniqueGift                    *UniqueGiftInfo                `json:"unique_gift,omitempty"`
	GiftUpgradeSent               *GiftInfo                      `json:"gift_upgrade_sent,omitempty"`
	ConnectedWebsite              string                         `json:"connected_website,omitempty"`
	WriteAccessAllowed            *WriteAccessAllowed            `json:"write_access_allowed,omitempty"`
	PassportData                  *PassportData                  `json:"passport_data,omitempty"`
	ProximityAlertTriggered       *ProximityAlertTriggered       `json:"proximity_alert_triggered,omitempty"`
	BoostAdded                    *ChatBoostAdded                `json:"boost_added,omitempty"`
	ChatBackgroundSet             *ChatBackground                `json:"chat_background_set,omitempty"`
	ChecklistTasksDone            *ChecklistTasksDone            `json:"checklist_tasks_done,omitempty"`
	ChecklistTasksAdded           *ChecklistTasksAdded           `json:"checklist_tasks_added,omitempty"`
	DirectMessagePriceChanged     *DirectMessagePriceChanged     `json:"direct_message_price_changed,omitempty"`
	ForumTopicCreated             *ForumTopicCreated             `json:"forum_topic_created,omitempty"`
	ForumTopicEdited              *ForumTopicEdited              `json:"forum_topic_edited,omitempty"`
	ForumTopicClosed              *ForumTopicClosed              `json:"forum_topic_closed,omitempty"`
	ForumTopicReopened            *ForumTopicReopened            `json:"forum_topic_reopened,omitempty"`
	GeneralForumTopicHidden       *GeneralForumTopicHidden       `json:"general_forum_topic_hidden,omitempty"`
	GeneralForumTopicUnhidden     *GeneralForumTopicUnhidden     `json:"general_forum_topic_unhidden,omitempty"`
	GiveawayCreated               *GiveawayCreated               `json:"giveaway_created,omitempty"`
	Giveaway                      *Giveaway                      `json:"giveaway,omitempty"`
	GiveawayWinners               *GiveawayWinners               `json:"giveaway_winners,omitempty"`
	GiveawayCompleted             *GiveawayCompleted             `json:"giveaway_completed,omitempty"`
	ManagedBotCreated             *ManagedBotCreated             `json:"managed_bot_created,omitempty"`
	PaidMessagePriceChanged       *PaidMessagePriceChanged       `json:"paid_message_price_changed,omitempty"`
	PollOptionAdded               *PollOptionAdded               `json:"poll_option_added,omitempty"`
	PollOptionDeleted             *PollOptionDeleted             `json:"poll_option_deleted,omitempty"`
	SuggestedPostApproved         *SuggestedPostApproved         `json:"suggested_post_approved,omitempty"`
	SuggestedPostApprovalFailed   *SuggestedPostApprovalFailed   `json:"suggested_post_approval_failed,omitempty"`
	SuggestedPostDeclined         *SuggestedPostDeclined         `json:"suggested_post_declined,omitempty"`
	SuggestedPostPaid             *SuggestedPostPaid             `json:"suggested_post_paid,omitempty"`
	SuggestedPostRefunded         *SuggestedPostRefunded         `json:"suggested_post_refunded,omitempty"`
	VideoChatScheduled            *VideoChatScheduled            `json:"video_chat_scheduled,omitempty"`
	VideoChatStarted              *VideoChatStarted              `json:"video_chat_started,omitempty"`
	VideoChatEnded                *VideoChatEnded                `json:"video_chat_ended,omitempty"`
	VideoChatParticipantsInvited  *VideoChatParticipantsInvited  `json:"video_chat_participants_invited,omitempty"`
	WebAppData                    *WebAppData                    `json:"web_app_data,omitempty"`
	ReplyMarkup                   *InlineKeyboardMarkup          `json:"reply_markup,omitempty"`

	// Deprecated: replaced by ForwardOrigin in Bot API 7.0. Telegram no longer
	// populates these fields; use ForwardOrigin and switch on its Type.
	ForwardFrom          *User  `json:"forward_from,omitempty"`
	ForwardFromChat      *Chat  `json:"forward_from_chat,omitempty"`
	ForwardFromMessageID int64  `json:"forward_from_message_id,omitempty"`
	ForwardSignature     string `json:"forward_signature,omitempty"`
	ForwardSenderName    string `json:"forward_sender_name,omitempty"`
	ForwardDate          int64  `json:"forward_date,omitempty"`

	// Deprecated: replaced by VideoChat* equivalents in Bot API 6.0.
	VoiceChatScheduled           *VoiceChatScheduled           `json:"voice_chat_scheduled,omitempty"`
	VoiceChatStarted             *VoiceChatStarted             `json:"voice_chat_started,omitempty"`
	VoiceChatEnded               *VoiceChatEnded               `json:"voice_chat_ended,omitempty"`
	VoiceChatParticipantsInvited *VoiceChatParticipantsInvited `json:"voice_chat_participants_invited,omitempty"`
}

// MessageId is the response shape returned by methods such as copyMessage that
// echo only the new message identifier.
type MessageId struct {
	MessageID int `json:"message_id"`
}

// InaccessibleMessage represents a message that cannot be accessed by the bot
// any more. The Date field is always 0 — use it to distinguish from Message.
type InaccessibleMessage struct {
	Chat      Chat  `json:"chat"`
	MessageID int   `json:"message_id"`
	Date      int64 `json:"date"`
}

// MaybeInaccessibleMessage is either a Message or an InaccessibleMessage. The
// concrete variant is determined by the Date field: 0 means inaccessible.
type MaybeInaccessibleMessage struct {
	Message
}

// IsInaccessible reports whether the message is an InaccessibleMessage — one
// the bot can no longer access. Such payloads always carry a zero Date.
func (m MaybeInaccessibleMessage) IsInaccessible() bool { return m.Date == 0 }

// TextQuote contains information about the quoted part of a message that is
// replied to by the given message.
type TextQuote struct {
	Text     string          `json:"text"`
	Entities []MessageEntity `json:"entities,omitempty"`
	Position int             `json:"position"`
	IsManual bool            `json:"is_manual,omitempty"`
}

// ExternalReplyInfo contains information about a message replied to from
// another chat or topic.
type ExternalReplyInfo struct {
	Origin             MessageOrigin       `json:"origin"`
	Chat               *Chat               `json:"chat,omitempty"`
	MessageID          int                 `json:"message_id,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	Animation          *Animation          `json:"animation,omitempty"`
	Audio              *Audio              `json:"audio,omitempty"`
	Document           *Document           `json:"document,omitempty"`
	PaidMedia          *PaidMediaInfo      `json:"paid_media,omitempty"`
	Photo              []PhotoSize         `json:"photo,omitempty"`
	Sticker            *Sticker            `json:"sticker,omitempty"`
	Story              *Story              `json:"story,omitempty"`
	Video              *Video              `json:"video,omitempty"`
	VideoNote          *VideoNote          `json:"video_note,omitempty"`
	Voice              *Voice              `json:"voice,omitempty"`
	HasMediaSpoiler    bool                `json:"has_media_spoiler,omitempty"`
	Checklist          *Checklist          `json:"checklist,omitempty"`
	Contact            *Contact            `json:"contact,omitempty"`
	Dice               *Dice               `json:"dice,omitempty"`
	Game               *Game               `json:"game,omitempty"`
	Giveaway           *Giveaway           `json:"giveaway,omitempty"`
	GiveawayWinners    *GiveawayWinners    `json:"giveaway_winners,omitempty"`
	Invoice            *Invoice            `json:"invoice,omitempty"`
	Location           *Location           `json:"location,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	Venue              *Venue              `json:"venue,omitempty"`
}

// ReplyParameters describes the message being replied to. Used as a value of
// the reply_parameters request field for sendMessage and friends.
type ReplyParameters struct {
	MessageID                int             `json:"message_id"`
	ChatID                   any             `json:"chat_id,omitempty"`
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply,omitempty"`
	Quote                    string          `json:"quote,omitempty"`
	QuoteParseMode           string          `json:"quote_parse_mode,omitempty"`
	QuoteEntities            []MessageEntity `json:"quote_entities,omitempty"`
	QuotePosition            int             `json:"quote_position,omitempty"`
	ChecklistTaskID          int             `json:"checklist_task_id,omitempty"`
	PollOptionID             string          `json:"poll_option_id,omitempty"`
}

// MessageOrigin describes where a forwarded message came from. The variant is
// indicated by the Type field, which is one of "user", "hidden_user", "chat"
// or "channel".
type MessageOrigin struct {
	Type            string `json:"type"`
	Date            int64  `json:"date"`
	SenderUser      *User  `json:"sender_user,omitempty"`
	SenderUserName  string `json:"sender_user_name,omitempty"`
	SenderChat      *Chat  `json:"sender_chat,omitempty"`
	Chat            *Chat  `json:"chat,omitempty"`
	MessageID       int    `json:"message_id,omitempty"`
	AuthorSignature string `json:"author_signature,omitempty"`
}

// MarshalJSON emits only fields valid for the selected message-origin variant.
func (m MessageOrigin) MarshalJSON() ([]byte, error) {
	type payload struct {
		Type            string `json:"type"`
		Date            int64  `json:"date"`
		SenderUser      *User  `json:"sender_user,omitempty"`
		SenderUserName  string `json:"sender_user_name,omitempty"`
		SenderChat      *Chat  `json:"sender_chat,omitempty"`
		Chat            *Chat  `json:"chat,omitempty"`
		MessageID       int    `json:"message_id,omitempty"`
		AuthorSignature string `json:"author_signature,omitempty"`
	}
	out := payload{
		Type: m.Type,
		Date: m.Date,
	}
	switch m.Type {
	case MessageOriginTypeUser:
		out.SenderUser = m.SenderUser
	case MessageOriginTypeHiddenUser:
		out.SenderUserName = m.SenderUserName
	case MessageOriginTypeChat:
		out.SenderChat = m.SenderChat
		out.AuthorSignature = m.AuthorSignature
	case MessageOriginTypeChannel:
		out.Chat = m.Chat
		out.MessageID = m.MessageID
		out.AuthorSignature = m.AuthorSignature
	}
	return json.Marshal(out)
}

// MessageOrigin Type values.
const (
	MessageOriginTypeUser       = "user"
	MessageOriginTypeHiddenUser = "hidden_user"
	MessageOriginTypeChat       = "chat"
	MessageOriginTypeChannel    = "channel"
)

// LinkPreviewOptions describes link preview generation for a message.
type LinkPreviewOptions struct {
	IsDisabled       bool   `json:"is_disabled,omitempty"`
	URL              string `json:"url,omitempty"`
	PreferSmallMedia bool   `json:"prefer_small_media,omitempty"`
	PreferLargeMedia bool   `json:"prefer_large_media,omitempty"`
	ShowAboveText    bool   `json:"show_above_text,omitempty"`
}

// MessageAutoDeleteTimerChanged is a service message about the auto-delete
// timer being changed.
type MessageAutoDeleteTimerChanged struct {
	MessageAutoDeleteTime int `json:"message_auto_delete_time"`
}

// ProximityAlertTriggered is a service message sent when a user triggers a
// proximity alert set by another user.
type ProximityAlertTriggered struct {
	Traveler User `json:"traveler"`
	Watcher  User `json:"watcher"`
	Distance int  `json:"distance"`
}
