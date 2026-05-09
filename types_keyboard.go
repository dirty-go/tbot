package tbot

import "encoding/json"

// KeyboardButton represents one button of the reply (custom) keyboard.
//
// Optional fields are mutually exclusive: at most one of RequestUsers,
// RequestChat, RequestManagedBot, RequestContact, RequestLocation,
// RequestPoll or WebApp may be set.
type KeyboardButton struct {
	Text               string                          `json:"text"`
	IconCustomEmojiID  string                          `json:"icon_custom_emoji_id,omitempty"`
	Style              string                          `json:"style,omitempty"`
	RequestUsers       *KeyboardButtonRequestUsers     `json:"request_users,omitempty"`
	RequestChat        *KeyboardButtonRequestChat      `json:"request_chat,omitempty"`
	RequestManagedBot  *KeyboardButtonRequestManagedBot `json:"request_managed_bot,omitempty"`
	RequestContact     bool                            `json:"request_contact,omitempty"`
	RequestLocation    bool                            `json:"request_location,omitempty"`
	RequestPoll        *KeyboardButtonPollType         `json:"request_poll,omitempty"`
	WebApp             *WebAppInfo                     `json:"web_app,omitempty"`

	// Deprecated: misspelled in early versions of this library. Telegram's
	// JSON tag is "request_poll"; this field is kept so existing Go code that
	// referenced the typo still compiles. New code should use RequestPoll.
	RequestPool *KeyboardButtonPollType `json:"-"`
}

// KeyboardButtonRequestUsers requests one or more users to be shared with the
// bot.
type KeyboardButtonRequestUsers struct {
	RequestID       int   `json:"request_id"`
	UserIsBot       *bool `json:"user_is_bot,omitempty"`
	UserIsPremium   *bool `json:"user_is_premium,omitempty"`
	MaxQuantity     int   `json:"max_quantity,omitempty"`
	RequestName     bool  `json:"request_name,omitempty"`
	RequestUsername bool  `json:"request_username,omitempty"`
	RequestPhoto    bool  `json:"request_photo,omitempty"`
}

// KeyboardButtonRequestChat requests a chat to be shared with the bot.
type KeyboardButtonRequestChat struct {
	RequestID               int                      `json:"request_id"`
	ChatIsChannel           bool                     `json:"chat_is_channel"`
	ChatIsForum             *bool                    `json:"chat_is_forum,omitempty"`
	ChatHasUsername         *bool                    `json:"chat_has_username,omitempty"`
	ChatIsCreated           bool                     `json:"chat_is_created,omitempty"`
	UserAdministratorRights *ChatAdministratorRights `json:"user_administrator_rights,omitempty"`
	BotAdministratorRights  *ChatAdministratorRights `json:"bot_administrator_rights,omitempty"`
	BotIsMember             bool                     `json:"bot_is_member,omitempty"`
	RequestTitle            bool                     `json:"request_title,omitempty"`
	RequestUsername         bool                     `json:"request_username,omitempty"`
	RequestPhoto            bool                     `json:"request_photo,omitempty"`
}

// KeyboardButtonRequestManagedBot requests the user to create a bot to be
// managed by the bot that sent the request.
type KeyboardButtonRequestManagedBot struct {
	RequestID         int    `json:"request_id"`
	SuggestedName     string `json:"suggested_name,omitempty"`
	SuggestedUsername string `json:"suggested_username,omitempty"`
}

// KeyboardButtonPollType is the request_poll value of a KeyboardButton. Empty
// Type allows any poll; "regular" or "quiz" restricts it.
type KeyboardButtonPollType struct {
	Type string `json:"type,omitempty"`
}

// KeyboardButtonPoolType is the misspelled alias kept for source
// compatibility.
//
// Deprecated: use KeyboardButtonPollType.
type KeyboardButtonPoolType = KeyboardButtonPollType

// ReplyKeyboardMarkup describes a custom keyboard with reply options.
//
// Note: Keyboard is a 2D array — each inner slice is a row of buttons.
type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	IsPersistent          bool               `json:"is_persistent,omitempty"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string             `json:"input_field_placeholder,omitempty"`
	Selective             bool               `json:"selective,omitempty"`
}

// ReplyKeyboardRemove removes the custom keyboard for the recipient.
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective,omitempty"`
}

