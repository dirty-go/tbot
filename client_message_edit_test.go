package tbot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEditMessageText_ReturnsEditedMessage(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":10,"date":1,"chat":{"id":42,"type":"private"},"text":"edited"}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	msg, err := c.EditMessageText(context.Background(), "edited", OptChatID(Int64(42)), OptMessageID(10))
	if err != nil {
		t.Fatalf("EditMessageText() error: %v", err)
	}
	if msg == nil || msg.Text != "edited" {
		t.Fatalf("unexpected message: %+v", msg)
	}
	for _, want := range []string{"chat_id=42", "message_id=10", "text=edited"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %q", want, body)
		}
	}
}

func TestEditMessageText_InlineReturnsNilMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	msg, err := c.EditMessageText(context.Background(), "edited", OptInlineMessageID("inline-1"))
	if err != nil {
		t.Fatalf("EditMessageText() error: %v", err)
	}
	if msg != nil {
		t.Fatalf("editing an inline message should return nil Message, got %+v", msg)
	}
}

func TestDeleteMessageAndMessages(t *testing.T) {
	var paths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	ctx := context.Background()
	if ok, err := c.DeleteMessage(ctx, Int64(42), 5); err != nil || !ok {
		t.Fatalf("DeleteMessage() = (%v, %v)", ok, err)
	}
	if ok, err := c.DeleteMessages(ctx, Int64(42), []int{5, 6, 7}); err != nil || !ok {
		t.Fatalf("DeleteMessages() = (%v, %v)", ok, err)
	}
	if len(paths) != 2 ||
		!strings.HasSuffix(paths[0], "/deleteMessage") ||
		!strings.HasSuffix(paths[1], "/deleteMessages") {
		t.Fatalf("unexpected paths: %v", paths)
	}
}

func TestSetMessageReaction(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	ok, err := c.SetMessageReaction(context.Background(), Int64(42), 9,
		[]ReactionType{EmojiReaction("👍")}, OptIsBig)
	if err != nil || !ok {
		t.Fatalf("SetMessageReaction() = (%v, %v)", ok, err)
	}
	for _, want := range []string{"chat_id=42", "message_id=9", "reaction=", "is_big=true"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %q", want, body)
		}
	}
}

func TestCopyMessage(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":99}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	id, err := c.CopyMessage(context.Background(), Int64(42), Int64(7), 3)
	if err != nil {
		t.Fatalf("CopyMessage() error: %v", err)
	}
	if id.MessageID != 99 {
		t.Fatalf("MessageID = %d, want 99", id.MessageID)
	}
	for _, want := range []string{"chat_id=42", "from_chat_id=7", "message_id=3"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %q", want, body)
		}
	}
}

func TestStopPoll(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/stopPoll") {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		fmt.Fprint(w, `{"ok":true,"result":{"id":"p1","question":"q","options":[],"total_voter_count":3,"is_closed":true,"is_anonymous":true,"type":"regular","allows_multiple_answers":false}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	poll, err := c.StopPoll(context.Background(), Int64(42), 8)
	if err != nil {
		t.Fatalf("StopPoll() error: %v", err)
	}
	if !poll.IsClosed || poll.ID != "p1" {
		t.Fatalf("unexpected poll: %+v", poll)
	}
}
