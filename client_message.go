package tbot

import (
	"context"
	"net/url"
)

// CopyMessage copies a single message of any kind and returns the new message's
// id. Unlike ForwardMessage the copy is not linked to the original.
func (c *Client) CopyMessage(ctx context.Context, chatID, fromChatID ChatID, messageID int, opts ...SendOption) (*MessageId, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("from_chat_id", fromChatID.String())
	req.Set("message_id", itoa(messageID))
	for _, opt := range opts {
		opt(req)
	}
	id := &MessageId{}
	err := c.sendRequest(ctx, "/copyMessage", req, id)
	return id, err
}

// CopyMessages copies multiple messages and returns the new message ids. The
// copies are not linked to the originals.
func (c *Client) CopyMessages(ctx context.Context, chatID, fromChatID ChatID, messageIDs []int, opts ...SendOption) ([]MessageId, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("from_chat_id", fromChatID.String())
	req.Set("message_ids", structString(messageIDs))
	for _, opt := range opts {
		opt(req)
	}
	var ids []MessageId
	err := c.sendRequest(ctx, "/copyMessages", req, &ids)
	return ids, err
}

// ForwardMessages forwards multiple messages of any kind and returns the new
// message ids. Messages that can't be forwarded are skipped.
func (c *Client) ForwardMessages(ctx context.Context, chatID, fromChatID ChatID, messageIDs []int, opts ...SendOption) ([]MessageId, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("from_chat_id", fromChatID.String())
	req.Set("message_ids", structString(messageIDs))
	for _, opt := range opts {
		opt(req)
	}
	var ids []MessageId
	err := c.sendRequest(ctx, "/forwardMessages", req, &ids)
	return ids, err
}

// DeleteMessage deletes a message, including service messages, subject to the
// Bot API's deletion-time limits.
func (c *Client) DeleteMessage(ctx context.Context, chatID ChatID, messageID int) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("message_id", itoa(messageID))
	var ok bool
	err := c.sendRequest(ctx, "/deleteMessage", req, &ok)
	return ok, err
}

// DeleteMessages deletes multiple messages simultaneously.
func (c *Client) DeleteMessages(ctx context.Context, chatID ChatID, messageIDs []int) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("message_ids", structString(messageIDs))
	var ok bool
	err := c.sendRequest(ctx, "/deleteMessages", req, &ok)
	return ok, err
}

// SetMessageReaction sets the bot's reactions on a message. Pass a nil or empty
// reaction slice to clear it; make the animation big with OptIsBig.
func (c *Client) SetMessageReaction(ctx context.Context, chatID ChatID, messageID int, reaction []ReactionType, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("message_id", itoa(messageID))
	if len(reaction) > 0 {
		req.Set("reaction", structString(reaction))
	}
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setMessageReaction", req, &ok)
	return ok, err
}
