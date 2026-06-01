package tbot

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
)

// editResult decodes the result of edit* / stopMessageLiveLocation calls. The
// Bot API returns the edited Message for chat messages, or the literal true for
// inline messages (in which case Message stays nil).
type editResult struct {
	Message *Message
}

func (e *editResult) UnmarshalJSON(b []byte) error {
	if bytes.Equal(bytes.TrimSpace(b), []byte("true")) {
		return nil
	}
	e.Message = &Message{}
	return json.Unmarshal(b, e.Message)
}

// edit runs an edit-style request, decoding either the edited Message or true.
func (c *Client) edit(ctx context.Context, method string, req url.Values) (*Message, error) {
	var res editResult
	if err := c.sendRequest(ctx, method, req, &res); err != nil {
		return nil, err
	}
	return res.Message, nil
}

// EditMessageText edits the text of a message. Target either a chat message
// (OptChatID + OptMessageID) or an inline message (OptInlineMessageID); editing
// an inline message returns (nil, nil).
func (c *Client) EditMessageText(ctx context.Context, text string, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("text", text)
	for _, opt := range opts {
		opt(req)
	}
	return c.edit(ctx, "/editMessageText", req)
}

// EditMessageCaption edits the caption of a media message. Set the new caption
// with OptCaption; target the message as in EditMessageText.
func (c *Client) EditMessageCaption(ctx context.Context, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	for _, opt := range opts {
		opt(req)
	}
	return c.edit(ctx, "/editMessageCaption", req)
}

// EditMessageMedia replaces the media of a message. The media's Media field
// must reference an existing file by file_id or URL.
func (c *Client) EditMessageMedia(ctx context.Context, media InputMedia, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("media", structString(media))
	for _, opt := range opts {
		opt(req)
	}
	return c.edit(ctx, "/editMessageMedia", req)
}

// EditMessageReplyMarkup edits only the reply markup of a message. Attach the
// new markup with OptInlineKeyboardMarkup.
func (c *Client) EditMessageReplyMarkup(ctx context.Context, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	for _, opt := range opts {
		opt(req)
	}
	return c.edit(ctx, "/editMessageReplyMarkup", req)
}

// EditMessageLiveLocation edits a live location message.
func (c *Client) EditMessageLiveLocation(ctx context.Context, latitude, longitude float64, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("latitude", formatFloat(latitude))
	req.Set("longitude", formatFloat(longitude))
	for _, opt := range opts {
		opt(req)
	}
	return c.edit(ctx, "/editMessageLiveLocation", req)
}

// StopMessageLiveLocation stops updating a live location message before
// live_period expires.
func (c *Client) StopMessageLiveLocation(ctx context.Context, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	for _, opt := range opts {
		opt(req)
	}
	return c.edit(ctx, "/stopMessageLiveLocation", req)
}

// StopPoll stops a poll the bot sent and returns the final Poll state.
func (c *Client) StopPoll(ctx context.Context, chatID ChatID, messageID int, opts ...SendOption) (*Poll, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("message_id", itoa(messageID))
	for _, opt := range opts {
		opt(req)
	}
	poll := &Poll{}
	err := c.sendRequest(ctx, "/stopPoll", req, poll)
	return poll, err
}
