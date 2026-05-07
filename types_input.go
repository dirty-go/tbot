package tbot

// InputFile is the Telegram input-file marker. The Bot API accepts three
// shapes:
//   - a previously-uploaded file_id (string)
//   - an HTTPS URL pointing at the file (string)
//   - a multipart upload, referenced from JSON as "attach://<name>"
//
// This package represents InputFile as a string in the JSON payload. For
// multipart uploads, code outside this file constructs the multipart body
// directly (see api.go).
type InputFile = string

// InputMedia describes a media item to send via sendMediaGroup or
// editMessageMedia.
//
// Variant is given by Type — one of "photo", "video", "animation", "audio",
// "document". Set Type and the fields documented for that variant.
type InputMedia struct {
	Type                        string          `json:"type"`
	Media                       string          `json:"media"`
	Thumbnail                   string          `json:"thumbnail,omitempty"`
	Cover                       string          `json:"cover,omitempty"`
	StartTimestamp              int             `json:"start_timestamp,omitempty"`
	Caption                     string          `json:"caption,omitempty"`
	ParseMode                   string          `json:"parse_mode,omitempty"`
	CaptionEntities             []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia       bool            `json:"show_caption_above_media,omitempty"`
	HasSpoiler                  bool            `json:"has_spoiler,omitempty"`
	Width                       int             `json:"width,omitempty"`
	Height                      int             `json:"height,omitempty"`
	Duration                    int             `json:"duration,omitempty"`
	SupportsStreaming           bool            `json:"supports_streaming,omitempty"`
	Performer                   string          `json:"performer,omitempty"`
	Title                       string          `json:"title,omitempty"`
	DisableContentTypeDetection bool            `json:"disable_content_type_detection,omitempty"`
}

// InputMedia Type values.
const (
	InputMediaTypePhoto     = "photo"
	InputMediaTypeVideo     = "video"
	InputMediaTypeAnimation = "animation"
	InputMediaTypeAudio     = "audio"
	InputMediaTypeDocument  = "document"
)

// InputPaidMedia describes a media item to be sent as paid media via
// sendPaidMedia. Variant is given by Type — "photo" or "video".
type InputPaidMedia struct {
	Type              string `json:"type"`
	Media             string `json:"media"`
	Thumbnail         string `json:"thumbnail,omitempty"`
	Cover             string `json:"cover,omitempty"`
	StartTimestamp    int    `json:"start_timestamp,omitempty"`
	Width             int    `json:"width,omitempty"`
	Height            int    `json:"height,omitempty"`
	Duration          int    `json:"duration,omitempty"`
	SupportsStreaming bool   `json:"supports_streaming,omitempty"`
}

// InputProfilePhoto describes a profile photo to upload via setUserProfilePhoto.
// Variant is given by Type — "static" (single photo) or "animated" (a video).
type InputProfilePhoto struct {
	Type              string  `json:"type"`
	Photo             string  `json:"photo,omitempty"`
	Animation         string  `json:"animation,omitempty"`
	MainFrameTimestamp float64 `json:"main_frame_timestamp,omitempty"`
}

// InputStoryContent describes the content of a story to publish via the API.
// Variant is given by Type — "photo" or "video".
type InputStoryContent struct {
	Type            string  `json:"type"`
	Photo           string  `json:"photo,omitempty"`
	Video           string  `json:"video,omitempty"`
	Duration        float64 `json:"duration,omitempty"`
	CoverFrameTimestamp float64 `json:"cover_frame_timestamp,omitempty"`
	IsAnimation     bool    `json:"is_animation,omitempty"`
}
