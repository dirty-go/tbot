package tbot

import (
	"net/url"
	"strconv"
)

// SendOption mutates an outgoing request's form values. Options are the
// discoverable, composable way to set the many optional Bot API fields without
// exploding method signatures.
type SendOption func(url.Values)

// Formatting and delivery options.
var (
	// OptParseModeHTML formats text/caption as HTML.
	OptParseModeHTML = func(r url.Values) { r.Set("parse_mode", "HTML") }
	// OptParseModeMarkdown formats text/caption as MarkdownV2.
	OptParseModeMarkdown = func(r url.Values) { r.Set("parse_mode", "MarkdownV2") }
	// OptDisableNotification sends the message silently.
	OptDisableNotification = func(r url.Values) { r.Set("disable_notification", "true") }
	// OptProtectContent forbids forwarding and saving of the sent content.
	OptProtectContent = func(r url.Values) { r.Set("protect_content", "true") }
	// OptAllowPaidBroadcast allows up to 1000 messages/second against a fee.
	OptAllowPaidBroadcast = func(r url.Values) { r.Set("allow_paid_broadcast", "true") }
	// OptDisableWebPagePreview disables link previews for the message.
	//
	// Deprecated: prefer OptLinkPreviewOptions, which the Bot API now favours.
	OptDisableWebPagePreview = func(r url.Values) { r.Set("disable_web_page_preview", "true") }
	// OptShowCaptionAboveMedia renders the caption above the media.
	OptShowCaptionAboveMedia = func(r url.Values) { r.Set("show_caption_above_media", "true") }
	// OptHasSpoiler covers the media with a spoiler animation.
	OptHasSpoiler = func(r url.Values) { r.Set("has_spoiler", "true") }
	// OptSupportsStreaming marks an uploaded video as suitable for streaming.
	OptSupportsStreaming = func(r url.Values) { r.Set("supports_streaming", "true") }
	// OptDisableContentTypeDetection disables server-side content-type sniffing
	// for an uploaded document.
	OptDisableContentTypeDetection = func(r url.Values) { r.Set("disable_content_type_detection", "true") }
	// OptSendingWithoutReply sends the message even if the replied-to message
	// is missing.
	OptSendingWithoutReply = func(r url.Values) { r.Set("allow_sending_without_reply", "true") }
	// OptIsBig makes a SetMessageReaction reaction animation big.
	OptIsBig = func(r url.Values) { r.Set("is_big", "true") }
)

// OptMessageThreadID targets a forum topic / thread within a chat.
func OptMessageThreadID(id int) SendOption {
	return func(r url.Values) { r.Set("message_thread_id", strconv.Itoa(id)) }
}

// OptBusinessConnectionID sends on behalf of a connected business account.
func OptBusinessConnectionID(id string) SendOption {
	return func(r url.Values) { r.Set("business_connection_id", id) }
}

// OptMessageEffectID applies a message effect (private chats only).
func OptMessageEffectID(id string) SendOption {
	return func(r url.Values) { r.Set("message_effect_id", id) }
}

// OptReplyToMessageID replies to a specific message.
//
// Deprecated: prefer OptReplyParameters, which also supports cross-chat
// replies and quotes.
func OptReplyToMessageID(id int) SendOption {
	return func(r url.Values) { r.Set("reply_to_message_id", strconv.Itoa(id)) }
}

// OptReplyParameters describes the message being replied to.
func OptReplyParameters(p *ReplyParameters) SendOption {
	return func(r url.Values) { r.Set("reply_parameters", structString(p)) }
}

// OptLinkPreviewOptions configures link-preview generation.
func OptLinkPreviewOptions(o *LinkPreviewOptions) SendOption {
	return func(r url.Values) { r.Set("link_preview_options", structString(o)) }
}

// OptCaption sets a media caption (0-1024 chars after entity parsing).
func OptCaption(text string) SendOption {
	return func(r url.Values) { r.Set("caption", text) }
}

