package tbot

import (
	"context"
	"net/url"
)

// SendInvoice sends an invoice to a chat. Required fields are set directly;
// optional Bot API fields (provider_token, max_tip_amount,
// suggested_tip_amounts, provider_data, photo_url, photo_size, photo_width,
// photo_height, need_name, need_phone_number, need_email,
// need_shipping_address, send_phone_number_to_provider,
// send_email_to_provider, is_flexible, …) can be set through send options.
func (c *Client) SendInvoice(ctx context.Context, chatID ChatID, title, description, payload, currency string, prices []LabeledPrice, opts ...SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("title", title)
	req.Set("description", description)
	req.Set("payload", payload)
	req.Set("currency", currency)
	req.Set("prices", structString(prices))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.sendRequest(ctx, "/sendInvoice", req, msg)
	return msg, err
}

// CreateInvoiceLink creates a link for an invoice. Required fields are set
// directly; optional Bot API fields follow the same pattern as SendInvoice.
func (c *Client) CreateInvoiceLink(ctx context.Context, title, description, payload, currency string, prices []LabeledPrice, opts ...SendOption) (string, error) {
	req := url.Values{}
	req.Set("title", title)
	req.Set("description", description)
	req.Set("payload", payload)
	req.Set("currency", currency)
	req.Set("prices", structString(prices))
	for _, opt := range opts {
		opt(req)
	}
	var link string
	err := c.sendRequest(ctx, "/createInvoiceLink", req, &link)
	return link, err
}

// AnswerShippingQuery replies to a shipping query. When ok is true, pass the
// available shippingOptions; when ok is false, pass a human-readable
// errorMessage explaining why the order cannot be completed.
func (c *Client) AnswerShippingQuery(ctx context.Context, shippingQueryID string, ok bool, shippingOptions []ShippingOption, errorMessage string) (bool, error) {
	req := url.Values{}
	req.Set("shipping_query_id", shippingQueryID)
	if ok {
		req.Set("ok", "true")
		if len(shippingOptions) > 0 {
			req.Set("shipping_options", structString(shippingOptions))
		}
	} else {
		req.Set("ok", "false")
		if errorMessage != "" {
			req.Set("error_message", errorMessage)
		}
	}
	var result bool
	err := c.sendRequest(ctx, "/answerShippingQuery", req, &result)
	return result, err
}

// AnswerPreCheckoutQuery confirms or rejects a pre-checkout query. When ok is
// false, pass a human-readable errorMessage explaining the rejection reason.
func (c *Client) AnswerPreCheckoutQuery(ctx context.Context, preCheckoutQueryID string, ok bool, errorMessage string) (bool, error) {
	req := url.Values{}
	req.Set("pre_checkout_query_id", preCheckoutQueryID)
	if ok {
		req.Set("ok", "true")
	} else {
		req.Set("ok", "false")
		if errorMessage != "" {
			req.Set("error_message", errorMessage)
		}
	}
	var result bool
	err := c.sendRequest(ctx, "/answerPreCheckoutQuery", req, &result)
	return result, err
}

// RefundStarPayment refunds a Stars payment to a user.
func (c *Client) RefundStarPayment(ctx context.Context, userID int64, telegramPaymentChargeID string) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("telegram_payment_charge_id", telegramPaymentChargeID)
	var ok bool
	err := c.sendRequest(ctx, "/refundStarPayment", req, &ok)
	return ok, err
}

// EditUserStarSubscription allows the bot to cancel or re-enable a Stars
// subscription that the user has purchased. Set isCanceled to true to cancel
// or false to re-enable the subscription.
func (c *Client) EditUserStarSubscription(ctx context.Context, userID int64, telegramPaymentChargeID string, isCanceled bool) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("telegram_payment_charge_id", telegramPaymentChargeID)
	if isCanceled {
		req.Set("is_canceled", "true")
	} else {
		req.Set("is_canceled", "false")
	}
	var ok bool
	err := c.sendRequest(ctx, "/editUserStarSubscription", req, &ok)
	return ok, err
}

