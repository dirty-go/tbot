package tbot

import (
	"context"
	"net/url"
)

// BanChatMember bans a user in a group, supergroup, or channel. In supergroups
// and channels the user will not be able to return to the chat on their own
// until they are unbanned. Use OptUntilDate to set a ban expiry.
func (c *Client) BanChatMember(ctx context.Context, chatID ChatID, userID int64, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("user_id", itoa64(userID))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/banChatMember", req, &ok)
	return ok, err
}

// UnbanChatMember unbans a previously banned user in a supergroup or channel.
// Pass opts with OptOnlyIfBanned to only unban if the user was banned.
func (c *Client) UnbanChatMember(ctx context.Context, chatID ChatID, userID int64, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("user_id", itoa64(userID))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/unbanChatMember", req, &ok)
	return ok, err
}

// RestrictChatMember restricts a user in a supergroup. Pass the permissions
// the user should have after the restriction. Use OptUntilDate to set an expiry.
func (c *Client) RestrictChatMember(ctx context.Context, chatID ChatID, userID int64, perms *ChatPermissions, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("user_id", itoa64(userID))
	req.Set("permissions", structString(perms))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/restrictChatMember", req, &ok)
	return ok, err
}

// PromoteChatMember promotes or demotes a user in a supergroup or channel.
// Pass rights to grant; omitted fields default to false (demote).
func (c *Client) PromoteChatMember(ctx context.Context, chatID ChatID, userID int64, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("user_id", itoa64(userID))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/promoteChatMember", req, &ok)
	return ok, err
}

// SetChatAdministratorCustomTitle sets a custom title for an administrator
// promoted by the bot.
func (c *Client) SetChatAdministratorCustomTitle(ctx context.Context, chatID ChatID, userID int64, customTitle string) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("user_id", itoa64(userID))
	req.Set("custom_title", customTitle)
	var ok bool
	err := c.sendRequest(ctx, "/setChatAdministratorCustomTitle", req, &ok)
	return ok, err
}

// BanChatSenderChat bans a channel chat in a supergroup or channel.
func (c *Client) BanChatSenderChat(ctx context.Context, chatID ChatID, senderChatID int64) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("sender_chat_id", itoa64(senderChatID))
	var ok bool
	err := c.sendRequest(ctx, "/banChatSenderChat", req, &ok)
	return ok, err
}

// UnbanChatSenderChat unbans a previously banned channel chat in a supergroup
// or channel.
func (c *Client) UnbanChatSenderChat(ctx context.Context, chatID ChatID, senderChatID int64) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("sender_chat_id", itoa64(senderChatID))
	var ok bool
	err := c.sendRequest(ctx, "/unbanChatSenderChat", req, &ok)
	return ok, err
}

// SetChatPermissions sets the default chat permissions for all members. The
// bot must be an administrator with CanRestrictMembers.
func (c *Client) SetChatPermissions(ctx context.Context, chatID ChatID, perms *ChatPermissions, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("permissions", structString(perms))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setChatPermissions", req, &ok)
	return ok, err
}

// ExportChatInviteLink generates a new primary invite link for a chat; any
// previously generated primary link is revoked.
func (c *Client) ExportChatInviteLink(ctx context.Context, chatID ChatID) (string, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var link string
	err := c.sendRequest(ctx, "/exportChatInviteLink", req, &link)
	return link, err
}

// CreateChatInviteLink creates an additional invite link for a chat.
func (c *Client) CreateChatInviteLink(ctx context.Context, chatID ChatID, opts ...SendOption) (*ChatInviteLink, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	for _, opt := range opts {
		opt(req)
	}
	link := &ChatInviteLink{}
	err := c.sendRequest(ctx, "/createChatInviteLink", req, link)
	return link, err
}

