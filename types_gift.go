package tbot

// Gift represents a regular Telegram gift available for purchase.
type Gift struct {
	ID                     string          `json:"id"`
	Sticker                Sticker         `json:"sticker"`
	StarCount              int             `json:"star_count"`
	UpgradeStarCount       int             `json:"upgrade_star_count,omitempty"`
	IsPremium              bool            `json:"is_premium,omitempty"`
	HasColors              bool            `json:"has_colors,omitempty"`
	TotalCount             int             `json:"total_count,omitempty"`
	RemainingCount         int             `json:"remaining_count,omitempty"`
	PersonalTotalCount     int             `json:"personal_total_count,omitempty"`
	PersonalRemainingCount int             `json:"personal_remaining_count,omitempty"`
	Background             *GiftBackground `json:"background,omitempty"`
	UniqueGiftVariantCount int             `json:"unique_gift_variant_count,omitempty"`
	PublisherChat          *Chat           `json:"publisher_chat,omitempty"`
}

// GiftBackground describes the background of a Gift. Telegram has not yet
// formalized this object's shape; treat it as opaque metadata.
type GiftBackground struct {
	Fill *BackgroundFill `json:"fill,omitempty"`
}

// Gifts is a list of Gift.
type Gifts struct {
	Gifts []Gift `json:"gifts"`
}

// UniqueGiftModel is the visual model used in a unique gift.
type UniqueGiftModel struct {
	Name           string  `json:"name"`
	Sticker        Sticker `json:"sticker"`
	RarityPerMille int     `json:"rarity_per_mille"`
	Rarity         string  `json:"rarity,omitempty"`
}

// UniqueGiftSymbol is the symbol used in a unique gift.
type UniqueGiftSymbol struct {
	Name           string  `json:"name"`
	Sticker        Sticker `json:"sticker"`
	RarityPerMille int     `json:"rarity_per_mille"`
}

// UniqueGiftBackdropColors describes the colors of a unique gift backdrop.
type UniqueGiftBackdropColors struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	SymbolColor int `json:"symbol_color"`
	TextColor   int `json:"text_color"`
}

// UniqueGiftBackdrop is the backdrop of a unique gift.
type UniqueGiftBackdrop struct {
	Name           string                   `json:"name"`
	Colors         UniqueGiftBackdropColors `json:"colors"`
	RarityPerMille int                      `json:"rarity_per_mille"`
}

// UniqueGift describes a unique gift that was upgraded from a regular one.
type UniqueGift struct {
	GiftID            string             `json:"gift_id"`
	BaseName          string             `json:"base_name"`
	Name              string             `json:"name"`
	Number            int                `json:"number"`
	Model             UniqueGiftModel    `json:"model"`
	Symbol            UniqueGiftSymbol   `json:"symbol"`
	Backdrop          UniqueGiftBackdrop `json:"backdrop"`
	IsPremium         bool               `json:"is_premium,omitempty"`
	IsBurned          bool               `json:"is_burned,omitempty"`
	IsFromBlockchain  bool               `json:"is_from_blockchain,omitempty"`
	Colors            *UniqueGiftColors  `json:"colors,omitempty"`
	PublisherChat     *Chat              `json:"publisher_chat,omitempty"`
}

// UniqueGiftColors describes the color scheme of a unique gift.
type UniqueGiftColors struct {
	ModelCustomEmojiID    string `json:"model_custom_emoji_id"`
	SymbolCustomEmojiID   string `json:"symbol_custom_emoji_id"`
	LightThemeMainColor   int    `json:"light_theme_main_color"`
	LightThemeOtherColors []int  `json:"light_theme_other_colors"`
	DarkThemeMainColor    int    `json:"dark_theme_main_color"`
	DarkThemeOtherColors  []int  `json:"dark_theme_other_colors"`
}

