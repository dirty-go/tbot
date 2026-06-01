package tbot

import (
	"context"
	"net/url"
)

// AnswerInlineQuery sends answers to an inline query. results may contain up
// to 50 items. Use opts to set cache_time, is_personal, next_offset, and the
// button shown above the results.
func (c *Client) AnswerInlineQuery(ctx context.Context, inlineQueryID string, results []InlineQueryResult, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("inline_query_id", inlineQueryID)
	req.Set("results", structString(results))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/answerInlineQuery", req, &ok)
	return ok, err
}

// AnswerWebAppQuery sets the result of an interaction with a Web App and sends
// a corresponding message on behalf of the user to the chat from which the
// query originated.
func (c *Client) AnswerWebAppQuery(ctx context.Context, webAppQueryID string, result InlineQueryResult) (*SentWebAppMessage, error) {
	req := url.Values{}
	req.Set("web_app_query_id", webAppQueryID)
	req.Set("result", structString(result))
	msg := &SentWebAppMessage{}
	err := c.sendRequest(ctx, "/answerWebAppQuery", req, msg)
	return msg, err
}

// SavePreparedInlineMessage stores a message that can be sent by a user of a
// Mini App. allowUserChats, allowBotChats, allowGroupChats, and
// allowChannelChats control which chat types the prepared message may be sent
// to. Returns a PreparedInlineMessage valid for 10 minutes.
func (c *Client) SavePreparedInlineMessage(ctx context.Context, userID int64, result InlineQueryResult, opts ...SendOption) (*PreparedInlineMessage, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("result", structString(result))
	for _, opt := range opts {
		opt(req)
	}
	msg := &PreparedInlineMessage{}
	err := c.sendRequest(ctx, "/savePreparedInlineMessage", req, msg)
	return msg, err
}
