package tbot

import (
	"context"
	"net/url"
)

// CreateForumTopic creates a topic in a forum supergroup chat.
func (c *Client) CreateForumTopic(ctx context.Context, chatID ChatID, name string, opts ...SendOption) (*ForumTopic, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("name", name)
	for _, opt := range opts {
		opt(req)
	}
	topic := &ForumTopic{}
	err := c.sendRequest(ctx, "/createForumTopic", req, topic)
	return topic, err
}

// EditForumTopic edits the name and icon of a forum topic.
func (c *Client) EditForumTopic(ctx context.Context, chatID ChatID, messageThreadID int, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("message_thread_id", itoa(messageThreadID))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/editForumTopic", req, &ok)
	return ok, err
}

// CloseForumTopic closes an open topic in a forum supergroup chat.
func (c *Client) CloseForumTopic(ctx context.Context, chatID ChatID, messageThreadID int) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("message_thread_id", itoa(messageThreadID))
	var ok bool
	err := c.sendRequest(ctx, "/closeForumTopic", req, &ok)
	return ok, err
}

// ReopenForumTopic reopens a closed topic in a forum supergroup chat.
func (c *Client) ReopenForumTopic(ctx context.Context, chatID ChatID, messageThreadID int) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("message_thread_id", itoa(messageThreadID))
	var ok bool
	err := c.sendRequest(ctx, "/reopenForumTopic", req, &ok)
	return ok, err
}

// DeleteForumTopic deletes a forum topic along with all its messages in a
// forum supergroup chat.
func (c *Client) DeleteForumTopic(ctx context.Context, chatID ChatID, messageThreadID int) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("message_thread_id", itoa(messageThreadID))
	var ok bool
	err := c.sendRequest(ctx, "/deleteForumTopic", req, &ok)
	return ok, err
}

// UnpinAllForumTopicMessages clears the list of pinned messages in a forum
// topic.
func (c *Client) UnpinAllForumTopicMessages(ctx context.Context, chatID ChatID, messageThreadID int) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("message_thread_id", itoa(messageThreadID))
	var ok bool
	err := c.sendRequest(ctx, "/unpinAllForumTopicMessages", req, &ok)
	return ok, err
}

// EditGeneralForumTopic edits the name of the General topic in a forum
// supergroup chat.
func (c *Client) EditGeneralForumTopic(ctx context.Context, chatID ChatID, name string) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("name", name)
	var ok bool
	err := c.sendRequest(ctx, "/editGeneralForumTopic", req, &ok)
	return ok, err
}

// CloseGeneralForumTopic closes the General topic in a forum supergroup chat.
func (c *Client) CloseGeneralForumTopic(ctx context.Context, chatID ChatID) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var ok bool
	err := c.sendRequest(ctx, "/closeGeneralForumTopic", req, &ok)
	return ok, err
}

// ReopenGeneralForumTopic reopens the General topic in a forum supergroup chat.
func (c *Client) ReopenGeneralForumTopic(ctx context.Context, chatID ChatID) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var ok bool
	err := c.sendRequest(ctx, "/reopenGeneralForumTopic", req, &ok)
	return ok, err
}

// HideGeneralForumTopic hides the General topic in a forum supergroup chat.
func (c *Client) HideGeneralForumTopic(ctx context.Context, chatID ChatID) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var ok bool
	err := c.sendRequest(ctx, "/hideGeneralForumTopic", req, &ok)
	return ok, err
}

// UnhideGeneralForumTopic unhides the General topic in a forum supergroup chat.
func (c *Client) UnhideGeneralForumTopic(ctx context.Context, chatID ChatID) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var ok bool
	err := c.sendRequest(ctx, "/unhideGeneralForumTopic", req, &ok)
	return ok, err
}

// UnpinAllGeneralForumTopicMessages clears the list of pinned messages in the
// General topic of a forum supergroup chat.
func (c *Client) UnpinAllGeneralForumTopicMessages(ctx context.Context, chatID ChatID) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var ok bool
	err := c.sendRequest(ctx, "/unpinAllGeneralForumTopicMessages", req, &ok)
	return ok, err
}
