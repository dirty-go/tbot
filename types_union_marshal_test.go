package tbot

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestReactionTypeMarshalJSON_WhitelistsVariantFields(t *testing.T) {
	got, err := json.Marshal(ReactionType{
		Type:          ReactionTypeEmoji,
		Emoji:         ":)",
		CustomEmojiID: "should-not-serialize",
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	s := string(got)
	if !strings.Contains(s, `"emoji":":)"`) {
		t.Fatalf("expected emoji field, got %s", s)
	}
	if strings.Contains(s, "custom_emoji_id") {
		t.Fatalf("custom_emoji_id must not be serialized for emoji variant: %s", s)
	}
}

func TestMessageOriginMarshalJSON_WhitelistsVariantFields(t *testing.T) {
	got, err := json.Marshal(MessageOrigin{
		Type:           MessageOriginTypeHiddenUser,
		Date:           1,
		SenderUserName: "hidden",
		Chat:           &Chat{ID: 42},
		MessageID:      99,
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	s := string(got)
	if !strings.Contains(s, `"sender_user_name":"hidden"`) {
		t.Fatalf("expected hidden user field, got %s", s)
	}
	if strings.Contains(s, `"chat"`) || strings.Contains(s, `"message_id"`) {
		t.Fatalf("channel-only fields must not be serialized for hidden_user: %s", s)
	}
}

func TestChatMemberMarshalJSON_WhitelistsVariantFields(t *testing.T) {
	got, err := json.Marshal(ChatMember{
		Status:          ChatMemberStatusMember,
		User:            User{ID: 1, IsBot: false, FirstName: "u"},
		CanManageChat:   true,
		CanSendMessages: true,
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	s := string(got)
	if strings.Contains(s, "can_manage_chat") || strings.Contains(s, "can_send_messages") {
		t.Fatalf("non-member fields must not be serialized for member variant: %s", s)
	}
}

func TestMenuButtonMarshalJSON_WhitelistsVariantFields(t *testing.T) {
	got, err := json.Marshal(MenuButton{
		Type:   MenuButtonTypeCommands,
		Text:   "ignored",
		WebApp: &WebAppInfo{URL: "https://example.com"},
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	s := string(got)
	if strings.Contains(s, `"text"`) || strings.Contains(s, `"web_app"`) {
		t.Fatalf("commands variant must not include web_app payload: %s", s)
	}
}

func TestChatBoostSourceMarshalJSON_WhitelistsVariantFields(t *testing.T) {
	got, err := json.Marshal(ChatBoostSource{
		Source:            ChatBoostSourcePremium,
		User:              &User{ID: 1, IsBot: false, FirstName: "u"},
		GiveawayMessageID: 10,
		PrizeStarCount:    20,
		IsUnclaimed:       true,
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	s := string(got)
	if !strings.Contains(s, `"user"`) {
		t.Fatalf("premium source should include user: %s", s)
	}
	if strings.Contains(s, "giveaway_message_id") || strings.Contains(s, "prize_star_count") || strings.Contains(s, "is_unclaimed") {
		t.Fatalf("giveaway-only fields must not be serialized for premium source: %s", s)
	}
}

func TestBackgroundFillMarshalJSON_WhitelistsVariantFields(t *testing.T) {
	got, err := json.Marshal(BackgroundFill{
		Type:          BackgroundFillSolid,
		Color:         1,
		TopColor:      2,
		BottomColor:   3,
		RotationAngle: 4,
		Colors:        []int{5, 6},
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	s := string(got)
	if !strings.Contains(s, `"color":1`) {
		t.Fatalf("solid fill should include color: %s", s)
	}
	if strings.Contains(s, "top_color") || strings.Contains(s, "colors") {
		t.Fatalf("solid fill must not include gradient/freeform fields: %s", s)
	}
}

func TestBackgroundTypeMarshalJSON_WhitelistsVariantFields(t *testing.T) {
	got, err := json.Marshal(BackgroundType{
		Type:             BackgroundTypeChatTheme,
		ThemeName:        "midnight",
		Document:         &Document{FileID: "doc", FileUniqueID: "uniq"},
		Fill:             &BackgroundFill{Type: BackgroundFillSolid, Color: 1},
		DarkThemeDimming: 50,
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	s := string(got)
	if !strings.Contains(s, `"theme_name":"midnight"`) {
		t.Fatalf("chat_theme should include theme_name: %s", s)
	}
	if strings.Contains(s, `"document"`) || strings.Contains(s, `"fill"`) || strings.Contains(s, `"dark_theme_dimming"`) {
		t.Fatalf("chat_theme must not include other variant fields: %s", s)
	}
}
