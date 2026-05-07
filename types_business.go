package tbot

// BusinessConnection describes the connection of the bot with a business
// account. Sent in Update.BusinessConnection when the connection is created
// or revoked.
type BusinessConnection struct {
	ID         string             `json:"id"`
	User       User               `json:"user"`
	UserChatID int64              `json:"user_chat_id"`
	Date       int64              `json:"date"`
	Rights     *BusinessBotRights `json:"rights,omitempty"`
	IsEnabled  bool               `json:"is_enabled"`
}

// BusinessBotRights describes the rights of a business bot.
type BusinessBotRights struct {
	CanReply                  bool `json:"can_reply,omitempty"`
	CanReadMessages           bool `json:"can_read_messages,omitempty"`
	CanDeleteSentMessages     bool `json:"can_delete_sent_messages,omitempty"`
	CanDeleteAllMessages      bool `json:"can_delete_all_messages,omitempty"`
	CanEditName               bool `json:"can_edit_name,omitempty"`
	CanEditBio                bool `json:"can_edit_bio,omitempty"`
	CanEditProfilePhoto       bool `json:"can_edit_profile_photo,omitempty"`
	CanEditUsername           bool `json:"can_edit_username,omitempty"`
	CanChangeGiftSettings     bool `json:"can_change_gift_settings,omitempty"`
	CanViewGiftsAndStars      bool `json:"can_view_gifts_and_stars,omitempty"`
	CanConvertGiftsToStars    bool `json:"can_convert_gifts_to_stars,omitempty"`
	CanTransferAndUpgradeGifts bool `json:"can_transfer_and_upgrade_gifts,omitempty"`
	CanTransferStars          bool `json:"can_transfer_stars,omitempty"`
	CanManageStories          bool `json:"can_manage_stories,omitempty"`
}

// BusinessIntro is the introduction shown on a business account profile.
type BusinessIntro struct {
	Title   string   `json:"title,omitempty"`
	Message string   `json:"message,omitempty"`
	Sticker *Sticker `json:"sticker,omitempty"`
}

// BusinessLocation is the location of a business account.
type BusinessLocation struct {
	Address  string    `json:"address"`
	Location *Location `json:"location,omitempty"`
}

// BusinessOpeningHoursInterval is one open interval in a business week.
//
// Telegram represents minutes as a sequence number from Monday 00:00. So
// Tuesday 14:00 = 1*1440 + 14*60 = 2280.
type BusinessOpeningHoursInterval struct {
	OpeningMinute int `json:"opening_minute"`
	ClosingMinute int `json:"closing_minute"`
}

// BusinessOpeningHours describes the opening hours of a business account.
type BusinessOpeningHours struct {
	TimeZoneName string                         `json:"time_zone_name"`
	OpeningHours []BusinessOpeningHoursInterval `json:"opening_hours"`
}

// BusinessMessagesDeleted is the update payload sent when messages were
// deleted from a connected business account.
type BusinessMessagesDeleted struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Chat                 Chat   `json:"chat"`
	MessageIDs           []int  `json:"message_ids"`
}