// OptThumbnail sets a thumbnail by file_id, URL, or "attach://" reference.
// Thumbnails cannot themselves be uploaded through this option.
func OptThumbnail(ref string) SendOption {
	return func(r url.Values) { r.Set("thumbnail", ref) }
}

// OptDuration sets a media duration in seconds.
func OptDuration(seconds int) SendOption {
	return func(r url.Values) { r.Set("duration", strconv.Itoa(seconds)) }
}

// OptWidth sets a media width in pixels.
func OptWidth(px int) SendOption {
	return func(r url.Values) { r.Set("width", strconv.Itoa(px)) }
}

// OptHeight sets a media height in pixels.
func OptHeight(px int) SendOption {
	return func(r url.Values) { r.Set("height", strconv.Itoa(px)) }
}

// OptPerformer sets the performer of an audio file.
func OptPerformer(name string) SendOption {
	return func(r url.Values) { r.Set("performer", name) }
}

// OptTitle sets the title of an audio file or venue.
func OptTitle(title string) SendOption {
	return func(r url.Values) { r.Set("title", title) }
}

// OptLivePeriod sets the period (seconds) a live location stays live.
func OptLivePeriod(seconds int) SendOption {
	return func(r url.Values) { r.Set("live_period", strconv.Itoa(seconds)) }
}

// OptHorizontalAccuracy sets a location's accuracy radius in metres.
func OptHorizontalAccuracy(metres float64) SendOption {
	return func(r url.Values) {
		r.Set("horizontal_accuracy", strconv.FormatFloat(metres, 'f', -1, 64))
	}
}

// OptHeading sets a live location's direction of movement in degrees (1-360).
func OptHeading(degrees int) SendOption {
	return func(r url.Values) { r.Set("heading", strconv.Itoa(degrees)) }
}

// OptProximityAlertRadius sets the proximity-alert radius in metres.
func OptProximityAlertRadius(metres int) SendOption {
	return func(r url.Values) { r.Set("proximity_alert_radius", strconv.Itoa(metres)) }
}

// OptChatID targets a chat message to edit (paired with OptMessageID). Edit
// methods need either this pair or OptInlineMessageID.
func OptChatID(chatID ChatID) SendOption {
	return func(r url.Values) { r.Set("chat_id", chatID.String()) }
}

// OptMessageID targets a specific message to edit.
func OptMessageID(id int) SendOption {
	return func(r url.Values) { r.Set("message_id", strconv.Itoa(id)) }
}

// OptInlineMessageID targets an inline message to edit.
func OptInlineMessageID(id string) SendOption {
	return func(r url.Values) { r.Set("inline_message_id", id) }
}

// Reply-markup options.
var (
	// OptReplyKeyboardRemove removes a custom keyboard.
	OptReplyKeyboardRemove = func(r url.Values) {
		r.Set("reply_markup", structString(&ReplyKeyboardRemove{RemoveKeyboard: true}))
	}
	// OptReplyKeyboardRemoveSelective removes the keyboard for targeted users.
	OptReplyKeyboardRemoveSelective = func(r url.Values) {
		r.Set("reply_markup", structString(&ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}))
	}
	// OptForceReply forces clients to show a reply interface.
	OptForceReply = func(r url.Values) {
		r.Set("reply_markup", structString(&ForceReply{ForceReply: true}))
	}
	// OptForceReplySelective forces a reply interface for targeted users.
	OptForceReplySelective = func(r url.Values) {
		r.Set("reply_markup", structString(&ForceReply{ForceReply: true, Selective: true}))
	}
)

// OptInlineKeyboardMarkup attaches an inline keyboard.
func OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup) SendOption {
	return func(r url.Values) { r.Set("reply_markup", structString(markup)) }
}

// OptReplyKeyboardMarkup attaches a custom reply keyboard.
func OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup) SendOption {
	return func(r url.Values) { r.Set("reply_markup", structString(markup)) }
}
