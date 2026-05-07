package tbot

// Update represents an incoming update. At most one of the optional fields
// will be set per Update — the receiver should switch on whichever is non-nil.
type Update struct {
	UpdateID                 int                          `json:"update_id"`
	Message                  *Message                     `json:"message,omitempty"`
	EditedMessage            *Message                     `json:"edited_message,omitempty"`
	ChannelPost              *Message                     `json:"channel_post,omitempty"`
	EditedChannelPost        *Message                     `json:"edited_channel_post,omitempty"`
	BusinessConnection       *BusinessConnection          `json:"business_connection,omitempty"`
	BusinessMessage          *Message                     `json:"business_message,omitempty"`
	EditedBusinessMessage    *Message                     `json:"edited_business_message,omitempty"`
	DeletedBusinessMessages  *BusinessMessagesDeleted     `json:"deleted_business_messages,omitempty"`
	MessageReaction          *MessageReactionUpdated      `json:"message_reaction,omitempty"`
	MessageReactionCount     *MessageReactionCountUpdated `json:"message_reaction_count,omitempty"`
	InlineQuery              *InlineQuery                 `json:"inline_query,omitempty"`
	ChosenInlineResult       *ChosenInlineResult          `json:"chosen_inline_result,omitempty"`
	CallbackQuery            *CallbackQuery               `json:"callback_query,omitempty"`
	ShippingQuery            *ShippingQuery               `json:"shipping_query,omitempty"`
	PreCheckoutQuery         *PreCheckoutQuery            `json:"pre_checkout_query,omitempty"`
	PurchasedPaidMedia       *PaidMediaPurchased          `json:"purchased_paid_media,omitempty"`
	Poll                     *Poll                        `json:"poll,omitempty"`
	PollAnswer               *PollAnswer                  `json:"poll_answer,omitempty"`
	MyChatMember             *ChatMemberUpdated           `json:"my_chat_member,omitempty"`
	ChatMember               *ChatMemberUpdated           `json:"chat_member,omitempty"`
	ChatJoinRequest          *ChatJoinRequest             `json:"chat_join_request,omitempty"`
	ChatBoost                *ChatBoostUpdated            `json:"chat_boost,omitempty"`
	RemovedChatBoost         *ChatBoostRemoved            `json:"removed_chat_boost,omitempty"`
	ManagedBot               *ManagedBotUpdated           `json:"managed_bot,omitempty"`
}

// WebhookInfo describes the current webhook status.
type WebhookInfo struct {
	URL                          string   `json:"url"`
	HasCustomCertificate         bool     `json:"has_custom_certificate"`
	PendingUpdateCount           int      `json:"pending_update_count"`
	IPAddress                    string   `json:"ip_address,omitempty"`
	LastErrorDate                int64    `json:"last_error_date,omitempty"`
	LastErrorMessage             string   `json:"last_error_message,omitempty"`
	LastSynchronizationErrorDate int64    `json:"last_synchronization_error_date,omitempty"`
	MaxConnections               int      `json:"max_connections,omitempty"`
	AllowedUpdates               []string `json:"allowed_updates,omitempty"`
}

// ResponseParameters carries optional metadata returned with API errors —
// most importantly retry_after on flood control and migrate_to_chat_id when a
// group has been upgraded to a supergroup.
type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`
	RetryAfter      int   `json:"retry_after,omitempty"`
}
