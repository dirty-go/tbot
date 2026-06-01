package tbot

import (
	"net/url"
	"testing"
)

// applyOpts runs a []SendOption against a fresh url.Values and returns the result.
func applyOpts(opts []SendOption) url.Values {
	v := url.Values{}
	for _, o := range opts {
		o(v)
	}
	return v
}

func TestSendMessageOptions_Apply_Empty(t *testing.T) {
	opts := SendMessageOptions{}
	if got := opts.Apply(); len(got) != 0 {
		t.Errorf("Apply() on zero value: want 0 options, got %d", len(got))
	}
}

func TestSendMessageOptions_Apply_ParseMode(t *testing.T) {
	cases := []struct {
		name      string
		parseMode string
		want      string
	}{
		{name: "html", parseMode: "HTML", want: "HTML"},
		{name: "markdownv2", parseMode: "MarkdownV2", want: "MarkdownV2"},
		{name: "empty_omitted", parseMode: "", want: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := SendMessageOptions{ParseMode: tc.parseMode}
			v := applyOpts(opts.Apply())
			if got := v.Get("parse_mode"); got != tc.want {
				t.Errorf("parse_mode = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSendMessageOptions_Apply_BoolFlags(t *testing.T) {
	cases := []struct {
		name  string
		opts  SendMessageOptions
		field string
		want  string
	}{
		{
			name:  "disable_notification",
			opts:  SendMessageOptions{DisableNotification: true},
			field: "disable_notification",
			want:  "true",
		},
		{
			name:  "protect_content",
			opts:  SendMessageOptions{ProtectContent: true},
			field: "protect_content",
			want:  "true",
		},
		{
			name:  "allow_paid_broadcast",
			opts:  SendMessageOptions{AllowPaidBroadcast: true},
			field: "allow_paid_broadcast",
			want:  "true",
		},
		{
			name:  "show_caption_above_media",
			opts:  SendMessageOptions{ShowCaptionAboveMedia: true},
			field: "show_caption_above_media",
			want:  "true",
		},
		{
			name:  "has_spoiler",
			opts:  SendMessageOptions{HasSpoiler: true},
			field: "has_spoiler",
			want:  "true",
		},
		{
			name:  "supports_streaming",
			opts:  SendMessageOptions{SupportsStreaming: true},
			field: "supports_streaming",
			want:  "true",
		},
		{
			name:  "disable_content_type_detection",
			opts:  SendMessageOptions{DisableContentTypeDetection: true},
			field: "disable_content_type_detection",
			want:  "true",
		},
		{
			name:  "is_anonymous",
			opts:  SendMessageOptions{IsAnonymous: true},
			field: "is_anonymous",
			want:  "true",
		},
		{
			name:  "allows_multiple_answers",
			opts:  SendMessageOptions{AllowsMultipleAnswers: true},
			field: "allows_multiple_answers",
			want:  "true",
		},
		{
			name:  "is_closed",
			opts:  SendMessageOptions{IsClosed: true},
			field: "is_closed",
			want:  "true",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := applyOpts(tc.opts.Apply())
			if got := v.Get(tc.field); got != tc.want {
				t.Errorf("%s = %q, want %q", tc.field, got, tc.want)
			}
		})
	}
}

func TestSendMessageOptions_Apply_StringFields(t *testing.T) {
	cases := []struct {
		name  string
		opts  SendMessageOptions
		field string
		want  string
	}{
		{name: "caption", opts: SendMessageOptions{Caption: "my caption"}, field: "caption", want: "my caption"},
		{name: "performer", opts: SendMessageOptions{Performer: "Artist"}, field: "performer", want: "Artist"},
		{name: "title", opts: SendMessageOptions{Title: "Song Name"}, field: "title", want: "Song Name"},
		{name: "thumbnail", opts: SendMessageOptions{Thumbnail: "file-id-thumb"}, field: "thumbnail", want: "file-id-thumb"},
		{name: "business_connection_id", opts: SendMessageOptions{BusinessConnectionID: "biz-123"}, field: "business_connection_id", want: "biz-123"},
		{name: "message_effect_id", opts: SendMessageOptions{MessageEffectID: "effect-abc"}, field: "message_effect_id", want: "effect-abc"},
		{name: "address", opts: SendMessageOptions{Address: "123 Main St"}, field: "address", want: "123 Main St"},
		{name: "foursquare_id", opts: SendMessageOptions{FoursquareID: "fsq-1"}, field: "foursquare_id", want: "fsq-1"},
		{name: "foursquare_type", opts: SendMessageOptions{FoursquareType: "food/restaurant"}, field: "foursquare_type", want: "food/restaurant"},
		{name: "google_place_id", opts: SendMessageOptions{GooglePlaceID: "gp-1"}, field: "google_place_id", want: "gp-1"},
		{name: "google_place_type", opts: SendMessageOptions{GooglePlaceType: "restaurant"}, field: "google_place_type", want: "restaurant"},
		{name: "phone_number", opts: SendMessageOptions{PhoneNumber: "+1234"}, field: "phone_number", want: "+1234"},
		{name: "first_name", opts: SendMessageOptions{FirstName: "Alice"}, field: "first_name", want: "Alice"},
		{name: "last_name", opts: SendMessageOptions{LastName: "Smith"}, field: "last_name", want: "Smith"},
		{name: "vcard", opts: SendMessageOptions{VCard: "BEGIN:VCARD"}, field: "vcard", want: "BEGIN:VCARD"},
		{name: "question", opts: SendMessageOptions{Question: "What?"}, field: "question", want: "What?"},
		{name: "type_poll", opts: SendMessageOptions{Type: "quiz"}, field: "type", want: "quiz"},
		{name: "explanation", opts: SendMessageOptions{Explanation: "Because."}, field: "explanation", want: "Because."},
		{name: "explanation_parse_mode", opts: SendMessageOptions{ExplanationParseMode: "HTML"}, field: "explanation_parse_mode", want: "HTML"},
		{name: "emoji", opts: SendMessageOptions{Emoji: "🎲"}, field: "emoji", want: "🎲"},
		{name: "inline_message_id", opts: SendMessageOptions{InlineMessageID: "inline-42"}, field: "inline_message_id", want: "inline-42"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := applyOpts(tc.opts.Apply())
			if got := v.Get(tc.field); got != tc.want {
				t.Errorf("%s = %q, want %q", tc.field, got, tc.want)
			}
		})
	}
}

func TestSendMessageOptions_Apply_IntFields(t *testing.T) {
	cases := []struct {
		name  string
		opts  SendMessageOptions
		field string
		want  string
	}{
		{name: "message_thread_id", opts: SendMessageOptions{MessageThreadID: 7}, field: "message_thread_id", want: "7"},
		{name: "duration", opts: SendMessageOptions{Duration: 30}, field: "duration", want: "30"},
		{name: "width", opts: SendMessageOptions{Width: 1280}, field: "width", want: "1280"},
		{name: "height", opts: SendMessageOptions{Height: 720}, field: "height", want: "720"},
		{name: "heading", opts: SendMessageOptions{Heading: 180}, field: "heading", want: "180"},
		{name: "proximity_alert_radius", opts: SendMessageOptions{ProximityAlertRadius: 500}, field: "proximity_alert_radius", want: "500"},
		{name: "live_period", opts: SendMessageOptions{LivePeriod: 600}, field: "live_period", want: "600"},
		{name: "correct_option_id", opts: SendMessageOptions{CorrectOptionID: 2}, field: "correct_option_id", want: "2"},
		{name: "open_period", opts: SendMessageOptions{OpenPeriod: 120}, field: "open_period", want: "120"},
		{name: "close_date", opts: SendMessageOptions{CloseDate: 1700000000}, field: "close_date", want: "1700000000"},
		{name: "message_id", opts: SendMessageOptions{MessageID: 99}, field: "message_id", want: "99"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := applyOpts(tc.opts.Apply())
			if got := v.Get(tc.field); got != tc.want {
				t.Errorf("%s = %q, want %q", tc.field, got, tc.want)
			}
		})
	}
}

func TestSendMessageOptions_Apply_ChatID(t *testing.T) {
	cases := []struct {
		name   string
		chatID ChatID
		want   string
	}{
		{name: "numeric", chatID: Int64(-100123), want: "-100123"},
		{name: "username", chatID: Username("@chan"), want: "@chan"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := SendMessageOptions{ChatID: tc.chatID}
			v := applyOpts(opts.Apply())
			if got := v.Get("chat_id"); got != tc.want {
				t.Errorf("chat_id = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSendMessageOptions_Apply_ReplyParameters(t *testing.T) {
	rp := &ReplyParameters{MessageID: 42, AllowSendingWithoutReply: true}
	opts := SendMessageOptions{ReplyParameters: rp}
	v := applyOpts(opts.Apply())
	if got := v.Get("reply_parameters"); got == "" {
		t.Error("reply_parameters: want non-empty JSON, got empty string")
	}
}

func TestSendMessageOptions_Apply_Entities(t *testing.T) {
	ents := []MessageEntity{{Type: "bold", Offset: 0, Length: 5}}
	opts := SendMessageOptions{Entities: ents}
	v := applyOpts(opts.Apply())
	if got := v.Get("entities"); got == "" {
		t.Error("entities: want non-empty JSON, got empty string")
	}
}

func TestSendMessageOptions_Apply_ReplyMarkup(t *testing.T) {
	markup := &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{
		{{Text: "Go", URL: "https://go.dev"}},
	}}
	opts := SendMessageOptions{ReplyMarkup: OptInlineKeyboardMarkup(markup)}
	v := applyOpts(opts.Apply())
	if got := v.Get("reply_markup"); got == "" {
		t.Error("reply_markup: want non-empty JSON, got empty string")
	}
}

func TestSendMessageOptions_Apply_BackwardCompat(t *testing.T) {
	// Ensure OptXxx functions still compile and produce the same output as the
	// struct-field equivalents, confirming backward compatibility.
	v1 := applyOpts([]SendOption{
		OptParseModeHTML,
		OptDisableNotification,
		OptMessageThreadID(3),
	})
	v2 := applyOpts(SendMessageOptions{
		ParseMode:           "HTML",
		DisableNotification: true,
		MessageThreadID:     3,
	}.Apply())
	for _, field := range []string{"parse_mode", "disable_notification", "message_thread_id"} {
		if v1.Get(field) != v2.Get(field) {
			t.Errorf("field %q: OptXxx=%q, struct=%q", field, v1.Get(field), v2.Get(field))
		}
	}
}

func TestSendMessageOptions_Apply_HorizontalAccuracy(t *testing.T) {
	opts := SendMessageOptions{HorizontalAccuracy: 10.5}
	v := applyOpts(opts.Apply())
	if got := v.Get("horizontal_accuracy"); got != "10.5" {
		t.Errorf("horizontal_accuracy = %q, want %q", got, "10.5")
	}
}

func TestSendMessageOptions_Apply_PollOptions(t *testing.T) {
	pollOpts := []InputPollOption{{Text: "Yes"}, {Text: "No"}}
	opts := SendMessageOptions{Options: pollOpts}
	v := applyOpts(opts.Apply())
	if got := v.Get("options"); got == "" {
		t.Error("options: want non-empty JSON, got empty string")
	}
}