// EditChatInviteLink edits a non-primary invite link created by the bot.
func (c *Client) EditChatInviteLink(ctx context.Context, chatID ChatID, inviteLink string, opts ...SendOption) (*ChatInviteLink, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("invite_link", inviteLink)
	for _, opt := range opts {
		opt(req)
	}
	link := &ChatInviteLink{}
	err := c.sendRequest(ctx, "/editChatInviteLink", req, link)
	return link, err
}

// RevokeChatInviteLink revokes an invite link created by the bot.
func (c *Client) RevokeChatInviteLink(ctx context.Context, chatID ChatID, inviteLink string) (*ChatInviteLink, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("invite_link", inviteLink)
	link := &ChatInviteLink{}
	err := c.sendRequest(ctx, "/revokeChatInviteLink", req, link)
	return link, err
}

// CreateChatSubscriptionInviteLink creates a subscription invite link for a
// channel chat. The bot must have CanInviteUsers rights. The link allows
// subscribing to the channel and paying subscriptionPrice Stars every
// subscriptionPeriod seconds.
func (c *Client) CreateChatSubscriptionInviteLink(ctx context.Context, chatID ChatID, subscriptionPeriod, subscriptionPrice int, opts ...SendOption) (*ChatInviteLink, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("subscription_period", itoa(subscriptionPeriod))
	req.Set("subscription_price", itoa(subscriptionPrice))
	for _, opt := range opts {
		opt(req)
	}
	link := &ChatInviteLink{}
	err := c.sendRequest(ctx, "/createChatSubscriptionInviteLink", req, link)
	return link, err
}

// EditChatSubscriptionInviteLink edits a subscription invite link created by
// the bot.
func (c *Client) EditChatSubscriptionInviteLink(ctx context.Context, chatID ChatID, inviteLink string, opts ...SendOption) (*ChatInviteLink, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("invite_link", inviteLink)
	for _, opt := range opts {
		opt(req)
	}
	link := &ChatInviteLink{}
	err := c.sendRequest(ctx, "/editChatSubscriptionInviteLink", req, link)
	return link, err
}

// ApproveChatJoinRequest approves a chat join request from a user.
func (c *Client) ApproveChatJoinRequest(ctx context.Context, chatID ChatID, userID int64) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("user_id", itoa64(userID))
	var ok bool
	err := c.sendRequest(ctx, "/approveChatJoinRequest", req, &ok)
	return ok, err
}

// DeclineChatJoinRequest declines a chat join request from a user.
func (c *Client) DeclineChatJoinRequest(ctx context.Context, chatID ChatID, userID int64) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("user_id", itoa64(userID))
	var ok bool
	err := c.sendRequest(ctx, "/declineChatJoinRequest", req, &ok)
	return ok, err
}

// SetChatPhoto sets a new profile photo for the chat. The photo must be an
// upload (FilePath / FileReader / FileBytes).
func (c *Client) SetChatPhoto(ctx context.Context, chatID ChatID, photo InputFile) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var files []fileField
	addFile(req, &files, "photo", photo)
	var ok bool
	err := c.send(ctx, "/setChatPhoto", req, files, &ok)
	return ok, err
}

// DeleteChatPhoto deletes a chat photo.
func (c *Client) DeleteChatPhoto(ctx context.Context, chatID ChatID) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var ok bool
	err := c.sendRequest(ctx, "/deleteChatPhoto", req, &ok)
	return ok, err
}

// SetChatTitle changes the title of a chat.
func (c *Client) SetChatTitle(ctx context.Context, chatID ChatID, title string) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("title", title)
	var ok bool
	err := c.sendRequest(ctx, "/setChatTitle", req, &ok)
	return ok, err
}

// SetChatDescription changes the description of a group, supergroup, or
// channel. Pass an empty string to remove the description.
func (c *Client) SetChatDescription(ctx context.Context, chatID ChatID, description string) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	if description != "" {
		req.Set("description", description)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setChatDescription", req, &ok)
	return ok, err
}

