package tbot

import (
	"context"
	"net/url"
)

// StickerFormat is the wire type for the sticker format field.
type StickerFormat string

// Valid values for StickerFormat.
const (
	StickerFormatStatic   StickerFormat = "static"
	StickerFormatAnimated StickerFormat = "animated"
	StickerFormatVideo    StickerFormat = "video"
)

// UploadStickerFile uploads a .WEBP, .PNG, .TGS, or .WEBM file with a sticker
// for later use in createNewStickerSet and addStickerToSet methods.
func (c *Client) UploadStickerFile(ctx context.Context, userID int64, sticker InputFile, stickerFormat StickerFormat) (*File, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("sticker_format", string(stickerFormat))
	var files []fileField
	addFile(req, &files, "sticker", sticker)
	var f File
	err := c.send(ctx, "/uploadStickerFile", req, files, &f)
	return &f, err
}

// GetStickerSet returns a sticker set by its name.
func (c *Client) GetStickerSet(ctx context.Context, name string) (*StickerSet, error) {
	req := url.Values{}
	req.Set("name", name)
	var set StickerSet
	err := c.sendRequest(ctx, "/getStickerSet", req, &set)
	return &set, err
}

// GetCustomEmojiStickers returns information about custom emoji stickers by
// their identifiers.
func (c *Client) GetCustomEmojiStickers(ctx context.Context, customEmojiIDs []string) ([]Sticker, error) {
	req := url.Values{}
	req.Set("custom_emoji_ids", structString(customEmojiIDs))
	var stickers []Sticker
	err := c.sendRequest(ctx, "/getCustomEmojiStickers", req, &stickers)
	return stickers, err
}

// CreateNewStickerSet creates a new sticker set owned by a user. The bot will
// be able to edit the sticker set. Use opts to set optional fields such as
// sticker_type and needs_repainting.
func (c *Client) CreateNewStickerSet(ctx context.Context, userID int64, name, title string, stickers []InputSticker, opts ...SendOption) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("name", name)
	req.Set("title", title)
	req.Set("stickers", structString(stickers))
	for _, opt := range opts {
		opt(req)
	}
	var ok bool
	err := c.sendRequest(ctx, "/createNewStickerSet", req, &ok)
	return ok, err
}

// AddStickerToSet adds a new sticker to a set created by the bot.
func (c *Client) AddStickerToSet(ctx context.Context, userID int64, name string, sticker InputSticker) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("name", name)
	req.Set("sticker", structString(sticker))
	var ok bool
	err := c.sendRequest(ctx, "/addStickerToSet", req, &ok)
	return ok, err
}

// SetStickerPositionInSet moves a sticker in a set created by the bot to a
// specific position in the set.
func (c *Client) SetStickerPositionInSet(ctx context.Context, sticker string, position int) (bool, error) {
	req := url.Values{}
	req.Set("sticker", sticker)
	req.Set("position", itoa(position))
	var ok bool
	err := c.sendRequest(ctx, "/setStickerPositionInSet", req, &ok)
	return ok, err
}

// DeleteStickerFromSet deletes a sticker from a set created by the bot.
func (c *Client) DeleteStickerFromSet(ctx context.Context, sticker string) (bool, error) {
	req := url.Values{}
	req.Set("sticker", sticker)
	var ok bool
	err := c.sendRequest(ctx, "/deleteStickerFromSet", req, &ok)
	return ok, err
}

// ReplaceStickerInSet replaces an existing sticker in a sticker set with a new
// one. The method is equivalent to calling deleteStickerFromSet, then
// addStickerToSet, but will not fail if the original sticker was not found.
func (c *Client) ReplaceStickerInSet(ctx context.Context, userID int64, name, oldSticker string, sticker InputSticker) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("name", name)
	req.Set("old_sticker", oldSticker)
	req.Set("sticker", structString(sticker))
	var ok bool
	err := c.sendRequest(ctx, "/replaceStickerInSet", req, &ok)
	return ok, err
}

// SetStickerEmojiList changes the list of emoji assigned to a regular or
// custom emoji sticker. The sticker must belong to a sticker set created by
// the bot.
func (c *Client) SetStickerEmojiList(ctx context.Context, sticker string, emojiList []string) (bool, error) {
	req := url.Values{}
	req.Set("sticker", sticker)
	req.Set("emoji_list", structString(emojiList))
	var ok bool
	err := c.sendRequest(ctx, "/setStickerEmojiList", req, &ok)
	return ok, err
}

// SetStickerKeywords changes the search keywords assigned to a regular or
// custom emoji sticker. Pass nil or an empty slice to clear the keywords.
func (c *Client) SetStickerKeywords(ctx context.Context, sticker string, keywords []string) (bool, error) {
	req := url.Values{}
	req.Set("sticker", sticker)
	if len(keywords) > 0 {
		req.Set("keywords", structString(keywords))
	}
	var ok bool
	err := c.sendRequest(ctx, "/setStickerKeywords", req, &ok)
	return ok, err
}

// SetStickerMaskPosition changes the mask position of a mask sticker. Pass
// nil to clear the mask position.
func (c *Client) SetStickerMaskPosition(ctx context.Context, sticker string, maskPosition *MaskPosition) (bool, error) {
	req := url.Values{}
	req.Set("sticker", sticker)
	if maskPosition != nil {
		req.Set("mask_position", structString(maskPosition))
	}
	var ok bool
	err := c.sendRequest(ctx, "/setStickerMaskPosition", req, &ok)
	return ok, err
}

// SetStickerSetTitle sets the title of a created sticker set.
func (c *Client) SetStickerSetTitle(ctx context.Context, name, title string) (bool, error) {
	req := url.Values{}
	req.Set("name", name)
	req.Set("title", title)
	var ok bool
	err := c.sendRequest(ctx, "/setStickerSetTitle", req, &ok)
	return ok, err
}

// SetStickerSetThumbnail sets the thumbnail of a regular or mask sticker set.
// Pass nil for thumbnail to remove the thumbnail; in that case the first
// sticker will be used as the thumbnail.
func (c *Client) SetStickerSetThumbnail(ctx context.Context, userID int64, name string, format StickerFormat, thumbnail InputFile) (bool, error) {
	req := url.Values{}
	req.Set("user_id", itoa64(userID))
	req.Set("name", name)
	req.Set("format", string(format))
	var files []fileField
	if thumbnail != nil {
		addFile(req, &files, "thumbnail", thumbnail)
	}
	var ok bool
	err := c.send(ctx, "/setStickerSetThumbnail", req, files, &ok)
	return ok, err
}

// SetCustomEmojiStickerSetThumbnail sets the thumbnail of a custom emoji
// sticker set. Pass an empty string for customEmojiID to remove the thumbnail
// and use the first sticker in the set.
func (c *Client) SetCustomEmojiStickerSetThumbnail(ctx context.Context, name, customEmojiID string) (bool, error) {
	req := url.Values{}
	req.Set("name", name)
	if customEmojiID != "" {
		req.Set("custom_emoji_id", customEmojiID)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setCustomEmojiStickerSetThumbnail", req, &ok)
	return ok, err
}

// DeleteStickerSet deletes a sticker set that was created by the bot.
func (c *Client) DeleteStickerSet(ctx context.Context, name string) (bool, error) {
	req := url.Values{}
	req.Set("name", name)
	var ok bool
	err := c.sendRequest(ctx, "/deleteStickerSet", req, &ok)
	return ok, err
}
