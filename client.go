package tbot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultRPS         = 30.0
	defaultBurst       = 30
	defaultMaxRetries  = 3
	defaultHTTPTimeout = 120 * time.Second
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

	httpClient  *http.Client
	rateLimiter *rateLimiter
	maxRetries  int
	sleepFn     func(time.Duration)
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

func NewClient(token string, baseURL string, opts ...ClientOptions) *Client {
	if baseURL == "" {
		baseURL = apiBaseURL
	}
	c := &Client{
		token:       token,
		baseURL:     baseURL,
		url:         fmt.Sprintf("%s/bot%s", baseURL, token) + "%s",
		logger:      nopLogger{},
		maxRetries:  defaultMaxRetries,
		rateLimiter: newRateLimiter(defaultRPS, defaultBurst),
		httpClient: &http.Client{
			Timeout:   defaultHTTPTimeout,
			Transport: netTransport,
		},
		sleepFn: time.Sleep,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.logger == nil {
		c.logger = nopLogger{}
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout:   defaultHTTPTimeout,
			Transport: netTransport,
		}
	}
	if c.sleepFn == nil {
		c.sleepFn = time.Sleep
	}
	return c
}

func structString(s any) string {
	str, _ := json.Marshal(s)
	return string(str)
}

// Me returns info about bot as a User object
func (c *Client) Me() (*User, error) {
	var me User
	err := c.sendRequest("/getMe", nil, &me)
	return &me, err
}

// GetUpdates receives incoming updates using long polling. It uses the client's
// configured update fields (`nextOffset`, `bufferSize`, `timeout`,
// `updateParams`) and advances `nextOffset` after a successful fetch.
func (c *Client) GetUpdates() ([]Update, error) {
	req := url.Values{}
	for key, values := range c.updateParams {
		for _, value := range values {
			req.Add(key, value)
		}
	}
	if c.nextOffset > 0 {
		req.Set("offset", strconv.Itoa(c.nextOffset))
	}
	if c.bufferSize > 0 {
		req.Set("limit", strconv.Itoa(c.bufferSize))
	}
	if c.timeout > 0 {
		req.Set("timeout", strconv.Itoa(c.timeout))
	}

	var updates []Update
	if err := c.sendRequest("/getUpdates", req, &updates); err != nil {
		return nil, err
	}
	for _, update := range updates {
		next := update.UpdateID + 1
		if next > c.nextOffset {
			c.nextOffset = next
		}
	}
	return updates, nil
}

// SetWebhook specifies a URL and receives incoming updates via webhook.
// Additional Bot API fields can be set through send options.
func (c *Client) SetWebhook(webhookURL string, opts ...sendOption) (bool, error) {
	req := url.Values{}
	req.Set("url", webhookURL)
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest("/setWebhook", req, &ok)
	return ok, err
}

// DeleteWebhook removes webhook integration and can optionally drop pending updates.
func (c *Client) DeleteWebhook(dropPendingUpdates bool) (bool, error) {
	req := url.Values{}
	if dropPendingUpdates {
		req.Set("drop_pending_updates", "true")
	}
	var ok bool
	err := c.sendRequest("/deleteWebhook", req, &ok)
	return ok, err
}

// GetWebhookInfo retrieves current webhook status.
func (c *Client) GetWebhookInfo() (*WebhookInfo, error) {
	var info WebhookInfo
	err := c.sendRequest("/getWebhookInfo", nil, &info)
	return &info, err
}

// AnswerCallbackQuery sends an answer to a callback query.
func (c *Client) AnswerCallbackQuery(callbackQueryID, text string, showAlert bool, opts ...sendOption) (bool, error) {
	req := url.Values{}
	req.Set("callback_query_id", callbackQueryID)
	if text != "" {
		req.Set("text", text)
	}
	if showAlert {
		req.Set("show_alert", "true")
	}
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest("/answerCallbackQuery", req, &ok)
	return ok, err
}

var (
	OptDisableWebPagePreview = func(r url.Values) { r.Set("disable_web_page_preview", "true") }
	OptReplyKeyboardRemove   = func(r url.Values) {
		r.Set("reply_markup", structString(&ReplyKeyboardRemove{RemoveKeyboard: true}))
	}
	OptInlineKeyboardMarkup = func(markup *InlineKeyboardMarkup) sendOption {
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
		r.Set("reply_markup", structString(&ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}))
	}
	OptForceReply = func(r url.Values) {
		r.Set("reply_markup", structString(&ForceReply{ForceReply: true}))
	}
	OptForceReplySelective = func(r url.Values) {
		r.Set("reply_markup", structString(&ForceReply{ForceReply: true, Selective: true}))
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

// SendLocation sends a geographic point.
func (c *Client) SendLocation(chatID string, latitude, longitude float64, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	req.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest("/sendLocation", req, msg)
	return msg, err
}

// SendVenue sends venue information.
func (c *Client) SendVenue(chatID string, latitude, longitude float64, title, address string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	req.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	req.Set("title", title)
	req.Set("address", address)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest("/sendVenue", req, msg)
	return msg, err
}

// SendContact sends a contact card.
func (c *Client) SendContact(chatID, phoneNumber, firstName string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("phone_number", phoneNumber)
	req.Set("first_name", firstName)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest("/sendContact", req, msg)
	return msg, err
}

// SendPoll sends a native Telegram poll.
func (c *Client) SendPoll(chatID, question string, options []InputPollOption, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("question", question)
	req.Set("options", structString(options))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest("/sendPoll", req, msg)
	return msg, err
}

// SendDice sends an animated emoji that displays a random value.
func (c *Client) SendDice(chatID string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest("/sendDice", req, msg)
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
