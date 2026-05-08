package tbot

// ChatBoostSource describes the origin of a chat boost. Variant is given by
// Source:
//   - "premium"   — User is set: a Premium subscriber boosted directly.
//   - "gift_code" — User is set: a Premium gift code was redeemed.
//   - "giveaway"  — GiveawayMessageID is set; User and IsUnclaimed may also be.
type ChatBoostSource struct {
	Source             string `json:"source"`
	User               *User  `json:"user,omitempty"`
	GiveawayMessageID  int    `json:"giveaway_message_id,omitempty"`
	PrizeStarCount     int    `json:"prize_star_count,omitempty"`
	IsUnclaimed        bool   `json:"is_unclaimed,omitempty"`
}

// ChatBoostSource Source values.
const (
	ChatBoostSourcePremium  = "premium"
	ChatBoostSourceGiftCode = "gift_code"
	ChatBoostSourceGiveaway = "giveaway"
)

// ChatBoost contains information about an individual boost.
type ChatBoost struct {
	BoostID        string          `json:"boost_id"`
	AddDate        int64           `json:"add_date"`
	ExpirationDate int64           `json:"expiration_date"`
	Source         ChatBoostSource `json:"source"`
}

// ChatBoostUpdated is the update payload sent when a chat boost is added or changed.
type ChatBoostUpdated struct {
	Chat  Chat      `json:"chat"`
	Boost ChatBoost `json:"boost"`
}

// ChatBoostRemoved is the update payload sent when a chat boost is removed.
type ChatBoostRemoved struct {
	Chat       Chat            `json:"chat"`
	BoostID    string          `json:"boost_id"`
	RemoveDate int64           `json:"remove_date"`
	Source     ChatBoostSource `json:"source"`
}

// UserChatBoosts is the list of boosts a user has applied to a chat.
type UserChatBoosts struct {
	Boosts []ChatBoost `json:"boosts"`
}

// ChatBoostAdded is a service message: user boosted the chat.
type ChatBoostAdded struct {
	BoostCount int `json:"boost_count"`
}

// ChatBackground describes a chat background.
type ChatBackground struct {
	Type BackgroundType `json:"type"`
}

// BackgroundFill describes a fill for a chat background. Variant is given by
// Type: "solid", "gradient" or "freeform_gradient".
type BackgroundFill struct {
	Type           string `json:"type"`
	Color          int    `json:"color,omitempty"`
	TopColor       int    `json:"top_color,omitempty"`
	BottomColor    int    `json:"bottom_color,omitempty"`
	RotationAngle  int    `json:"rotation_angle,omitempty"`
	Colors         []int  `json:"colors,omitempty"`
}

// BackgroundFill Type values.
const (
	BackgroundFillSolid            = "solid"
	BackgroundFillGradient         = "gradient"
	BackgroundFillFreeformGradient = "freeform_gradient"
)

// BackgroundType describes a chat background. Variant is given by Type:
//   - "fill"       — Fill, DarkThemeDimming
//   - "wallpaper"  — Document, DarkThemeDimming, IsBlurred, IsMoving
//   - "pattern"    — Document, Fill, Intensity, IsInverted, IsMoving
//   - "chat_theme" — ThemeName
type BackgroundType struct {
	Type             string          `json:"type"`
	Fill             *BackgroundFill `json:"fill,omitempty"`
	DarkThemeDimming int             `json:"dark_theme_dimming,omitempty"`
	Document         *Document       `json:"document,omitempty"`
	IsBlurred        bool            `json:"is_blurred,omitempty"`
	IsMoving         bool            `json:"is_moving,omitempty"`
	Intensity        int             `json:"intensity,omitempty"`
	IsInverted       bool            `json:"is_inverted,omitempty"`
	ThemeName        string          `json:"theme_name,omitempty"`
}

// BackgroundType Type values.
const (
	BackgroundTypeFill      = "fill"
	BackgroundTypeWallpaper = "wallpaper"
	BackgroundTypePattern   = "pattern"
	BackgroundTypeChatTheme = "chat_theme"
)