// GiftInfo is a service message about a regular gift sent or received.
type GiftInfo struct {
	Gift                    Gift            `json:"gift"`
	OwnedGiftID             string          `json:"owned_gift_id,omitempty"`
	ConvertStarCount        int             `json:"convert_star_count,omitempty"`
	PrepaidUpgradeStarCount int             `json:"prepaid_upgrade_star_count,omitempty"`
	IsUpgradeSeparate       bool            `json:"is_upgrade_separate,omitempty"`
	CanBeUpgraded           bool            `json:"can_be_upgraded,omitempty"`
	Text                    string          `json:"text,omitempty"`
	Entities                []MessageEntity `json:"entities,omitempty"`
	IsPrivate               bool            `json:"is_private,omitempty"`
	UniqueGiftNumber        int             `json:"unique_gift_number,omitempty"`
}

// UniqueGiftInfo is a service message about a unique gift sent or received.
type UniqueGiftInfo struct {
	Gift                UniqueGift `json:"gift"`
	Origin              string     `json:"origin"`
	LastResaleCurrency  string     `json:"last_resale_currency,omitempty"`
	LastResaleAmount    int64      `json:"last_resale_amount,omitempty"`
	OwnedGiftID         string     `json:"owned_gift_id,omitempty"`
	TransferStarCount   int        `json:"transfer_star_count,omitempty"`
	NextTransferDate    int64      `json:"next_transfer_date,omitempty"`
}

// OwnedGift represents a gift owned by the bot or a chat. Variant is given by
// Type: "regular" → Gift is set; "unique" → UniqueGift is set.
type OwnedGift struct {
	Type                    string          `json:"type"`
	Gift                    *Gift           `json:"gift,omitempty"`
	UniqueGift              *UniqueGift     `json:"unique_gift,omitempty"`
	OwnedGiftID             string          `json:"owned_gift_id,omitempty"`
	SenderUser              *User           `json:"sender_user,omitempty"`
	SendDate                int64           `json:"send_date"`
	Text                    string          `json:"text,omitempty"`
	Entities                []MessageEntity `json:"entities,omitempty"`
	IsPrivate               bool            `json:"is_private,omitempty"`
	IsSaved                 bool            `json:"is_saved,omitempty"`
	CanBeUpgraded           bool            `json:"can_be_upgraded,omitempty"`
	CanBeTransferred        bool            `json:"can_be_transferred,omitempty"`
	WasRefunded             bool            `json:"was_refunded,omitempty"`
	ConvertStarCount        int             `json:"convert_star_count,omitempty"`
	PrepaidUpgradeStarCount int             `json:"prepaid_upgrade_star_count,omitempty"`
	IsUpgradeSeparate       bool            `json:"is_upgrade_separate,omitempty"`
	UniqueGiftNumber        int             `json:"unique_gift_number,omitempty"`
	TransferStarCount       int             `json:"transfer_star_count,omitempty"`
	NextTransferDate        int64           `json:"next_transfer_date,omitempty"`
}

// OwnedGift Type values.
const (
	OwnedGiftTypeRegular = "regular"
	OwnedGiftTypeUnique  = "unique"
)

// OwnedGifts is the response to getOwnedGifts.
type OwnedGifts struct {
	TotalCount int         `json:"total_count"`
	Gifts      []OwnedGift `json:"gifts"`
	NextOffset string      `json:"next_offset,omitempty"`
}

// AcceptedGiftTypes describes which gift types are accepted by a chat.
type AcceptedGiftTypes struct {
	UnlimitedGifts       bool `json:"unlimited_gifts"`
	LimitedGifts         bool `json:"limited_gifts"`
	UniqueGifts          bool `json:"unique_gifts"`
	PremiumSubscription  bool `json:"premium_subscription"`
	GiftsFromChannels    bool `json:"gifts_from_channels"`
}

// UserRating is the rating of a user (used in ChatFullInfo).
type UserRating struct {
	Level               int `json:"level"`
	Rating              int `json:"rating"`
	CurrentLevelRating  int `json:"current_level_rating"`
	NextLevelRating     int `json:"next_level_rating,omitempty"`
}