// InlineKeyboardMarkup describes an inline keyboard attached to a message.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton describes one button of an inline keyboard.
//
// At most one of URL, CallbackData, WebApp, LoginURL, SwitchInlineQuery,
// SwitchInlineQueryCurrentChat, SwitchInlineQueryChosenChat, CopyText,
// CallbackGame or Pay should be set.
type InlineKeyboardButton struct {
	Text                         string                       `json:"text"`
	IconCustomEmojiID            string                       `json:"icon_custom_emoji_id,omitempty"`
	Style                        string                       `json:"style,omitempty"`
	URL                          string                       `json:"url,omitempty"`
	CallbackData                 string                       `json:"callback_data,omitempty"`
	WebApp                       *WebAppInfo                  `json:"web_app,omitempty"`
	LoginURL                     *LoginURL                    `json:"login_url,omitempty"`
	SwitchInlineQuery            *string                      `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat *string                      `json:"switch_inline_query_current_chat,omitempty"`
	SwitchInlineQueryChosenChat  *SwitchInlineQueryChosenChat `json:"switch_inline_query_chosen_chat,omitempty"`
	CopyText                     *CopyTextButton              `json:"copy_text,omitempty"`
	CallbackGame                 *CallbackGame                `json:"callback_game,omitempty"`
	Pay                          bool                         `json:"pay,omitempty"`
}

// LoginURL describes the parameters of a Telegram-Login button.
type LoginURL struct {
	URL                string `json:"url"`
	ForwardText        string `json:"forward_text,omitempty"`
	BotUsername        string `json:"bot_username,omitempty"`
	RequestWriteAccess bool   `json:"request_write_access,omitempty"`
}

// SwitchInlineQueryChosenChat lets the user pick a chat in which to
// auto-fill an inline query when an inline-keyboard button is tapped.
type SwitchInlineQueryChosenChat struct {
	Query             string `json:"query,omitempty"`
	AllowUserChats    bool   `json:"allow_user_chats,omitempty"`
	AllowBotChats     bool   `json:"allow_bot_chats,omitempty"`
	AllowGroupChats   bool   `json:"allow_group_chats,omitempty"`
	AllowChannelChats bool   `json:"allow_channel_chats,omitempty"`
}

// CopyTextButton represents a button that copies a given text to the clipboard.
type CopyTextButton struct {
	Text string `json:"text"`
}

// ForceReply requests the client to display a reply interface for the user.
type ForceReply struct {
	ForceReply            bool   `json:"force_reply"`
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	Selective             bool   `json:"selective,omitempty"`
}

// CallbackQuery represents a callback query from a callback button in an
// inline keyboard.
type CallbackQuery struct {
	ID              string                    `json:"id"`
	From            User                      `json:"from"`
	Message         *MaybeInaccessibleMessage `json:"message,omitempty"`
	InlineMessageID string                    `json:"inline_message_id,omitempty"`
	ChatInstance    string                    `json:"chat_instance"`
	Data            string                    `json:"data,omitempty"`
	GameShortName   string                    `json:"game_short_name,omitempty"`
}

// WebAppInfo describes a Web App.
type WebAppInfo struct {
	URL string `json:"url"`
}

// WebAppData is the data sent from a Web App to the bot.
type WebAppData struct {
	Data       string `json:"data"`
	ButtonText string `json:"button_text"`
}

// MenuButton describes the bot's menu button in a private chat. Variant is
// given by Type: "default", "commands" or "web_app".
type MenuButton struct {
	Type   string      `json:"type"`
	Text   string      `json:"text,omitempty"`
	WebApp *WebAppInfo `json:"web_app,omitempty"`
}

// MarshalJSON emits only fields valid for the selected menu-button variant.
func (m MenuButton) MarshalJSON() ([]byte, error) {
	type payload struct {
		Type   string      `json:"type"`
		Text   string      `json:"text,omitempty"`
		WebApp *WebAppInfo `json:"web_app,omitempty"`
	}
	out := payload{Type: m.Type}
	if m.Type == MenuButtonTypeWebApp {
		out.Text = m.Text
		out.WebApp = m.WebApp
	}
	return json.Marshal(out)
}

// MenuButton Type values.
const (
	MenuButtonTypeDefault  = "default"
	MenuButtonTypeCommands = "commands"
	MenuButtonTypeWebApp   = "web_app"
)