// PinChatMessage pins a message in a group, supergroup, or channel.
// Pass OptDisableNotification to pin silently.
func (c *Client) PinChatMessage(ctx context.Context, chatID ChatID, messageID int, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("message_id", itoa(messageID))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/pinChatMessage", req, &ok)
	return ok, err
}

// UnpinChatMessage unpins a message in a group, supergroup, or channel.
// If messageID is 0 the most recent pinned message is unpinned.
func (c *Client) UnpinChatMessage(ctx context.Context, chatID ChatID, messageID int) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	if messageID != 0 {
		req.Set("message_id", itoa(messageID))
	}
	var ok bool
	err := c.sendRequest(ctx, "/unpinChatMessage", req, &ok)
	return ok, err
}

// UnpinAllChatMessages clears the list of pinned messages in a chat.
func (c *Client) UnpinAllChatMessages(ctx context.Context, chatID ChatID) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var ok bool
	err := c.sendRequest(ctx, "/unpinAllChatMessages", req, &ok)
	return ok, err
}

// LeaveChat makes the bot leave a group, supergroup, or channel.
func (c *Client) LeaveChat(ctx context.Context, chatID ChatID) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var ok bool
	err := c.sendRequest(ctx, "/leaveChat", req, &ok)
	return ok, err
}

// GetChat returns up-to-date information about the chat as ChatFullInfo.
func (c *Client) GetChat(ctx context.Context, chatID ChatID) (*ChatFullInfo, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	chat := &ChatFullInfo{}
	err := c.sendRequest(ctx, "/getChat", req, chat)
	return chat, err
}

// GetChatAdministrators returns a list of administrators in a chat.
func (c *Client) GetChatAdministrators(ctx context.Context, chatID ChatID) ([]ChatMember, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var members []ChatMember
	err := c.sendRequest(ctx, "/getChatAdministrators", req, &members)
	return members, err
}

// GetChatMemberCount returns the number of members in a chat.
func (c *Client) GetChatMemberCount(ctx context.Context, chatID ChatID) (int, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var count int
	err := c.sendRequest(ctx, "/getChatMemberCount", req, &count)
	return count, err
}

// GetChatMember returns information about a specific member of a chat.
func (c *Client) GetChatMember(ctx context.Context, chatID ChatID, userID int64) (*ChatMember, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("user_id", itoa64(userID))
	member := &ChatMember{}
	err := c.sendRequest(ctx, "/getChatMember", req, member)
	return member, err
}

// SetChatStickerSet sets a new group sticker set for a supergroup.
func (c *Client) SetChatStickerSet(ctx context.Context, chatID ChatID, stickerSetName string) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("sticker_set_name", stickerSetName)
	var ok bool
	err := c.sendRequest(ctx, "/setChatStickerSet", req, &ok)
	return ok, err
}

// DeleteChatStickerSet removes the sticker set linked to the supergroup.
func (c *Client) DeleteChatStickerSet(ctx context.Context, chatID ChatID) (bool, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var ok bool
	err := c.sendRequest(ctx, "/deleteChatStickerSet", req, &ok)
	return ok, err
}

// GetForumTopicIconStickers returns custom emoji stickers that can be used as
// forum topic icons by any user.
func (c *Client) GetForumTopicIconStickers(ctx context.Context) ([]Sticker, error) {
	var stickers []Sticker
	err := c.sendRequest(ctx, "/getForumTopicIconStickers", nil, &stickers)
	return stickers, err
}

// GetUserProfilePhotos returns a portion of the specified user's profile
// photos. offset is the sequential number of the first photo to return (0
// for all); limit caps the result (1–100, default 100).
func (c *Client) GetUserProfilePhotos(ctx context.Context, userID int64, offset, limit int) (*UserProfilePhotos, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	if offset > 0 {
		req.Set("offset", itoa(offset))
	}
	if limit > 0 {
		req.Set("limit", itoa(limit))
	}
	photos := &UserProfilePhotos{}
	err := c.sendRequest(ctx, "/getUserProfilePhotos", req, photos)
	return photos, err
}
