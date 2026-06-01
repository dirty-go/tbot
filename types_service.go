package tbot

// VideoChatScheduled is a service message about a video chat scheduled in the chat.
type VideoChatScheduled struct {
	StartDate int64 `json:"start_date"`
}

// VideoChatStarted is a service message about a video chat started in the chat.
type VideoChatStarted struct {
	Duration int `json:"duration,omitempty"`
}

// VideoChatEnded is a service message about a video chat ended in the chat.
type VideoChatEnded struct {
	Duration int `json:"duration"`
}

// VideoChatParticipantsInvited is a service message about new members invited
// to a video chat.
type VideoChatParticipantsInvited struct {
	Users []User `json:"users"`
}

// UsersShared contains information about users shared with the bot via a
// KeyboardButtonRequestUsers button.
type UsersShared struct {
	RequestID int          `json:"request_id"`
	Users     []SharedUser `json:"users"`
}

// SharedUser is a single user from UsersShared. Telegram only populates
// FirstName/LastName/Username/Photo if the bot requested them.
type SharedUser struct {
	UserID    int64       `json:"user_id"`
	FirstName string      `json:"first_name,omitempty"`
	LastName  string      `json:"last_name,omitempty"`
	Username  string      `json:"username,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
}

// ChatShared contains information about a chat shared with the bot via a
// KeyboardButtonRequestChat button.
type ChatShared struct {
	RequestID int         `json:"request_id"`
	ChatID    int64       `json:"chat_id"`
	Title     string      `json:"title,omitempty"`
	Username  string      `json:"username,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
}

// WriteAccessAllowed is a service message indicating that the user allowed the
// bot to write messages.
type WriteAccessAllowed struct {
	FromRequest        bool   `json:"from_request,omitempty"`
	WebAppName         string `json:"web_app_name,omitempty"`
	FromAttachmentMenu bool   `json:"from_attachment_menu,omitempty"`
}

// ChatOwnerLeft is a service message indicating the chat owner has left.
type ChatOwnerLeft struct {
	NewOwner *User `json:"new_owner,omitempty"`
}

// ChatOwnerChanged is a service message indicating the chat owner has changed.
type ChatOwnerChanged struct {
	NewOwner User `json:"new_owner"`
}

// ManagedBotCreated is a service message: the user created a managed bot.
type ManagedBotCreated struct {
	Bot User `json:"bot"`
}

// ManagedBotUpdated is the update payload sent when a managed bot is created
// or its token/owner changes.
type ManagedBotUpdated struct {
	User User `json:"user"`
	Bot  User `json:"bot"`
}

// PaidMessagePriceChanged is a service message indicating the price for paid
// messages has changed in the chat.
type PaidMessagePriceChanged struct {
	PaidMessageStarCount int `json:"paid_message_star_count"`
}

// DirectMessagePriceChanged is a service message indicating the price for paid
// direct messages to a channel has changed.
type DirectMessagePriceChanged struct {
	AreDirectMessagesEnabled bool `json:"are_direct_messages_enabled"`
	DirectMessageStarCount   int  `json:"direct_message_star_count,omitempty"`
}

// DirectMessagesTopic is information about the direct messages chat topic that
// contains a message.
type DirectMessagesTopic struct {
	TopicID int64 `json:"topic_id"`
	User    *User `json:"user,omitempty"`
}

// PollOptionAdded is a service message: an answer option was added to a poll.
type PollOptionAdded struct {
	PollMessage         *MaybeInaccessibleMessage `json:"poll_message,omitempty"`
	OptionPersistentID  string                    `json:"option_persistent_id"`
	OptionText          string                    `json:"option_text"`
	OptionTextEntities  []MessageEntity           `json:"option_text_entities,omitempty"`
}

// PollOptionDeleted is a service message: an answer option was deleted from a
// poll.
type PollOptionDeleted struct {
	PollMessage         *MaybeInaccessibleMessage `json:"poll_message,omitempty"`
	OptionPersistentID  string                    `json:"option_persistent_id"`
	OptionText          string                    `json:"option_text"`
	OptionTextEntities  []MessageEntity           `json:"option_text_entities,omitempty"`
}

// ChecklistTask is a single task within a Checklist.
type ChecklistTask struct {
	ID               int             `json:"id"`
	Text             string          `json:"text"`
	TextEntities     []MessageEntity `json:"text_entities,omitempty"`
	CompletedByUser  *User           `json:"completed_by_user,omitempty"`
	CompletedByChat  *Chat           `json:"completed_by_chat,omitempty"`
	CompletionDate   int64           `json:"completion_date,omitempty"`
}

// Checklist represents a checklist message.
type Checklist struct {
	Title                     string          `json:"title"`
	TitleEntities             []MessageEntity `json:"title_entities,omitempty"`
	Tasks                     []ChecklistTask `json:"tasks"`
	OthersCanAddTasks         bool            `json:"others_can_add_tasks,omitempty"`
	OthersCanMarkTasksAsDone  bool            `json:"others_can_mark_tasks_as_done,omitempty"`
}

// ChecklistTasksDone is a service message indicating that some tasks in a
// checklist were marked done or not done.
type ChecklistTasksDone struct {
	ChecklistMessage      *Message `json:"checklist_message,omitempty"`
	MarkedAsDoneTaskIDs   []int    `json:"marked_as_done_task_ids,omitempty"`
	MarkedAsNotDoneTaskIDs []int   `json:"marked_as_not_done_task_ids,omitempty"`
}

// ChecklistTasksAdded is a service message indicating that tasks were added to
// a checklist.
type ChecklistTasksAdded struct {
	ChecklistMessage *Message        `json:"checklist_message,omitempty"`
	Tasks            []ChecklistTask `json:"tasks"`
}

// SuggestedPostInfo describes the suggested-post parameters of a message.
type SuggestedPostInfo struct {
	State    string              `json:"state"`
	Price    *SuggestedPostPrice `json:"price,omitempty"`
	SendDate int64               `json:"send_date,omitempty"`
}

// SuggestedPostPrice represents the price of a suggested post. Currency is
// either "XTR" (Telegram Stars) or "TON" (toncoin nanotokens).
type SuggestedPostPrice struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount,omitempty"`
}

// SuggestedPostApproved is a service message: a suggested post was approved.
type SuggestedPostApproved struct {
	SuggestedPostMessage *Message            `json:"suggested_post_message,omitempty"`
	Price                *SuggestedPostPrice `json:"price,omitempty"`
	SendDate             int64               `json:"send_date"`
}

// SuggestedPostApprovalFailed is a service message: approval of a suggested
// post failed (e.g. the channel had insufficient stars at send time).
type SuggestedPostApprovalFailed struct {
	SuggestedPostMessage *Message           `json:"suggested_post_message,omitempty"`
	Price                SuggestedPostPrice `json:"price"`
}

// SuggestedPostDeclined is a service message: a suggested post was declined.
type SuggestedPostDeclined struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`
	Comment              string   `json:"comment,omitempty"`
}

// SuggestedPostPaid is a service message: payment for a suggested post was received.
type SuggestedPostPaid struct {
	SuggestedPostMessage *Message    `json:"suggested_post_message,omitempty"`
	Currency             string      `json:"currency"`
	Amount               int64       `json:"amount,omitempty"`
	StarAmount           *StarAmount `json:"star_amount,omitempty"`
}

// SuggestedPostRefunded is a service message: payment for a suggested post was refunded.
type SuggestedPostRefunded struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`
	Reason               string   `json:"reason"`
}
