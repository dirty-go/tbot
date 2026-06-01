package tbot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultRPS         = 30.0
	defaultBurst       = 30
	defaultMaxRetries  = 3
	defaultHTTPTimeout = 120 * time.Second
)

// Client is a Telegram Bot API client. Its methods are safe for concurrent use
// by multiple goroutines, with one exception: GetUpdates mutates the client's
// long-poll offset and must be driven from a single goroutine.
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

// NewClient builds a Client for the given bot token. baseURL overrides the Bot
// API host (pass "" for the public api.telegram.org). Behaviour is tuned with
// ClientOptions.
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
	return c
}

// structString marshals v to a compact JSON string for a form field. The
// supplied request types never fail to marshal, so the error is discarded.
func structString(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// send dispatches a request as multipart when files are present and as a
// form-encoded body otherwise.
func (c *Client) send(ctx context.Context, method string, req url.Values, files []fileField, out any) error {
	if len(files) > 0 {
		return c.sendMultipart(ctx, method, req, files, out)
	}
	return c.sendRequest(ctx, method, req, out)
}

// addFile routes an InputFile into the request: a reference (file_id/URL) is
// written as the form field, while an upload is queued as a multipart part.
func addFile(req url.Values, files *[]fileField, field string, f InputFile) {
	if f == nil {
		return
	}
	if f.needsUpload() {
		*files = append(*files, fileField{field: field, file: f})
		return
	}
	if v := f.value(); v != "" {
		req.Set(field, v)
	}
}

// Me returns basic information about the bot as a User object.
func (c *Client) Me(ctx context.Context) (*User, error) {
	var me User
	err := c.sendRequest(ctx, "/getMe", nil, &me)
	return &me, err
}

// GetUpdates receives incoming updates using long polling. It uses the client's
// configured update fields (nextOffset, bufferSize, timeout, updateParams) and
// advances nextOffset after a successful fetch. It must be called from a single
// goroutine; concurrent calls race on the offset.
func (c *Client) GetUpdates(ctx context.Context) ([]Update, error) {
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
	if err := c.sendRequest(ctx, "/getUpdates", req, &updates); err != nil {
		return nil, err
	}
	for _, update := range updates {
		if next := update.UpdateID + 1; next > c.nextOffset {
			c.nextOffset = next
		}
	}
	return updates, nil
}

// SetWebhook specifies a URL to receive incoming updates via webhook.
// Additional Bot API fields can be set through send options.
func (c *Client) SetWebhook(ctx context.Context, webhookURL string, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("url", webhookURL)
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setWebhook", req, &ok)
	return ok, err
}

// DeleteWebhook removes webhook integration and can optionally drop pending updates.
func (c *Client) DeleteWebhook(ctx context.Context, dropPendingUpdates bool) (bool, error) {
	req := url.Values{}
	if dropPendingUpdates {
		req.Set("drop_pending_updates", "true")
	}
	var ok bool
	err := c.sendRequest(ctx, "/deleteWebhook", req, &ok)
	return ok, err
}

// GetWebhookInfo retrieves the current webhook status.
func (c *Client) GetWebhookInfo(ctx context.Context) (*WebhookInfo, error) {
	var info WebhookInfo
	err := c.sendRequest(ctx, "/getWebhookInfo", nil, &info)
	return &info, err
}

// AnswerCallbackQuery sends an answer to a callback query from an inline keyboard.
func (c *Client) AnswerCallbackQuery(ctx context.Context, callbackQueryID, text string, showAlert bool, opts ...SendOption) (bool, error) {
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
	err := c.sendRequest(ctx, "/answerCallbackQuery", req, &ok)
	return ok, err
}

// SendMessage sends a text message to a chat. Target a forum topic with
// OptMessageThreadID; format text with OptParseModeHTML / OptParseModeMarkdown;
// attach keyboards with the reply-markup options.
func (c *Client) SendMessage(ctx context.Context, chatID ChatID, text string, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("text", text)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest(ctx, "/sendMessage", req, msg)
	return msg, err
}

// ForwardMessage forwards a single message from one chat to another.
func (c *Client) ForwardMessage(ctx context.Context, chatID, fromChatID ChatID, messageID int, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("from_chat_id", fromChatID.String())
	req.Set("message_id", strconv.Itoa(messageID))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest(ctx, "/forwardMessage", req, msg)
	return msg, err
}

// SendSticker sends a sticker. The sticker may be an existing file_id or URL
// (FileID / FileURL) or a fresh .webp/.tgs/.webm upload (FilePath / FileReader
// / FileBytes).
func (c *Client) SendSticker(ctx context.Context, chatID ChatID, sticker InputFile, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var files []fileField
	addFile(req, &files, "sticker", sticker)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.send(ctx, "/sendSticker", req, files, msg)
	return msg, err
}

// SendLocation sends a geographic point.
func (c *Client) SendLocation(ctx context.Context, chatID ChatID, latitude, longitude float64, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	req.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest(ctx, "/sendLocation", req, msg)
	return msg, err
}

// SendVenue sends information about a venue.
func (c *Client) SendVenue(ctx context.Context, chatID ChatID, latitude, longitude float64, title, address string, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	req.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	req.Set("title", title)
	req.Set("address", address)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest(ctx, "/sendVenue", req, msg)
	return msg, err
}

// SendContact sends a phone contact.
func (c *Client) SendContact(ctx context.Context, chatID ChatID, phoneNumber, firstName string, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("phone_number", phoneNumber)
	req.Set("first_name", firstName)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest(ctx, "/sendContact", req, msg)
	return msg, err
}

// SendPoll sends a native Telegram poll.
func (c *Client) SendPoll(ctx context.Context, chatID ChatID, question string, options []InputPollOption, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("question", question)
	req.Set("options", structString(options))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest(ctx, "/sendPoll", req, msg)
	return msg, err
}

// SendDice sends an animated emoji that displays a random value. Choose the
// emoji with OptCaption-style raw options or rely on the default die.
func (c *Client) SendDice(ctx context.Context, chatID ChatID, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest(ctx, "/sendDice", req, msg)
	return msg, err
}

type chatAction string

// Actions for SendChatAction.
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

// SendChatAction tells the user that something is happening on the bot's side.
func (c *Client) SendChatAction(ctx context.Context, chatID ChatID, action chatAction, opts ...SendOption) error {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("action", string(action))
	for _, opt := range opts {
		opt(req)
	}
	var sent bool
	return c.sendRequest(ctx, "/sendChatAction", req, &sent)
}
