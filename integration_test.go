//go:build integration

package tbot_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	. "github.com/dirty-go/tbot" // dot-import for test readability is acceptable in _test packages
)

func TestIntegration_GetMe(t *testing.T) {
	token := os.Getenv("TBOT_TEST_TOKEN")
	if token == "" {
		t.Skip("TBOT_TEST_TOKEN not set")
	}
	c := NewClient(token, "")
	me, err := c.Me(context.Background())
	if err != nil {
		t.Fatalf("Me() error: %v", err)
	}
	if me.ID == 0 {
		t.Fatal("expected non-zero bot ID")
	}
	t.Logf("bot: @%s (id=%d)", me.Username, me.ID)
}

func TestIntegration_GetUpdates(t *testing.T) {
	token := os.Getenv("TBOT_TEST_TOKEN")
	if token == "" {
		t.Skip("TBOT_TEST_TOKEN not set")
	}
	c := NewClient(token, "",
		WithRateLimit(0, 0),
		WithMaxRetries(0),
	)
	// limit=1, timeout=0 (short-poll) — may return 0 updates which is fine.
	c2 := NewClient(token, "",
		WithRateLimit(0, 0),
		WithMaxRetries(0),
	)
	_ = c
	updates, err := c2.GetUpdates(context.Background())
	if err != nil {
		t.Fatalf("GetUpdates() error: %v", err)
	}
	// A valid response is a non-nil slice (may be empty).
	if updates == nil {
		t.Fatal("expected non-nil updates slice")
	}
	t.Logf("GetUpdates returned %d updates", len(updates))
}

func TestIntegration_SendMessage(t *testing.T) {
	token := os.Getenv("TBOT_TEST_TOKEN")
	if token == "" {
		t.Skip("TBOT_TEST_TOKEN not set")
	}
	rawChatID := os.Getenv("TBOT_TEST_CHAT_ID")
	if rawChatID == "" {
		t.Skip("TBOT_TEST_CHAT_ID not set")
	}
	chatIDInt, err := strconv.ParseInt(rawChatID, 10, 64)
	if err != nil {
		t.Fatalf("TBOT_TEST_CHAT_ID %q is not a valid int64: %v", rawChatID, err)
	}
	c := NewClient(token, "")
	msg, err := c.SendMessage(context.Background(), Int64(chatIDInt), "integration test ping")
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}
	if msg.MessageID == 0 {
		t.Fatal("expected non-zero message_id")
	}
	t.Logf("sent message_id=%d to chat_id=%d", msg.MessageID, chatIDInt)
}
