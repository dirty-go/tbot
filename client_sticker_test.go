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

func TestGetStickerSet_GoldenPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/getStickerSet") {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		fmt.Fprint(w, `{"ok":true,"result":{"name":"test","title":"Test Set","sticker_type":"regular","stickers":[]}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	set, err := c.GetStickerSet(context.Background(), "test")
	if err != nil {
		t.Fatalf("GetStickerSet() error: %v", err)
	}
	if set.Name != "test" || set.Title != "Test Set" || set.StickerType != StickerTypeRegular {
		t.Fatalf("unexpected StickerSet: %+v", set)
	}
	if set.Stickers == nil {
		t.Fatal("expected non-nil Stickers slice")
	}
}

func TestGetCustomEmojiStickers_RequestBody(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":[]}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	stickers, err := c.GetCustomEmojiStickers(context.Background(), []string{"emoji1", "emoji2"})
	if err != nil {
		t.Fatalf("GetCustomEmojiStickers() error: %v", err)
	}
	if stickers == nil {
		t.Fatal("expected non-nil stickers slice")
	}
	if !strings.Contains(body, "custom_emoji_ids=") {
		t.Fatalf("request body missing custom_emoji_ids, got: %q", body)
	}
	// The JSON array should contain both IDs.
	for _, id := range []string{"emoji1", "emoji2"} {
		if !strings.Contains(body, id) {
			t.Fatalf("request body missing emoji id %q, got: %q", id, body)
		}
	}
}

func TestUploadStickerFile_MultipartBody(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":{"file_id":"f1","file_unique_id":"fu1","file_size":42}}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	f, err := c.UploadStickerFile(
		context.Background(),
		123,
		FileBytes("sticker.webp", []byte("data")),
		StickerFormatStatic,
	)
	if err != nil {
		t.Fatalf("UploadStickerFile() error: %v", err)
	}
	if f.FileID != "f1" {
		t.Fatalf("unexpected FileID: %q", f.FileID)
	}
	if !strings.Contains(body, "sticker") {
		t.Fatalf("multipart body missing sticker field, got: %q", body)
	}
}

func TestDeleteStickerFromSet_RequestBody(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	ok, err := c.DeleteStickerFromSet(context.Background(), "sticker_file_id_1")
	if err != nil || !ok {
		t.Fatalf("DeleteStickerFromSet() = (%v, %v)", ok, err)
	}
	if !strings.Contains(body, "sticker=") {
		t.Fatalf("request body missing sticker param, got: %q", body)
	}
	if !strings.Contains(body, "sticker_file_id_1") {
		t.Fatalf("request body missing sticker value, got: %q", body)
	}
}
