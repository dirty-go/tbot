package tbot

// Sticker represents a sticker.
type Sticker struct {
	FileID           string        `json:"file_id"`
	FileUniqueID     string        `json:"file_unique_id"`
	Type             string        `json:"type"`
	Width            int           `json:"width"`
	Height           int           `json:"height"`
	IsAnimated       bool          `json:"is_animated"`
	IsVideo          bool          `json:"is_video"`
	Thumbnail        *PhotoSize    `json:"thumbnail,omitempty"`
	Emoji            string        `json:"emoji,omitempty"`
	SetName          string        `json:"set_name,omitempty"`
	PremiumAnimation *File         `json:"premium_animation,omitempty"`
	MaskPosition     *MaskPosition `json:"mask_position,omitempty"`
	CustomEmojiID    string        `json:"custom_emoji_id,omitempty"`
	NeedsRepainting  bool          `json:"needs_repainting,omitempty"`
	FileSize         int64         `json:"file_size,omitempty"`
}

// Sticker Type values.
const (
	StickerTypeRegular     = "regular"
	StickerTypeMask        = "mask"
	StickerTypeCustomEmoji = "custom_emoji"
)

// StickerSet represents a sticker set.
type StickerSet struct {
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	StickerType string     `json:"sticker_type"`
	Stickers    []Sticker  `json:"stickers"`
	Thumbnail   *PhotoSize `json:"thumbnail,omitempty"`
}

// MaskPosition describes the position on faces where a mask should be placed.
type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float32 `json:"x_shift"`
	YShift float32 `json:"y_shift"`
	Scale  float32 `json:"scale"`
}

// InputSticker is the parameter shape for createNewStickerSet / addStickerToSet.
type InputSticker struct {
	Sticker      string        `json:"sticker"`
	Format       string        `json:"format"`
	EmojiList    []string      `json:"emoji_list"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	Keywords     []string      `json:"keywords,omitempty"`
}
