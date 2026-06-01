package tbot

import "encoding/json"

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

// MarshalJSON emits only fields valid for the selected boost source variant.
func (s ChatBoostSource) MarshalJSON() ([]byte, error) {
	type payload struct {
		Source            string `json:"source"`
		User              *User  `json:"user,omitempty"`
		GiveawayMessageID int    `json:"giveaway_message_id,omitempty"`
		PrizeStarCount    int    `json:"prize_star_count,omitempty"`
		IsUnclaimed       bool   `json:"is_unclaimed,omitempty"`
	}
	out := payload{Source: s.Source}
	switch s.Source {
	case ChatBoostSourcePremium, ChatBoostSourceGiftCode:
		out.User = s.User
	case ChatBoostSourceGiveaway:
		out.GiveawayMessageID = s.GiveawayMessageID
		out.User = s.User
		out.PrizeStarCount = s.PrizeStarCount
		out.IsUnclaimed = s.IsUnclaimed
	}
	return json.Marshal(out)
}

// ChatBoostSource Source values.
const (
	ChatBoostSourcePremium  = "premium"
	ChatBoostSourceGiftCode = "gift_code"
	ChatBoostSourceGiveaway = "giveaway"
)

// IsPremium reports whether the boost was purchased with a Telegram Premium subscription.
func (s ChatBoostSource) IsPremium() bool { return s.Source == ChatBoostSourcePremium }

// IsGiftCode reports whether the boost came from a redeemed Premium gift code.
func (s ChatBoostSource) IsGiftCode() bool { return s.Source == ChatBoostSourceGiftCode }

// IsGiveaway reports whether the boost came from a giveaway.
func (s ChatBoostSource) IsGiveaway() bool { return s.Source == ChatBoostSourceGiveaway }

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

// MarshalJSON emits only fields valid for the selected fill variant.
func (b BackgroundFill) MarshalJSON() ([]byte, error) {
	type payload struct {
		Type          string `json:"type"`
		Color         int    `json:"color,omitempty"`
		TopColor      int    `json:"top_color,omitempty"`
		BottomColor   int    `json:"bottom_color,omitempty"`
		RotationAngle int    `json:"rotation_angle,omitempty"`
		Colors        []int  `json:"colors,omitempty"`
	}
	out := payload{Type: b.Type}
	switch b.Type {
	case BackgroundFillSolid:
		out.Color = b.Color
	case BackgroundFillGradient:
		out.TopColor = b.TopColor
		out.BottomColor = b.BottomColor
		out.RotationAngle = b.RotationAngle
	case BackgroundFillFreeformGradient:
		out.Colors = b.Colors
	}
	return json.Marshal(out)
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

// MarshalJSON emits only fields valid for the selected background variant.
func (b BackgroundType) MarshalJSON() ([]byte, error) {
	type payload struct {
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
	out := payload{Type: b.Type}
	switch b.Type {
	case BackgroundTypeFill:
		out.Fill = b.Fill
		out.DarkThemeDimming = b.DarkThemeDimming
	case BackgroundTypeWallpaper:
		out.Document = b.Document
		out.DarkThemeDimming = b.DarkThemeDimming
		out.IsBlurred = b.IsBlurred
		out.IsMoving = b.IsMoving
	case BackgroundTypePattern:
		out.Document = b.Document
		out.Fill = b.Fill
		out.Intensity = b.Intensity
		out.IsInverted = b.IsInverted
		out.IsMoving = b.IsMoving
	case BackgroundTypeChatTheme:
		out.ThemeName = b.ThemeName
	}
	return json.Marshal(out)
}

// BackgroundType Type values.
const (
	BackgroundTypeFill      = "fill"
	BackgroundTypeWallpaper = "wallpaper"
	BackgroundTypePattern   = "pattern"
	BackgroundTypeChatTheme = "chat_theme"
)