// GetStarTransactions returns the bot's Star transactions in reverse
// chronological order. Optional offset and limit can be set through send
// options.
func (c *Client) GetStarTransactions(ctx context.Context, opts ...SendOption) (*StarTransactions, error) {
	req := url.Values{}
	for _, opt := range opts {
		opt(req)
	}
	txns := &StarTransactions{}
	err := c.sendRequest(ctx, "/getStarTransactions", req, txns)
	return txns, err
}

// GetMyStarBalance returns the bot's current Star balance.
func (c *Client) GetMyStarBalance(ctx context.Context) (*StarAmount, error) {
	balance := &StarAmount{}
	err := c.sendRequest(ctx, "/getMyStarBalance", nil, balance)
	return balance, err
}

// GetAvailableGifts returns the list of gifts that the bot can send.
func (c *Client) GetAvailableGifts(ctx context.Context) (*Gifts, error) {
	gifts := &Gifts{}
	err := c.sendRequest(ctx, "/getAvailableGifts", nil, gifts)
	return gifts, err
}

// SendGift sends a gift to a user. Optional fields (text, entities,
// pay_for_upgrade) can be set through send options.
func (c *Client) SendGift(ctx context.Context, userID int64, giftID string, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("gift_id", giftID)
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/sendGift", req, &ok)
	return ok, err
}

// GiftPremiumSubscription gifts a Telegram Premium subscription to a user.
// Optional fields (text, entities) can be set through send options.
func (c *Client) GiftPremiumSubscription(ctx context.Context, userID int64, monthCount, starCount int, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("month_count", itoa(monthCount))
	req.Set("star_count", itoa(starCount))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/giftPremiumSubscription", req, &ok)
	return ok, err
}

// VerifyUserGift verifies ownership of a gift on behalf of a user.
func (c *Client) VerifyUserGift(ctx context.Context, userID int64, ownedGiftID string) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("owned_gift_id", ownedGiftID)
	var ok bool
	err := c.sendRequest(ctx, "/verifyUserGift", req, &ok)
	return ok, err
}

// ConvertGiftToStars converts an owned gift to Telegram Stars on behalf of
// the user.
func (c *Client) ConvertGiftToStars(ctx context.Context, userID int64, ownedGiftID string) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("owned_gift_id", ownedGiftID)
	var ok bool
	err := c.sendRequest(ctx, "/convertGiftToStars", req, &ok)
	return ok, err
}

// UpgradeGift upgrades a regular gift to a unique gift. Optional fields
// (keep_original_details, star_count) can be set through send options.
func (c *Client) UpgradeGift(ctx context.Context, userID int64, ownedGiftID string, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("owned_gift_id", ownedGiftID)
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/upgradeGift", req, &ok)
	return ok, err
}

// TransferGift transfers an owned gift to a new owner. Optional fields
// (star_count) can be set through send options.
func (c *Client) TransferGift(ctx context.Context, ownedGiftID string, newOwnerChatID ChatID, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("owned_gift_id", ownedGiftID)
	req.Set("new_owner_chat_id", newOwnerChatID.String())
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/transferGift", req, &ok)
	return ok, err
}

// GetReceivedGifts returns the list of gifts received by the bot or a
// business account. businessConnectionID identifies the business account;
// pass an empty string for the bot itself. Optional filters and pagination
// (exclude_unsaved, exclude_saved, exclude_unlimited, exclude_limited,
// exclude_unique, sort_by_price, offset, limit) can be set through send
// options.
func (c *Client) GetReceivedGifts(ctx context.Context, businessConnectionID string, opts ...SendOption) (*OwnedGifts, error) {
	req := url.Values{}
	if businessConnectionID != "" {
		req.Set("business_connection_id", businessConnectionID)
	}
	for _, opt := range opts {
		opt(req)
	}
	gifts := &OwnedGifts{}
	err := c.sendRequest(ctx, "/getReceivedGifts", req, gifts)
	return gifts, err
}

// SaveGift saves or unsaves an owned gift. Pass opts with the "saved" field
// to control the saved state.
func (c *Client) SaveGift(ctx context.Context, userID int64, ownedGiftID string, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("owned_gift_id", ownedGiftID)
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/saveGift", req, &ok)
	return ok, err
}
