//go:build integration

package tbot_test

import (
	"context"
	"os"
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
