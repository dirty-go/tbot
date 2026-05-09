// Package tbot — Telegram Bot API types.
//
// Field names and JSON tags follow the official Bot API 9.6 reference at
// https://core.telegram.org/bots/api. Optional fields use `,omitempty` and are
// pointers or zero-able primitives. Integer fields that carry Telegram chat or
// user identifiers are `int64` because Telegram explicitly notes those values
// can exceed 32-bit precision.
//
// File-size fields are also `int64` since the API documents that they may
// exceed 2^31.
//
// Some fields and types from earlier API versions (Bot API ≤5.x) have been
// renamed by Telegram. To preserve source compatibility for callers of this
// library, the older Go fields/types are retained and annotated with a
// `// Deprecated:` comment that points to the current replacement. New code
// should use the replacement.
//
// Type definitions are split across several files in this package by domain:
//
//   - types.go            — small media/data leaves (PhotoSize, Audio, Poll, …)
//   - types_user.go       — User
//   - types_chat.go       — Chat / ChatFullInfo / ChatMember* / permissions
//   - types_message.go    — Message and reply/origin/preview metadata
//   - types_service.go    — service messages embedded in Message
//   - types_keyboard.go   — keyboards, web-app buttons, callback queries
//   - types_sticker.go    — Sticker / StickerSet
//   - types_reaction.go   — ReactionType / message reaction updates
//   - types_forum.go      — forum topic messages
//   - types_business.go   — business connections and bot rights
//   - types_boost.go      — chat boosts and chat backgrounds
//   - types_giveaway.go   — giveaway service messages
//   - types_payment.go    — invoices, payments, paid media
//   - types_gift.go       — Gifts and unique gifts
//   - types_stars.go      — Star transactions and revenue
//   - types_passport.go   — Telegram Passport
//   - types_inline.go     — inline mode (queries, results, content)
//   - types_input.go      — InputMedia / InputPaidMedia / InputFile
//   - types_bot.go        — bot commands, scopes, names, descriptions, menu button
//   - types_update.go     — Update / WebhookInfo

package tbot

// MessageEntity represents one special entity in a text message — a hashtag,
// bot command, URL, mention, formatted span, or custom emoji.
type MessageEntity struct {
	Type           string `json:"type"`
	Offset         int    `json:"offset"`
	Length         int    `json:"length"`
	URL            string `json:"url,omitempty"`
	User           *User  `json:"user,omitempty"`
	Language       string `json:"language,omitempty"`
	CustomEmojiID  string `json:"custom_emoji_id,omitempty"`
	UnixTime       int64  `json:"unix_time,omitempty"`
	DateTimeFormat string `json:"date_time_format,omitempty"`
}

// PhotoSize represents one size of a photo or a file/sticker thumbnail.
type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// Animation represents an animation file (GIF or H.264/MPEG-4 AVC video without sound).
type Animation struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MIMEType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`

	// Deprecated: legacy alias for MIMEType. Telegram only ever populates
	// mime_type; this field is kept so that older Go callers compile.
	MimeType string `json:"-"`
	// Deprecated: legacy field present in earlier versions of this library.
	// Telegram renamed the JSON tag from "thumb" to "thumbnail" in Bot API 6.6;
	// this Go field is now never populated. Use Thumbnail instead.
	Thumb *PhotoSize `json:"-"`
}

// Audio represents an audio file to be treated as music by Telegram clients.
type Audio struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Duration     int        `json:"duration"`
	Performer    string     `json:"performer,omitempty"`
	Title        string     `json:"title,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MIMEType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
}

