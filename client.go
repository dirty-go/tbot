package tbot

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	token        string
	baseURL      string
	url          string
	updateParams url.Values
	timeout      int
	bufferSize   int
	nextOffset   int
	logger       Logger
}

type sendOption func(url.Values)

var (
	OptParseModeHTML       = func(r url.Values) { r.Set("parse_mode", "HTML") }
	OptParseModeMarkdown   = func(r url.Values) { r.Set("parse_mode", "MarkdownV2") }
	OptDisableNotification = func(r url.Values) { r.Set("disable_notification", "true") }
	OptReplyToMessageID    = func(id int) sendOption {
		return func(r url.Values) {
			r.Set("reply_to_message_id", strconv.Itoa(id))
		}
	}
	OptSendingWithoutReply = func(r url.Values) { r.Set("allow_sending_without_reply", "true") }
)

func NewClient(token string, baseURL string) *Client {
	if baseURL == "" {
		baseURL = apiBaseURL
	}
	return &Client{
		token:   token,
		baseURL: baseURL,
		url:     fmt.Sprintf("%s/bot%s", baseURL, token) + "%s",
	}
}

func structString(s any) string {
	str, _ := json.Marshal(s)
	return string(str)
}

type User struct {
	ID                      int    `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Username                string `json:"username"`
	LanguageCode            string `json:"language_code"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
}

// Me returns info about bot as a User object
func (c *Client) Me() (*User, error) {
	var me User
	err := c.sendRequest("/getMe", nil, &me)
	return &me, err
}

// ChatPhoto represents a chat photo
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

// ChatPermissions describes actions that a non-administrator user is allowed to take in a chat.
type ChatPermissions struct {
	// CanSendMessages True, if the user is allowed to send text messages, contacts, locations and venues
	CanSendMessages bool `json:"can_send_messages,omitempty"`
	// CanSendMediaMessages True, if the user is allowed to send audios, documents, photos, videos, video notes and voice notes, implies can_send_messages
	CanSendMediaMessages bool `json:"can_send_media_messages,omitempty"`
	// CanSendPolls True, if the user is allowed to send polls, implies can_send_messages
	CanSendPolls bool `json:"can_send_polls,omitempty"`
	// CanSendOtherMessages True, if the user is allowed to send animations, games, stickers and use inline bots, implies can_send_media_messages
	CanSendOtherMessages bool `json:"can_send_other_messages,omitempty"`
	// CanAddWebPagePreviews True, if the user is allowed to add web page previews to their messages, implies can_send_media_messages
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	// CanChaneInfo True, if the user is allowed to change the chat title, photo and other settings. Ignored in public supergroups
	CanChangeInfo bool `json:"can_change_info,omitempty"`
	// CanInviteUsers  True, if the user is allowed to invite new users to the chat
	CanInviteUsers bool `json:"can_invite_users,omitempty"`
	// CanPinMessages True, if the user is allowed to pin messages. Ignored in public supergroups
	CanPinMessages bool `json:"can_pin_messages,omitempty"`
}

type replyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

type LoginURL struct {
	URL                string  `json:"url"`
	ForwardText        *string `json:"forward_text,omitempty"`
	BotUsername        string  `json:"bot_username,omitempty"`
	RequestWriteAccess string  `json:"request_write_access,omitempty"`
}

type forceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"selective"`
}

var (
	OptDisableWebPagePreview = func(r url.Values) { r.Set("disable_web_page_preview", "true") }
	OptReplyKeyboardRemove   = func(r url.Values) { r.Set("reply_markup", structString(&replyKeyboardRemove{RemoveKeyboard: true})) }
	OptInlineKeyboardMarkup  = func(markup *InlineKeyboardMarkup) sendOption {
		return func(r url.Values) {
			r.Set("reply_markup", structString(markup))
		}
	}
	OptReplyKeyboardMarkup = func(markup *ReplyKeyboardMarkup) sendOption {
		return func(r url.Values) {
			r.Set("reply_markup", structString(markup))
		}
	}
	OptReplyKeyboardRemoveSelective = func(r url.Values) {
		r.Set("reply_markup", structString(&replyKeyboardRemove{RemoveKeyboard: true, Selective: true}))
	}
	OptForceReply = func(r url.Values) {
		r.Set("reply_markup", structString(&forceReply{ForceReply: true}))
	}
	OptForceReplySelective = func(r url.Values) {
		r.Set("reply_markup", structString(&forceReply{ForceReply: true, Selective: true}))
	}
)

// SendMessage sends message to telegram chat. Available options
//   - OptParseModeHTML
//   - OptParseModeMarkdown
//   - OptDisableWebPagePreview
//   - OptDisableNotification
//   - OptReplyToMessageID(id int)
//   - OptSendingWithoutReply
//   - OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
//   - OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
//   - OptReplyKeyboardRemove
//   - OptReplyKeyboardRemoveSelective
//   - OptForceReply
//   - OptForceReplySelective
func (c *Client) SendMessage(chatID string, thread_id string, text string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	thread_id = strings.TrimSpace(thread_id)
	if thread_id != "" {
		req.Set("message_thread_id", thread_id)
	}
	req.Set("text", text)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest("/sendMessage", req, msg)
	return msg, err
}

// ForwardMessage forwards message from one chat to another. Available options:
//   - OptDisableNotification
func (c *Client) ForwardMessage(chatID, fromChatID string, messageID int, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("from_chat_id", fromChatID)
	req.Set("message_id", strconv.Itoa(messageID))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequestWithFiles("/forwardMessage", req, msg)
	return msg, err
}

type inputFile struct {
	field string
	name  string
}

// SendStickerFile send .webp file sticker. Available options:
//   - OptDisableNotification
//   - OptReplyToMessageID(id int)
//   - OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
//   - OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
//   - OptReplyKeyboardRemove
//   - OptReplyKeyboardRemoveSelective
//   - OptForceReply
//   - OptForceReplySelective
func (c *Client) SendStickerFile(chatID string, filename string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequestWithFiles("/sendSticker", req, msg, inputFile{field: "sticker", name: filename})
	return msg, err
}

// SendSticker send previously uploaded sticker. Available options:
//   - OptDisableNotification
//   - OptReplyToMessageID(id int)
//   - OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
//   - OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
//   - OptReplyKeyboardRemove
//   - OptReplyKeyboardRemoveSelective
//   - OptForceReply
//   - OptForceReplySelective
func (c *Client) SendSticker(chatID, fileID string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("sticker", fileID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest("/sendSticker", req, msg)
	return msg, err
}

type chatAction string

// Actions for SendChatAction
const (
	ActionTyping          chatAction = "typing"
	ActionUploadPhoto     chatAction = "upload_photo"
	ActionRecordVideo     chatAction = "record_video"
	ActionUploadVideo     chatAction = "upload_video"
	ActionRecordAudio     chatAction = "record_audio"
	ActionUploadAudio     chatAction = "upload_audio"
	ActionUploadDocument  chatAction = "upload_document"
	ActionFindLocation    chatAction = "find_location"
	ActionRecordVideoNote chatAction = "record_video_note"
	ActionUploadVideoNote chatAction = "upload_video_note"
)

func (c *Client) SendChatAction(chatID string, action chatAction) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("action", string(action))
	var sent bool
	return c.sendRequest("sendChatAction", req, &sent)
}
