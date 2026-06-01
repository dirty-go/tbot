package tbot

import (
	"net/url"
	"strconv"
)

// SendMessageOptions collects all common optional fields for send methods into a
// single struct. Callers who prefer struct literals over building a []SendOption
// slice can populate this struct and call Apply to obtain the equivalent slice.
//
// Each non-zero field maps to the equivalent OptXxx functional option. Zero
// values (false, 0, "", nil) are ignored, producing no option in the output.
//
// Example usage:
//
//	opts := tbot.SendMessageOptions{
//	    ParseMode:           "HTML",
//	    DisableNotification: true,
//	}
//	msg, err := c.SendMessage(ctx, chatID, "hello", opts.Apply()...)
//
// Both styles (struct and OptXxx functions) remain fully supported and may be
// freely mixed:
//
//	extra := append(opts.Apply(), tbot.OptProtectContent)
type SendMessageOptions struct {
	// Text / caption formatting.
	ParseMode string // "HTML", "MarkdownV2", or "Markdown"

	// Delivery flags — each maps to the corresponding OptXxx var when true.
	DisableNotification bool
	ProtectContent      bool
	AllowPaidBroadcast  bool

	// Structured reply, preview, and targeting options.
	ReplyParameters    *ReplyParameters
	LinkPreviewOptions *LinkPreviewOptions

	// Entities to attach to the text or caption.
	Entities        []MessageEntity
	CaptionEntities []MessageEntity

	// Thread / business / effect routing.
	MessageThreadID      int
	BusinessConnectionID string
	MessageEffectID      string

	// Reply markup — store a pre-built SendOption from OptInlineKeyboardMarkup,
	// OptReplyKeyboardMarkup, OptReplyKeyboardRemove, or OptForceReply.
	ReplyMarkup SendOption

	// Media caption.
	Caption string

	// Caption display.
	ShowCaptionAboveMedia bool
	HasSpoiler            bool

	// Media dimensions and duration.
	Duration  int
	Width     int
	Height    int
	Thumbnail string // file_id or URL

	// Audio / video extras.
	SupportsStreaming           bool
	DisableContentTypeDetection bool
	Performer                   string
	Title                       string

	// Location extras.
	Latitude             float64
	Longitude            float64
	HorizontalAccuracy   float64
	Heading              int
	ProximityAlertRadius int
	LivePeriod           int

	// Venue extras.
	Address         string
	FoursquareID    string
	FoursquareType  string
	GooglePlaceID   string
	GooglePlaceType string

	// Contact extras.
	PhoneNumber string
	FirstName   string
	LastName    string
	VCard       string

	// Poll extras.
	Question             string
	Options              []InputPollOption
	IsAnonymous          bool
	Type                 string // "quiz" or "regular"
	AllowsMultipleAnswers bool
	Explanation          string
	ExplanationParseMode string
	ExplanationEntities  []MessageEntity
	OpenPeriod           int
	CloseDate            int64
	IsClosed             bool

	// Dice extras.
	Emoji string

	// Edit targeting (paired with chatID/messageID or inlineMessageID).
	ChatID          ChatID
	MessageID       int
	InlineMessageID string
}