// Document represents a general file (as opposed to photos, voice messages,
// audio files, etc.).
type Document struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MIMEType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`

	// Deprecated: see Animation.Thumb.
	Thumb *PhotoSize `json:"-"`
}

// Story represents a story posted in a chat.
type Story struct {
	Chat Chat `json:"chat"`
	ID   int  `json:"id"`
}

// Video represents a video file.
type Video struct {
	FileID         string         `json:"file_id"`
	FileUniqueID   string         `json:"file_unique_id"`
	Width          int            `json:"width"`
	Height         int            `json:"height"`
	Duration       int            `json:"duration"`
	Thumbnail      *PhotoSize     `json:"thumbnail,omitempty"`
	Cover          []PhotoSize    `json:"cover,omitempty"`
	StartTimestamp int            `json:"start_timestamp,omitempty"`
	Qualities      []VideoQuality `json:"qualities,omitempty"`
	FileName       string         `json:"file_name,omitempty"`
	MIMEType       string         `json:"mime_type,omitempty"`
	FileSize       int64          `json:"file_size,omitempty"`

	// Deprecated: legacy alias for MIMEType.
	MimeType string `json:"-"`
	// Deprecated: see Animation.Thumb.
	Thumb *PhotoSize `json:"-"`
}

// VideoQuality describes one of the available qualities of a Video.
type VideoQuality struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Codec        string `json:"codec"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// VideoNote represents a video message (round, recorded with the camera).
type VideoNote struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Length       int        `json:"length"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`

	// Deprecated: see Animation.Thumb.
	Thumb *PhotoSize `json:"-"`
}

// Voice represents a voice note.
type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MIMEType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`

	// Deprecated: legacy alias for MIMEType.
	MimeType string `json:"-"`
}

// Contact represents a phone contact.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
	VCard       string `json:"vcard,omitempty"`
}

// Location represents a point on the map.
type Location struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
}

// Venue represents a venue.
type Venue struct {
	Location        Location `json:"location"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	FoursquareID    string   `json:"foursquare_id,omitempty"`
	FoursquareType  string   `json:"foursquare_type,omitempty"`
	GooglePlaceID   string   `json:"google_place_id,omitempty"`
	GooglePlaceType string   `json:"google_place_type,omitempty"`
}

// Dice represents an animated emoji that displays a random value.
type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

// PollOption is one of the choices in a Poll.
type PollOption struct {
	PersistentID string          `json:"persistent_id"`
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
	VoterCount   int             `json:"voter_count"`
	AddedByUser  *User           `json:"added_by_user,omitempty"`
	AddedByChat  *Chat           `json:"added_by_chat,omitempty"`
	AdditionDate int64           `json:"addition_date,omitempty"`
}

// InputPollOption is the option passed to sendPoll / addPollOption.
type InputPollOption struct {
	Text          string          `json:"text"`
	TextParseMode string          `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity `json:"text_entities,omitempty"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll.
type PollAnswer struct {
	PollID              string   `json:"poll_id"`
	VoterChat           *Chat    `json:"voter_chat,omitempty"`
	User                *User    `json:"user,omitempty"`
	OptionIDs           []int    `json:"option_ids"`
	OptionPersistentIDs []string `json:"option_persistent_ids"`
}

// Poll represents a native Telegram poll.
type Poll struct {
	ID                    string          `json:"id"`
	Question              string          `json:"question"`
	QuestionEntities      []MessageEntity `json:"question_entities,omitempty"`
	Options               []PollOption    `json:"options"`
	TotalVoterCount       int             `json:"total_voter_count"`
	IsClosed              bool            `json:"is_closed"`
	IsAnonymous           bool            `json:"is_anonymous"`
	Type                  string          `json:"type"`
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers"`
	AllowsRevoting        bool            `json:"allows_revoting"`
	CorrectOptionIDs      []int           `json:"correct_option_ids,omitempty"`
	Explanation           string          `json:"explanation,omitempty"`
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"`
	OpenPeriod            int             `json:"open_period,omitempty"`
	CloseDate             int64           `json:"close_date,omitempty"`
	Description           string          `json:"description,omitempty"`
	DescriptionEntities   []MessageEntity `json:"description_entities,omitempty"`

	// Deprecated: replaced by CorrectOptionIDs (now an array). Older code that
	// reads a single value should switch to the new field.
	CorrectOptionID int `json:"-"`
}

// File represents a file ready to be downloaded via getFile.
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

// Game represents a HTML5 game registered via @BotFather.
type Game struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Photo        []PhotoSize     `json:"photo"`
	Text         string          `json:"text,omitempty"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
	Animation    *Animation      `json:"animation,omitempty"`
}

// GameHighScore is one row of a game's high-score table.
type GameHighScore struct {
	Position int  `json:"position"`
	User     User `json:"user"`
	Score    int  `json:"score"`
}

// CallbackGame is a placeholder, currently holds no information. Use
// @BotFather to set up your game.
type CallbackGame struct{}
