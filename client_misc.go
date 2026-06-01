package tbot

import (
	"context"
	"net/url"
)

// GetUserChatBoosts returns the list of boosts the given user has applied to
// the chat.
func (c *Client) GetUserChatBoosts(ctx context.Context, chatID ChatID, userID int64) (*UserChatBoosts, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("user_id", itoa64(userID))
	boosts := &UserChatBoosts{}
	err := c.sendRequest(ctx, "/getUserChatBoosts", req, boosts)
	return boosts, err
}

// GetBusinessConnection returns information about the connection of the bot
// with a business account.
func (c *Client) GetBusinessConnection(ctx context.Context, businessConnectionID string) (*BusinessConnection, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	conn := &BusinessConnection{}
	err := c.sendRequest(ctx, "/getBusinessConnection", req, conn)
	return conn, err
}

// SetUserEmojiStatus sets the emoji status for the given user. Optional
// fields (emoji_status_custom_emoji_id, emoji_status_expiration_date) can be
// passed as send options.
func (c *Client) SetUserEmojiStatus(ctx context.Context, userID int64, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setUserEmojiStatus", req, &ok)
	return ok, err
}

// ReadBusinessMessage marks a message in a connected business account as read.
func (c *Client) ReadBusinessMessage(ctx context.Context, businessConnectionID string, chatID ChatID, messageID int) (bool, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	req.Set("chat_id", chatID.String())
	req.Set("message_id", itoa(messageID))
	var ok bool
	err := c.sendRequest(ctx, "/readBusinessMessage", req, &ok)
	return ok, err
}

// DeleteBusinessMessages deletes messages on behalf of a connected business
// account.
func (c *Client) DeleteBusinessMessages(ctx context.Context, businessConnectionID string, messageIDs []int) (bool, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	req.Set("message_ids", structString(messageIDs))
	var ok bool
	err := c.sendRequest(ctx, "/deleteBusinessMessages", req, &ok)
	return ok, err
}

// SetBusinessAccountName changes the first (and optionally last) name of a
// managed business account. Optional last_name can be set via send options.
func (c *Client) SetBusinessAccountName(ctx context.Context, businessConnectionID, firstName string, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	req.Set("first_name", firstName)
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setBusinessAccountName", req, &ok)
	return ok, err
}

// SetBusinessAccountUsername changes the username of a managed business
// account.
func (c *Client) SetBusinessAccountUsername(ctx context.Context, businessConnectionID, username string) (bool, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	req.Set("username", username)
	var ok bool
	err := c.sendRequest(ctx, "/setBusinessAccountUsername", req, &ok)
	return ok, err
}

// SetBusinessAccountBio changes the bio of a managed business account.
// Optional bio can be set via send options.
func (c *Client) SetBusinessAccountBio(ctx context.Context, businessConnectionID string, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setBusinessAccountBio", req, &ok)
	return ok, err
}

// SetBusinessAccountProfilePhoto sets a profile photo for a managed business
// account. The photo may be a file_id, URL, or a fresh upload. Optional
// is_personal flag can be set via send options.
func (c *Client) SetBusinessAccountProfilePhoto(ctx context.Context, businessConnectionID string, photo InputFile, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	var files []fileField
	addFile(req, &files, "photo", photo)
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.send(ctx, "/setBusinessAccountProfilePhoto", req, files, &ok)
	return ok, err
}

// RemoveBusinessAccountProfilePhoto removes the profile photo of a managed
// business account. Pass OptIsPersonal via opts to target the personal photo.
func (c *Client) RemoveBusinessAccountProfilePhoto(ctx context.Context, businessConnectionID string, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/removeBusinessAccountProfilePhoto", req, &ok)
	return ok, err
}

// SetBusinessAccountGiftSettings updates the gift settings of a managed
// business account.
func (c *Client) SetBusinessAccountGiftSettings(ctx context.Context, businessConnectionID string, showGiftButton bool, acceptedGiftTypes AcceptedGiftTypes) (bool, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	if showGiftButton {
		req.Set("show_gift_button", "true")
	} else {
		req.Set("show_gift_button", "false")
	}
	req.Set("accepted_gift_types", structString(acceptedGiftTypes))
	var ok bool
	err := c.sendRequest(ctx, "/setBusinessAccountGiftSettings", req, &ok)
	return ok, err
}

// GetBusinessAccountStarBalance returns the current Star balance of a managed
// business account.
func (c *Client) GetBusinessAccountStarBalance(ctx context.Context, businessConnectionID string) (*StarAmount, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	balance := &StarAmount{}
	err := c.sendRequest(ctx, "/getBusinessAccountStarBalance", req, balance)
	return balance, err
}

// GetBusinessAccountGifts returns the gifts owned by a managed business
// account. Optional filters and pagination can be set through send options
// (exclude_unsaved, exclude_saved, exclude_unlimited, exclude_limited,
// exclude_unique, sort_by_price, offset, limit).
func (c *Client) GetBusinessAccountGifts(ctx context.Context, businessConnectionID string, opts ...SendOption) (*OwnedGifts, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	for _, opt := range opts {
		opt(req)
	}
	gifts := &OwnedGifts{}
	err := c.sendRequest(ctx, "/getBusinessAccountGifts", req, gifts)
	return gifts, err
}

// TransferBusinessAccountStars transfers Telegram Stars from the bot to a
// managed business account.
func (c *Client) TransferBusinessAccountStars(ctx context.Context, businessConnectionID string, starCount int) (bool, error) {
	req := url.Values{}
	req.Set("business_connection_id", businessConnectionID)
	req.Set("star_count", itoa(starCount))
	var ok bool
	err := c.sendRequest(ctx, "/transferBusinessAccountStars", req, &ok)
	return ok, err
}

// SetGameScore sets the score of the specified user in a game. Target either
// a chat message (OptChatID + OptMessageID) or an inline message
// (OptInlineMessageID). Editing an inline message returns (nil, nil) on
// success.
func (c *Client) SetGameScore(ctx context.Context, userID int64, score int, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("score", itoa(score))
	for _, opt := range opts {
		opt(req)
	}
	return c.edit(ctx, "/setGameScore", req)
}

// GetGameHighScores returns the high scores table for a game. Target either a
// chat message (OptChatID + OptMessageID) or an inline message
// (OptInlineMessageID).
func (c *Client) GetGameHighScores(ctx context.Context, userID int64, opts ...SendOption) ([]GameHighScore, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	for _, opt := range opts {
		opt(req)
	}
	var scores []GameHighScore
	err := c.sendRequest(ctx, "/getGameHighScores", req, &scores)
	return scores, err
}
