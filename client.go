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
