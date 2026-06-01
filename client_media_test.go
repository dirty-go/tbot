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

const okMessage = `{"ok":true,"result":{"message_id":1,"date":1,"chat":{"id":42,"type":"private"}}}`

func TestSendPhoto_ByReferenceUsesFormBody(t *testing.T) {
	var ct, body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/sendPhoto") {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		ct = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, okMessage)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	msg, err := c.SendPhoto(context.Background(), Int64(42), FileID("photo-id-123"), OptCaption("hi"))
	if err != nil {
		t.Fatalf("SendPhoto() error: %v", err)
	}
	if msg.MessageID != 1 {
		t.Fatalf("unexpected message: %+v", msg)
	}
	if ct != "application/x-www-form-urlencoded" {
		t.Fatalf("content-type = %q, want form-encoded", ct)
	}
	for _, want := range []string{"chat_id=42", "photo=photo-id-123", "caption=hi"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %q", want, body)
		}
	}
}

func TestSendPhoto_UploadUsesMultipart(t *testing.T) {
	var (
		gotChatID  string
		gotName    string
		gotContent string
		gotCT      string
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCT = r.Header.Get("Content-Type")
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("ParseMultipartForm: %v", err)
		}
		gotChatID = r.FormValue("chat_id")
		file, hdr, err := r.FormFile("photo")
		if err != nil {
			t.Fatalf("FormFile(photo): %v", err)
		}
		defer func() { _ = file.Close() }()
		gotName = hdr.Filename
		b, _ := io.ReadAll(file)
		gotContent = string(b)
		fmt.Fprint(w, okMessage)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	_, err := c.SendPhoto(context.Background(), Int64(42), FileBytes("pic.jpg", []byte("RAWPHOTO")))
	if err != nil {
		t.Fatalf("SendPhoto() error: %v", err)
	}
	if !strings.HasPrefix(gotCT, "multipart/form-data") {
		t.Fatalf("content-type = %q, want multipart", gotCT)
	}
	if gotChatID != "42" {
		t.Fatalf("chat_id = %q, want 42", gotChatID)
	}
	if gotName != "pic.jpg" {
		t.Fatalf("filename = %q, want pic.jpg", gotName)
	}
	if gotContent != "RAWPHOTO" {
		t.Fatalf("file content = %q, want RAWPHOTO", gotContent)
	}
}

func TestSendMediaGroup(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/sendMediaGroup") {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":[{"message_id":1,"date":1,"chat":{"id":42,"type":"private"}},{"message_id":2,"date":1,"chat":{"id":42,"type":"private"}}]}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	media := []InputMedia{
		{Type: InputMediaTypePhoto, Media: "https://x/1.jpg"},
		{Type: InputMediaTypePhoto, Media: "https://x/2.jpg"},
	}
	msgs, err := c.SendMediaGroup(context.Background(), Int64(42), media)
	if err != nil {
		t.Fatalf("SendMediaGroup() error: %v", err)
	}
	if len(msgs) != 2 || msgs[1].MessageID != 2 {
		t.Fatalf("unexpected messages: %+v", msgs)
	}
	if !strings.Contains(body, "chat_id=42") || !strings.Contains(body, "media=") {
		t.Fatalf("unexpected body: %q", body)
	}
}
