package tbot

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClient_GetUpdates_UsesClientFieldsAndAdvancesOffset(t *testing.T) {
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/getUpdates") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, `{"ok":true,"result":[{"update_id":1000},{"update_id":1001}]}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	c.nextOffset = 5
	c.bufferSize = 50
	c.timeout = 30
	c.updateParams = url.Values{}
	c.updateParams.Set("allowed_updates", `["message","callback_query"]`)

	updates, err := c.GetUpdates()
	if err != nil {
		t.Fatalf("GetUpdates() error: %v", err)
	}
	if len(updates) != 2 {
		t.Fatalf("unexpected updates count: %d", len(updates))
	}
	if c.nextOffset != 1002 {
		t.Fatalf("expected nextOffset to advance to 1002, got %d", c.nextOffset)
	}
	for _, want := range []string{"offset=5", "limit=50", "timeout=30", "allowed_updates=%5B%22message%22%2C%22callback_query%22%5D"} {
		if !strings.Contains(gotBody, want) {
			t.Fatalf("expected request body to contain %q, got %q", want, gotBody)
		}
	}
}

func TestClient_SetDeleteGetWebhook(t *testing.T) {
	var setBody, deleteBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/setWebhook"):
			b, _ := io.ReadAll(r.Body)
			setBody = string(b)
			fmt.Fprint(w, `{"ok":true,"result":true}`)
		case strings.HasSuffix(r.URL.Path, "/deleteWebhook"):
			b, _ := io.ReadAll(r.Body)
			deleteBody = string(b)
			fmt.Fprint(w, `{"ok":true,"result":true}`)
		case strings.HasSuffix(r.URL.Path, "/getWebhookInfo"):
			fmt.Fprint(w, `{"ok":true,"result":{"url":"https://example.com/hook","has_custom_certificate":false,"pending_update_count":3}}`)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	ok, err := c.SetWebhook("https://example.com/hook", func(v url.Values) { v.Set("max_connections", "10") })
	if err != nil || !ok {
		t.Fatalf("SetWebhook() = (%v, %v), want (true, nil)", ok, err)
	}
	if !strings.Contains(setBody, "url=https%3A%2F%2Fexample.com%2Fhook") || !strings.Contains(setBody, "max_connections=10") {
		t.Fatalf("unexpected setWebhook body: %q", setBody)
	}

	ok, err = c.DeleteWebhook(true)
	if err != nil || !ok {
		t.Fatalf("DeleteWebhook() = (%v, %v), want (true, nil)", ok, err)
	}
	if !strings.Contains(deleteBody, "drop_pending_updates=true") {
		t.Fatalf("unexpected deleteWebhook body: %q", deleteBody)
	}

	info, err := c.GetWebhookInfo()
	if err != nil {
		t.Fatalf("GetWebhookInfo() error: %v", err)
	}
	if info.URL != "https://example.com/hook" || info.PendingUpdateCount != 3 {
		t.Fatalf("unexpected webhook info: %+v", info)
	}
}

func TestClient_AnswerCallbackQuery(t *testing.T) {
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/answerCallbackQuery") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	ok, err := c.AnswerCallbackQuery("cbq-1", "done", true, func(v url.Values) { v.Set("cache_time", "5") })
	if err != nil || !ok {
		t.Fatalf("AnswerCallbackQuery() = (%v, %v), want (true, nil)", ok, err)
	}
	for _, want := range []string{"callback_query_id=cbq-1", "text=done", "show_alert=true", "cache_time=5"} {
		if !strings.Contains(gotBody, want) {
			t.Fatalf("expected callback body to contain %q, got %q", want, gotBody)
		}
	}
}

func TestClient_SendLocationSendVenueSendContactSendPollSendDice(t *testing.T) {
	type captured struct {
		path string
		body string
	}
	var seen []captured
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		seen = append(seen, captured{path: r.URL.Path, body: string(b)})
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"date":0,"chat":{"id":42,"type":"private"}}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	if _, err := c.SendLocation("42", 1.23, 4.56); err != nil {
		t.Fatalf("SendLocation() error: %v", err)
	}
	if _, err := c.SendVenue("42", 1.23, 4.56, "title", "address"); err != nil {
		t.Fatalf("SendVenue() error: %v", err)
	}
	if _, err := c.SendContact("42", "+123", "Alice"); err != nil {
		t.Fatalf("SendContact() error: %v", err)
	}
	if _, err := c.SendPoll("42", "q?", []InputPollOption{{Text: "A"}, {Text: "B"}}); err != nil {
		t.Fatalf("SendPoll() error: %v", err)
	}
	if _, err := c.SendDice("42", func(v url.Values) { v.Set("emoji", "dice") }); err != nil {
		t.Fatalf("SendDice() error: %v", err)
	}

	if len(seen) != 5 {
		t.Fatalf("expected 5 requests, got %d", len(seen))
	}
	checks := []struct {
		idx        int
		pathSuffix string
		fields     []string
	}{
		{0, "/sendLocation", []string{"chat_id=42", "latitude=1.23", "longitude=4.56"}},
		{1, "/sendVenue", []string{"chat_id=42", "title=title", "address=address"}},
		{2, "/sendContact", []string{"chat_id=42", "phone_number=%2B123", "first_name=Alice"}},
		{3, "/sendPoll", []string{"chat_id=42", "question=q%3F", "options=%5B%7B%22text%22%3A%22A%22%7D%2C%7B%22text%22%3A%22B%22%7D%5D"}},
		{4, "/sendDice", []string{"chat_id=42", "emoji=dice"}},
	}
	for _, tc := range checks {
		got := seen[tc.idx]
		if !strings.HasSuffix(got.path, tc.pathSuffix) {
			t.Fatalf("request %d path = %q, want suffix %q", tc.idx, got.path, tc.pathSuffix)
		}
		for _, field := range tc.fields {
			if !strings.Contains(got.body, field) {
				t.Fatalf("request %d expected body to contain %q, got %q", tc.idx, field, got.body)
			}
		}
	}
}
