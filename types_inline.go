package tbot

// InlineQuery represents an incoming inline query.
type InlineQuery struct {
	ID       string    `json:"id"`
	From     User      `json:"from"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
	ChatType string    `json:"chat_type,omitempty"`
	Location *Location `json:"location,omitempty"`
}

// InlineQueryResultsButton describes a button to be shown above inline query
// results, e.g. an authorization Web App.
type InlineQueryResultsButton struct {
	Text           string      `json:"text"`
	WebApp         *WebAppInfo `json:"web_app,omitempty"`
	StartParameter string      `json:"start_parameter,omitempty"`
}

// ChosenInlineResult is the result chosen by the user from an InlineQuery.
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`
	From            User      `json:"from"`
	Location        *Location `json:"location,omitempty"`
	InlineMessageID string    `json:"inline_message_id,omitempty"`
	Query           string    `json:"query"`
}

// SentWebAppMessage describes the message sent by the bot via a Web App.
type SentWebAppMessage struct {
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// PreparedInlineMessage is a prepared inline message returned by
// savePreparedInlineMessage.
type PreparedInlineMessage struct {
	ID             string `json:"id"`
	ExpirationDate int64  `json:"expiration_date"`
}

// InlineQueryResult is the union shape used to send an inline result.
//
// Telegram defines many concrete InlineQueryResult* subtypes (Article, Photo,
// Gif, Mpeg4Gif, Video, Audio, Voice, Document, Location, Venue, Contact,
// Game, plus Cached* variants for previously uploaded files). Rather than
// model each one as a separate struct, we expose a single shape that holds
// the discriminator (Type), the result identifier (ID), and the union of all
// fields used by the variants. Set Type and the fields documented for that
// variant; leave others zero.
//
// Reference: https://core.telegram.org/bots/api#inlinequeryresult
type InlineQueryResult struct {
	Type string `json:"type"`
	ID   string `json:"id"`

	// Common fields.
	Title               string                `json:"title,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	CaptionEntities     []MessageEntity       `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                `json:"show_caption_above_media,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`

	// Article.
	URL         string `json:"url,omitempty"`
	HideURL     bool   `json:"hide_url,omitempty"`
	Description string `json:"description,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	ThumbnailWidth int  `json:"thumbnail_width,omitempty"`
	ThumbnailHeight int `json:"thumbnail_height,omitempty"`

	// Photo / Gif / Video / Audio / Voice / Document URLs and dimensions.
	PhotoURL     string `json:"photo_url,omitempty"`
	PhotoFileID  string `json:"photo_file_id,omitempty"`
	PhotoWidth   int    `json:"photo_width,omitempty"`
	PhotoHeight  int    `json:"photo_height,omitempty"`
	GifURL       string `json:"gif_url,omitempty"`
	GifFileID    string `json:"gif_file_id,omitempty"`
	GifWidth     int    `json:"gif_width,omitempty"`
	GifHeight    int    `json:"gif_height,omitempty"`
	GifDuration  int    `json:"gif_duration,omitempty"`
	Mpeg4URL     string `json:"mpeg4_url,omitempty"`
	Mpeg4FileID  string `json:"mpeg4_file_id,omitempty"`
	Mpeg4Width   int    `json:"mpeg4_width,omitempty"`
	Mpeg4Height  int    `json:"mpeg4_height,omitempty"`
	Mpeg4Duration int   `json:"mpeg4_duration,omitempty"`
	VideoURL     string `json:"video_url,omitempty"`
	VideoFileID  string `json:"video_file_id,omitempty"`
	VideoWidth   int    `json:"video_width,omitempty"`
	VideoHeight  int    `json:"video_height,omitempty"`
	VideoDuration int   `json:"video_duration,omitempty"`
	AudioURL     string `json:"audio_url,omitempty"`
	AudioFileID  string `json:"audio_file_id,omitempty"`
	AudioDuration int   `json:"audio_duration,omitempty"`
	VoiceURL     string `json:"voice_url,omitempty"`
	VoiceFileID  string `json:"voice_file_id,omitempty"`
	VoiceDuration int   `json:"voice_duration,omitempty"`
	DocumentURL  string `json:"document_url,omitempty"`
	DocumentFileID string `json:"document_file_id,omitempty"`
	StickerFileID string `json:"sticker_file_id,omitempty"`
	MIMEType     string `json:"mime_type,omitempty"`

	// Audio / video / document.
	Performer string `json:"performer,omitempty"`

	// Thumbnails for non-cached results.
	ThumbnailMimeType string `json:"thumbnail_mime_type,omitempty"`

	// Location / Venue.
	Latitude             float64 `json:"latitude,omitempty"`
	Longitude            float64 `json:"longitude,omitempty"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
	Address              string  `json:"address,omitempty"`
	FoursquareID         string  `json:"foursquare_id,omitempty"`
	FoursquareType       string  `json:"foursquare_type,omitempty"`
	GooglePlaceID        string  `json:"google_place_id,omitempty"`
	GooglePlaceType      string  `json:"google_place_type,omitempty"`

	// Contact.
	PhoneNumber string `json:"phone_number,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	VCard       string `json:"vcard,omitempty"`

	// Game.
	GameShortName string `json:"game_short_name,omitempty"`
}

// InlineQueryResult Type values.
const (
	InlineQueryResultTypeArticle  = "article"
	InlineQueryResultTypePhoto    = "photo"
	InlineQueryResultTypeGif      = "gif"
	InlineQueryResultTypeMpeg4Gif = "mpeg4_gif"
	InlineQueryResultTypeVideo    = "video"
	InlineQueryResultTypeAudio    = "audio"
	InlineQueryResultTypeVoice    = "voice"
	InlineQueryResultTypeDocument = "document"
	InlineQueryResultTypeLocation = "location"
	InlineQueryResultTypeVenue    = "venue"
	InlineQueryResultTypeContact  = "contact"
	InlineQueryResultTypeGame     = "game"
	InlineQueryResultTypeSticker  = "sticker"
)

// InputMessageContent describes the content of a message to be sent via the
// bot as the result of an inline query. Variant is implied by which fields
// are set; in practice callers populate the fields for one of:
//   - text       — MessageText, ParseMode, Entities, LinkPreviewOptions
//   - location   — Latitude, Longitude, HorizontalAccuracy, LivePeriod, …
//   - venue      — Latitude, Longitude, Title, Address, FoursquareID, …
//   - contact    — PhoneNumber, FirstName, LastName, VCard
//   - invoice    — Title, Description, Payload, …
type InputMessageContent struct {
	// Text.
	MessageText        string              `json:"message_text,omitempty"`
	ParseMode          string              `json:"parse_mode,omitempty"`
	Entities           []MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`

	// Location / Venue.
	Latitude             float64 `json:"latitude,omitempty"`
	Longitude            float64 `json:"longitude,omitempty"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
	Title                string  `json:"title,omitempty"`
	Address              string  `json:"address,omitempty"`
	FoursquareID         string  `json:"foursquare_id,omitempty"`
	FoursquareType       string  `json:"foursquare_type,omitempty"`
	GooglePlaceID        string  `json:"google_place_id,omitempty"`
	GooglePlaceType      string  `json:"google_place_type,omitempty"`

	// Contact.
	PhoneNumber string `json:"phone_number,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	VCard       string `json:"vcard,omitempty"`

	// Invoice (InputInvoiceMessageContent).
	Description               string         `json:"description,omitempty"`
	Payload                   string         `json:"payload,omitempty"`
	ProviderToken             string         `json:"provider_token,omitempty"`
	Currency                  string         `json:"currency,omitempty"`
	Prices                    []LabeledPrice `json:"prices,omitempty"`
	MaxTipAmount              int            `json:"max_tip_amount,omitempty"`
	SuggestedTipAmounts       []int          `json:"suggested_tip_amounts,omitempty"`
	ProviderData              string         `json:"provider_data,omitempty"`
	PhotoURL                  string         `json:"photo_url,omitempty"`
	PhotoSize                 int            `json:"photo_size,omitempty"`
	PhotoWidth                int            `json:"photo_width,omitempty"`
	PhotoHeight               int            `json:"photo_height,omitempty"`
	NeedName                  bool           `json:"need_name,omitempty"`
	NeedPhoneNumber           bool           `json:"need_phone_number,omitempty"`
	NeedEmail                 bool           `json:"need_email,omitempty"`
	NeedShippingAddress       bool           `json:"need_shipping_address,omitempty"`
	SendPhoneNumberToProvider bool           `json:"send_phone_number_to_provider,omitempty"`
	SendEmailToProvider       bool           `json:"send_email_to_provider,omitempty"`
	IsFlexible                bool           `json:"is_flexible,omitempty"`
}
