package tbot

import (
	"encoding/json"
	"testing"
)

// roundTrip marshals v then unmarshals into a zero value of the same type.
func roundTripJSON[T any](t *testing.T, v T) T {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out T
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func TestReactionTypeRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		in    ReactionType
		check func(*testing.T, ReactionType)
	}{
		{
			name: "emoji",
			in:   ReactionType{Type: ReactionTypeEmoji, Emoji: "👍"},
			check: func(t *testing.T, got ReactionType) {
				if got.Type != ReactionTypeEmoji {
					t.Errorf("Type: got %q want %q", got.Type, ReactionTypeEmoji)
				}
				if got.Emoji != "👍" {
					t.Errorf("Emoji: got %q want %q", got.Emoji, "👍")
				}
			},
		},
		{
			name: "custom_emoji",
			in:   ReactionType{Type: ReactionTypeCustomEmoji, CustomEmojiID: "custom-123"},
			check: func(t *testing.T, got ReactionType) {
				if got.Type != ReactionTypeCustomEmoji {
					t.Errorf("Type: got %q want %q", got.Type, ReactionTypeCustomEmoji)
				}
				if got.CustomEmojiID != "custom-123" {
					t.Errorf("CustomEmojiID: got %q want %q", got.CustomEmojiID, "custom-123")
				}
			},
		},
		{
			name: "paid",
			in:   ReactionType{Type: ReactionTypePaid},
			check: func(t *testing.T, got ReactionType) {
				if got.Type != ReactionTypePaid {
					t.Errorf("Type: got %q want %q", got.Type, ReactionTypePaid)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestMessageOriginRoundTrip(t *testing.T) {
	user := &User{ID: 42, IsBot: false, FirstName: "Alice"}
	chat := &Chat{ID: 99, Type: "channel"}
	cases := []struct {
		name  string
		in    MessageOrigin
		check func(*testing.T, MessageOrigin)
	}{
		{
			name: "user",
			in:   MessageOrigin{Type: MessageOriginTypeUser, Date: 1000, SenderUser: user},
			check: func(t *testing.T, got MessageOrigin) {
				if got.Type != MessageOriginTypeUser {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.SenderUser == nil || got.SenderUser.ID != 42 {
					t.Errorf("SenderUser.ID: got %v", got.SenderUser)
				}
			},
		},
		{
			name: "hidden_user",
			in:   MessageOrigin{Type: MessageOriginTypeHiddenUser, Date: 1000, SenderUserName: "hidden"},
			check: func(t *testing.T, got MessageOrigin) {
				if got.Type != MessageOriginTypeHiddenUser {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.SenderUserName != "hidden" {
					t.Errorf("SenderUserName: got %q", got.SenderUserName)
				}
			},
		},
		{
			name: "chat",
			in:   MessageOrigin{Type: MessageOriginTypeChat, Date: 1000, SenderChat: chat, AuthorSignature: "sig"},
			check: func(t *testing.T, got MessageOrigin) {
				if got.Type != MessageOriginTypeChat {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.SenderChat == nil || got.SenderChat.ID != 99 {
					t.Errorf("SenderChat.ID: got %v", got.SenderChat)
				}
				if got.AuthorSignature != "sig" {
					t.Errorf("AuthorSignature: got %q", got.AuthorSignature)
				}
			},
		},
		{
			name: "channel",
			in:   MessageOrigin{Type: MessageOriginTypeChannel, Date: 1000, Chat: chat, MessageID: 7, AuthorSignature: "chan"},
			check: func(t *testing.T, got MessageOrigin) {
				if got.Type != MessageOriginTypeChannel {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Chat == nil || got.Chat.ID != 99 {
					t.Errorf("Chat.ID: got %v", got.Chat)
				}
				if got.MessageID != 7 {
					t.Errorf("MessageID: got %d", got.MessageID)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestChatMemberRoundTrip(t *testing.T) {
	u := User{ID: 1, IsBot: false, FirstName: "Bob"}
	cases := []struct {
		name  string
		in    ChatMember
		check func(*testing.T, ChatMember)
	}{
		{
			name: "creator",
			in:   ChatMember{Status: ChatMemberStatusCreator, User: u, IsAnonymous: true, CustomTitle: "Boss"},
			check: func(t *testing.T, got ChatMember) {
				if got.Status != ChatMemberStatusCreator {
					t.Errorf("Status: got %q", got.Status)
				}
				if !got.IsAnonymous {
					t.Error("IsAnonymous: want true")
				}
				if got.CustomTitle != "Boss" {
					t.Errorf("CustomTitle: got %q", got.CustomTitle)
				}
			},
		},
		{
			name: "administrator",
			in:   ChatMember{Status: ChatMemberStatusAdministrator, User: u, CanManageChat: true, CanDeleteMessages: true},
			check: func(t *testing.T, got ChatMember) {
				if got.Status != ChatMemberStatusAdministrator {
					t.Errorf("Status: got %q", got.Status)
				}
				if !got.CanManageChat {
					t.Error("CanManageChat: want true")
				}
				if !got.CanDeleteMessages {
					t.Error("CanDeleteMessages: want true")
				}
			},
		},
		{
			name: "restricted",
			in:   ChatMember{Status: ChatMemberStatusRestricted, User: u, IsMember: true, CanSendMessages: true, UntilDate: 9999},
			check: func(t *testing.T, got ChatMember) {
				if got.Status != ChatMemberStatusRestricted {
					t.Errorf("Status: got %q", got.Status)
				}
				if !got.IsMember {
					t.Error("IsMember: want true")
				}
				if !got.CanSendMessages {
					t.Error("CanSendMessages: want true")
				}
				if got.UntilDate != 9999 {
					t.Errorf("UntilDate: got %d", got.UntilDate)
				}
			},
		},
		{
			name: "member",
			in:   ChatMember{Status: ChatMemberStatusMember, User: u},
			check: func(t *testing.T, got ChatMember) {
				if got.Status != ChatMemberStatusMember {
					t.Errorf("Status: got %q", got.Status)
				}
				if got.User.ID != 1 {
					t.Errorf("User.ID: got %d", got.User.ID)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestMenuButtonRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		in    MenuButton
		check func(*testing.T, MenuButton)
	}{
		{
			name: "default",
			in:   MenuButton{Type: MenuButtonTypeDefault},
			check: func(t *testing.T, got MenuButton) {
				if got.Type != MenuButtonTypeDefault {
					t.Errorf("Type: got %q", got.Type)
				}
			},
		},
		{
			name: "commands",
			in:   MenuButton{Type: MenuButtonTypeCommands},
			check: func(t *testing.T, got MenuButton) {
				if got.Type != MenuButtonTypeCommands {
					t.Errorf("Type: got %q", got.Type)
				}
			},
		},
		{
			name: "web_app",
			in:   MenuButton{Type: MenuButtonTypeWebApp, Text: "Open", WebApp: &WebAppInfo{URL: "https://example.com"}},
			check: func(t *testing.T, got MenuButton) {
				if got.Type != MenuButtonTypeWebApp {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Text != "Open" {
					t.Errorf("Text: got %q", got.Text)
				}
				if got.WebApp == nil || got.WebApp.URL != "https://example.com" {
					t.Errorf("WebApp.URL: got %v", got.WebApp)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestChatBoostSourceRoundTrip(t *testing.T) {
	u := &User{ID: 5, IsBot: false, FirstName: "Carol"}
	cases := []struct {
		name  string
		in    ChatBoostSource
		check func(*testing.T, ChatBoostSource)
	}{
		{
			name: "premium",
			in:   ChatBoostSource{Source: ChatBoostSourcePremium, User: u},
			check: func(t *testing.T, got ChatBoostSource) {
				if got.Source != ChatBoostSourcePremium {
					t.Errorf("Source: got %q", got.Source)
				}
				if got.User == nil || got.User.ID != 5 {
					t.Errorf("User.ID: got %v", got.User)
				}
			},
		},
		{
			name: "gift_code",
			in:   ChatBoostSource{Source: ChatBoostSourceGiftCode, User: u},
			check: func(t *testing.T, got ChatBoostSource) {
				if got.Source != ChatBoostSourceGiftCode {
					t.Errorf("Source: got %q", got.Source)
				}
				if got.User == nil || got.User.ID != 5 {
					t.Errorf("User.ID: got %v", got.User)
				}
			},
		},
		{
			name: "giveaway",
			in:   ChatBoostSource{Source: ChatBoostSourceGiveaway, GiveawayMessageID: 77, PrizeStarCount: 10, IsUnclaimed: true},
			check: func(t *testing.T, got ChatBoostSource) {
				if got.Source != ChatBoostSourceGiveaway {
					t.Errorf("Source: got %q", got.Source)
				}
				if got.GiveawayMessageID != 77 {
					t.Errorf("GiveawayMessageID: got %d", got.GiveawayMessageID)
				}
				if got.PrizeStarCount != 10 {
					t.Errorf("PrizeStarCount: got %d", got.PrizeStarCount)
				}
				if !got.IsUnclaimed {
					t.Error("IsUnclaimed: want true")
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestBackgroundFillRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		in    BackgroundFill
		check func(*testing.T, BackgroundFill)
	}{
		{
			name: "solid",
			in:   BackgroundFill{Type: BackgroundFillSolid, Color: 0xFF0000},
			check: func(t *testing.T, got BackgroundFill) {
				if got.Type != BackgroundFillSolid {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Color != 0xFF0000 {
					t.Errorf("Color: got %d", got.Color)
				}
			},
		},
		{
			name: "gradient",
			in:   BackgroundFill{Type: BackgroundFillGradient, TopColor: 1, BottomColor: 2, RotationAngle: 45},
			check: func(t *testing.T, got BackgroundFill) {
				if got.Type != BackgroundFillGradient {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.TopColor != 1 || got.BottomColor != 2 {
					t.Errorf("TopColor/BottomColor: got %d/%d", got.TopColor, got.BottomColor)
				}
				if got.RotationAngle != 45 {
					t.Errorf("RotationAngle: got %d", got.RotationAngle)
				}
			},
		},
		{
			name: "freeform_gradient",
			in:   BackgroundFill{Type: BackgroundFillFreeformGradient, Colors: []int{1, 2, 3}},
			check: func(t *testing.T, got BackgroundFill) {
				if got.Type != BackgroundFillFreeformGradient {
					t.Errorf("Type: got %q", got.Type)
				}
				if len(got.Colors) != 3 || got.Colors[0] != 1 {
					t.Errorf("Colors: got %v", got.Colors)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestChatMemberLeftAndBannedRoundTrip(t *testing.T) {
	u := User{ID: 9, IsBot: false, FirstName: "Dave"}
	cases := []struct {
		name  string
		in    ChatMember
		check func(*testing.T, ChatMember)
	}{
		{
			name: "left",
			in:   ChatMember{Status: ChatMemberStatusLeft, User: u},
			check: func(t *testing.T, got ChatMember) {
				if got.Status != ChatMemberStatusLeft {
					t.Errorf("Status: got %q", got.Status)
				}
				if got.User.ID != 9 {
					t.Errorf("User.ID: got %d", got.User.ID)
				}
			},
		},
		{
			name: "kicked",
			in:   ChatMember{Status: ChatMemberStatusBanned, User: u},
			check: func(t *testing.T, got ChatMember) {
				if got.Status != ChatMemberStatusBanned {
					t.Errorf("Status: got %q", got.Status)
				}
				if got.User.ID != 9 {
					t.Errorf("User.ID: got %d", got.User.ID)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestMaybeInaccessibleMessageRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		in    MaybeInaccessibleMessage
		check func(*testing.T, MaybeInaccessibleMessage)
	}{
		{
			name: "accessible",
			in: MaybeInaccessibleMessage{Message: Message{
				MessageID: 5,
				Date:      1700000000,
				Chat:      Chat{ID: 42, Type: "private"},
				Text:      "hello",
			}},
			check: func(t *testing.T, got MaybeInaccessibleMessage) {
				if got.IsInaccessible() {
					t.Error("IsInaccessible: want false for non-zero Date")
				}
				if got.MessageID != 5 {
					t.Errorf("MessageID: got %d", got.MessageID)
				}
				if got.Text != "hello" {
					t.Errorf("Text: got %q", got.Text)
				}
			},
		},
		{
			name: "inaccessible",
			in: MaybeInaccessibleMessage{Message: Message{
				MessageID: 7,
				Date:      0,
				Chat:      Chat{ID: 99, Type: "group"},
			}},
			check: func(t *testing.T, got MaybeInaccessibleMessage) {
				if !got.IsInaccessible() {
					t.Error("IsInaccessible: want true for zero Date")
				}
				if got.MessageID != 7 {
					t.Errorf("MessageID: got %d", got.MessageID)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestBotCommandScopeRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		in    BotCommandScope
		check func(*testing.T, BotCommandScope)
	}{
		{
			name: "default",
			in:   BotCommandScope{Type: BotCommandScopeTypeDefault},
			check: func(t *testing.T, got BotCommandScope) {
				if got.Type != BotCommandScopeTypeDefault {
					t.Errorf("Type: got %q", got.Type)
				}
			},
		},
		{
			name: "all_private_chats",
			in:   BotCommandScope{Type: BotCommandScopeTypeAllPrivateChats},
			check: func(t *testing.T, got BotCommandScope) {
				if got.Type != BotCommandScopeTypeAllPrivateChats {
					t.Errorf("Type: got %q", got.Type)
				}
			},
		},
		{
			name: "all_group_chats",
			in:   BotCommandScope{Type: BotCommandScopeTypeAllGroupChats},
			check: func(t *testing.T, got BotCommandScope) {
				if got.Type != BotCommandScopeTypeAllGroupChats {
					t.Errorf("Type: got %q", got.Type)
				}
			},
		},
		{
			name: "all_chat_administrators",
			in:   BotCommandScope{Type: BotCommandScopeTypeAllChatAdministrators},
			check: func(t *testing.T, got BotCommandScope) {
				if got.Type != BotCommandScopeTypeAllChatAdministrators {
					t.Errorf("Type: got %q", got.Type)
				}
			},
		},
		{
			name: "chat",
			in:   BotCommandScope{Type: BotCommandScopeTypeChat, ChatID: int64(123)},
			check: func(t *testing.T, got BotCommandScope) {
				if got.Type != BotCommandScopeTypeChat {
					t.Errorf("Type: got %q", got.Type)
				}
			},
		},
		{
			name: "chat_administrators",
			in:   BotCommandScope{Type: BotCommandScopeTypeChatAdministrators, ChatID: "@mychan"},
			check: func(t *testing.T, got BotCommandScope) {
				if got.Type != BotCommandScopeTypeChatAdministrators {
					t.Errorf("Type: got %q", got.Type)
				}
			},
		},
		{
			name: "chat_member",
			in:   BotCommandScope{Type: BotCommandScopeTypeChatMember, ChatID: int64(456), UserID: 789},
			check: func(t *testing.T, got BotCommandScope) {
				if got.Type != BotCommandScopeTypeChatMember {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.UserID != 789 {
					t.Errorf("UserID: got %d", got.UserID)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestInputMediaRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		in    InputMedia
		check func(*testing.T, InputMedia)
	}{
		{
			name: "photo",
			in:   InputMedia{Type: InputMediaTypePhoto, Media: "https://x/photo.jpg", Caption: "a photo"},
			check: func(t *testing.T, got InputMedia) {
				if got.Type != InputMediaTypePhoto {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Media != "https://x/photo.jpg" {
					t.Errorf("Media: got %q", got.Media)
				}
				if got.Caption != "a photo" {
					t.Errorf("Caption: got %q", got.Caption)
				}
			},
		},
		{
			name: "video",
			in:   InputMedia{Type: InputMediaTypeVideo, Media: "file-id-vid", Duration: 30, Width: 1920, Height: 1080, SupportsStreaming: true},
			check: func(t *testing.T, got InputMedia) {
				if got.Type != InputMediaTypeVideo {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Duration != 30 {
					t.Errorf("Duration: got %d", got.Duration)
				}
				if !got.SupportsStreaming {
					t.Error("SupportsStreaming: want true")
				}
			},
		},
		{
			name: "animation",
			in:   InputMedia{Type: InputMediaTypeAnimation, Media: "file-id-anim", HasSpoiler: true},
			check: func(t *testing.T, got InputMedia) {
				if got.Type != InputMediaTypeAnimation {
					t.Errorf("Type: got %q", got.Type)
				}
				if !got.HasSpoiler {
					t.Error("HasSpoiler: want true")
				}
			},
		},
		{
			name: "audio",
			in:   InputMedia{Type: InputMediaTypeAudio, Media: "file-id-aud", Performer: "The Band", Title: "Song"},
			check: func(t *testing.T, got InputMedia) {
				if got.Type != InputMediaTypeAudio {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Performer != "The Band" {
					t.Errorf("Performer: got %q", got.Performer)
				}
				if got.Title != "Song" {
					t.Errorf("Title: got %q", got.Title)
				}
			},
		},
		{
			name: "document",
			in:   InputMedia{Type: InputMediaTypeDocument, Media: "file-id-doc", DisableContentTypeDetection: true},
			check: func(t *testing.T, got InputMedia) {
				if got.Type != InputMediaTypeDocument {
					t.Errorf("Type: got %q", got.Type)
				}
				if !got.DisableContentTypeDetection {
					t.Error("DisableContentTypeDetection: want true")
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestInputPaidMediaRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		in    InputPaidMedia
		check func(*testing.T, InputPaidMedia)
	}{
		{
			name: "photo",
			in:   InputPaidMedia{Type: "photo", Media: "file-id-photo"},
			check: func(t *testing.T, got InputPaidMedia) {
				if got.Type != "photo" {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Media != "file-id-photo" {
					t.Errorf("Media: got %q", got.Media)
				}
			},
		},
		{
			name: "video",
			in:   InputPaidMedia{Type: "video", Media: "file-id-vid", Width: 1280, Height: 720, Duration: 60, SupportsStreaming: true},
			check: func(t *testing.T, got InputPaidMedia) {
				if got.Type != "video" {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Width != 1280 || got.Height != 720 {
					t.Errorf("Width/Height: got %d/%d", got.Width, got.Height)
				}
				if !got.SupportsStreaming {
					t.Error("SupportsStreaming: want true")
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}

func TestBackgroundTypeRoundTrip(t *testing.T) {
	fill := &BackgroundFill{Type: BackgroundFillSolid, Color: 1}
	doc := &Document{FileID: "fid", FileUniqueID: "uid"}
	cases := []struct {
		name  string
		in    BackgroundType
		check func(*testing.T, BackgroundType)
	}{
		{
			name: "fill",
			in:   BackgroundType{Type: BackgroundTypeFill, Fill: fill, DarkThemeDimming: 30},
			check: func(t *testing.T, got BackgroundType) {
				if got.Type != BackgroundTypeFill {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Fill == nil || got.Fill.Type != BackgroundFillSolid {
					t.Errorf("Fill: got %v", got.Fill)
				}
				if got.DarkThemeDimming != 30 {
					t.Errorf("DarkThemeDimming: got %d", got.DarkThemeDimming)
				}
			},
		},
		{
			name: "wallpaper",
			in:   BackgroundType{Type: BackgroundTypeWallpaper, Document: doc, DarkThemeDimming: 50, IsBlurred: true, IsMoving: true},
			check: func(t *testing.T, got BackgroundType) {
				if got.Type != BackgroundTypeWallpaper {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Document == nil || got.Document.FileID != "fid" {
					t.Errorf("Document: got %v", got.Document)
				}
				if !got.IsBlurred || !got.IsMoving {
					t.Errorf("IsBlurred/IsMoving: got %v/%v", got.IsBlurred, got.IsMoving)
				}
			},
		},
		{
			name: "pattern",
			in:   BackgroundType{Type: BackgroundTypePattern, Document: doc, Fill: fill, Intensity: 80, IsInverted: true},
			check: func(t *testing.T, got BackgroundType) {
				if got.Type != BackgroundTypePattern {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.Intensity != 80 {
					t.Errorf("Intensity: got %d", got.Intensity)
				}
				if !got.IsInverted {
					t.Error("IsInverted: want true")
				}
			},
		},
		{
			name: "chat_theme",
			in:   BackgroundType{Type: BackgroundTypeChatTheme, ThemeName: "midnight"},
			check: func(t *testing.T, got BackgroundType) {
				if got.Type != BackgroundTypeChatTheme {
					t.Errorf("Type: got %q", got.Type)
				}
				if got.ThemeName != "midnight" {
					t.Errorf("ThemeName: got %q", got.ThemeName)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, roundTripJSON(t, tc.in))
		})
	}
}