// Apply converts the populated fields of o into the equivalent []SendOption
// slice. Fields at their zero value (false, 0, "", nil) produce no option.
// The result is a fresh slice each call and is safe to append to.
func (o SendMessageOptions) Apply() []SendOption {
	var opts []SendOption

	add := func(opt SendOption) { opts = append(opts, opt) }

	if o.ParseMode != "" {
		add(func(v url.Values) { v.Set("parse_mode", o.ParseMode) })
	}
	if o.DisableNotification {
		add(OptDisableNotification)
	}
	if o.ProtectContent {
		add(OptProtectContent)
	}
	if o.AllowPaidBroadcast {
		add(OptAllowPaidBroadcast)
	}
	if o.ReplyParameters != nil {
		add(OptReplyParameters(o.ReplyParameters))
	}
	if o.LinkPreviewOptions != nil {
		add(OptLinkPreviewOptions(o.LinkPreviewOptions))
	}
	if len(o.Entities) > 0 {
		add(func(v url.Values) { v.Set("entities", structString(o.Entities)) })
	}
	if len(o.CaptionEntities) > 0 {
		add(func(v url.Values) { v.Set("caption_entities", structString(o.CaptionEntities)) })
	}
	if o.MessageThreadID != 0 {
		add(OptMessageThreadID(o.MessageThreadID))
	}
	if o.BusinessConnectionID != "" {
		add(OptBusinessConnectionID(o.BusinessConnectionID))
	}
	if o.MessageEffectID != "" {
		add(OptMessageEffectID(o.MessageEffectID))
	}
	if o.ReplyMarkup != nil {
		add(o.ReplyMarkup)
	}
	if o.Caption != "" {
		add(OptCaption(o.Caption))
	}
	if o.ShowCaptionAboveMedia {
		add(OptShowCaptionAboveMedia)
	}
	if o.HasSpoiler {
		add(OptHasSpoiler)
	}
	if o.Duration != 0 {
		add(OptDuration(o.Duration))
	}
	if o.Width != 0 {
		add(OptWidth(o.Width))
	}
	if o.Height != 0 {
		add(OptHeight(o.Height))
	}
	if o.Thumbnail != "" {
		add(OptThumbnail(o.Thumbnail))
	}
	if o.SupportsStreaming {
		add(OptSupportsStreaming)
	}
	if o.DisableContentTypeDetection {
		add(OptDisableContentTypeDetection)
	}
	if o.Performer != "" {
		add(OptPerformer(o.Performer))
	}
	if o.Title != "" {
		add(OptTitle(o.Title))
	}
	if o.Latitude != 0 {
		add(func(v url.Values) { v.Set("latitude", formatFloat(o.Latitude)) })
	}
	if o.Longitude != 0 {
		add(func(v url.Values) { v.Set("longitude", formatFloat(o.Longitude)) })
	}
	if o.HorizontalAccuracy != 0 {
		add(OptHorizontalAccuracy(o.HorizontalAccuracy))
	}
	if o.Heading != 0 {
		add(OptHeading(o.Heading))
	}
	if o.ProximityAlertRadius != 0 {
		add(OptProximityAlertRadius(o.ProximityAlertRadius))
	}
	if o.LivePeriod != 0 {
		add(OptLivePeriod(o.LivePeriod))
	}
	if o.Address != "" {
		add(func(v url.Values) { v.Set("address", o.Address) })
	}
	if o.FoursquareID != "" {
		add(func(v url.Values) { v.Set("foursquare_id", o.FoursquareID) })
	}
	if o.FoursquareType != "" {
		add(func(v url.Values) { v.Set("foursquare_type", o.FoursquareType) })
	}
	if o.GooglePlaceID != "" {
		add(func(v url.Values) { v.Set("google_place_id", o.GooglePlaceID) })
	}
	if o.GooglePlaceType != "" {
		add(func(v url.Values) { v.Set("google_place_type", o.GooglePlaceType) })
	}
	if o.PhoneNumber != "" {
		add(func(v url.Values) { v.Set("phone_number", o.PhoneNumber) })
	}
	if o.FirstName != "" {
		add(func(v url.Values) { v.Set("first_name", o.FirstName) })
	}
	if o.LastName != "" {
		add(func(v url.Values) { v.Set("last_name", o.LastName) })
	}
	if o.VCard != "" {
		add(func(v url.Values) { v.Set("vcard", o.VCard) })
	}
	if o.Question != "" {
		add(func(v url.Values) { v.Set("question", o.Question) })
	}
	if len(o.Options) > 0 {
		add(func(v url.Values) { v.Set("options", structString(o.Options)) })
	}
	if o.IsAnonymous {
		add(func(v url.Values) { v.Set("is_anonymous", "true") })
	}
	if o.Type != "" {
		add(func(v url.Values) { v.Set("type", o.Type) })
	}
	if o.AllowsMultipleAnswers {
		add(func(v url.Values) { v.Set("allows_multiple_answers", "true") })
	}
	if o.Explanation != "" {
		add(func(v url.Values) { v.Set("explanation", o.Explanation) })
	}
	if o.ExplanationParseMode != "" {
		add(func(v url.Values) { v.Set("explanation_parse_mode", o.ExplanationParseMode) })
	}
	if len(o.ExplanationEntities) > 0 {
		add(func(v url.Values) { v.Set("explanation_entities", structString(o.ExplanationEntities)) })
	}
	if o.OpenPeriod != 0 {
		add(func(v url.Values) { v.Set("open_period", strconv.Itoa(o.OpenPeriod)) })
	}
	if o.CloseDate != 0 {
		add(func(v url.Values) { v.Set("close_date", strconv.FormatInt(o.CloseDate, 10)) })
	}
	if o.IsClosed {
		add(func(v url.Values) { v.Set("is_closed", "true") })
	}
	if o.Emoji != "" {
		add(func(v url.Values) { v.Set("emoji", o.Emoji) })
	}
	if !o.ChatID.IsZero() {
		add(OptChatID(o.ChatID))
	}
	if o.MessageID != 0 {
		add(OptMessageID(o.MessageID))
	}
	if o.InlineMessageID != "" {
		add(OptInlineMessageID(o.InlineMessageID))
	}

	return opts
}
