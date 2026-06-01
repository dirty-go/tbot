package tbot

import (
	"context"
	"net/url"
)

// SendPhoto sends a photo. The photo may be a file_id, URL, or fresh upload
// (see InputFile). Set a caption with OptCaption and OptParseModeHTML.
func (c *Client) SendPhoto(ctx context.Context, chatID ChatID, photo InputFile, opts ...SendOption) (*Message, error) {
	return c.sendMedia(ctx, "/sendPhoto", chatID, "photo", photo, opts)
}

// SendAudio sends an audio file to be displayed as music. Set metadata with
// OptPerformer, OptTitle, and OptDuration.
func (c *Client) SendAudio(ctx context.Context, chatID ChatID, audio InputFile, opts ...SendOption) (*Message, error) {
	return c.sendMedia(ctx, "/sendAudio", chatID, "audio", audio, opts)
}

// SendDocument sends a general file. Disable content-type sniffing with
// OptDisableContentTypeDetection.
func (c *Client) SendDocument(ctx context.Context, chatID ChatID, document InputFile, opts ...SendOption) (*Message, error) {
	return c.sendMedia(ctx, "/sendDocument", chatID, "document", document, opts)
}

// SendVideo sends a video (MPEG-4). Mark streamable uploads with
// OptSupportsStreaming and set dimensions with OptWidth / OptHeight.
func (c *Client) SendVideo(ctx context.Context, chatID ChatID, video InputFile, opts ...SendOption) (*Message, error) {
	return c.sendMedia(ctx, "/sendVideo", chatID, "video", video, opts)
}

// SendAnimation sends an animation (GIF or soundless H.264/MPEG-4).
func (c *Client) SendAnimation(ctx context.Context, chatID ChatID, animation InputFile, opts ...SendOption) (*Message, error) {
	return c.sendMedia(ctx, "/sendAnimation", chatID, "animation", animation, opts)
}

// SendVoice sends a voice note (OGG/OPUS, or MP3/M4A as a file).
func (c *Client) SendVoice(ctx context.Context, chatID ChatID, voice InputFile, opts ...SendOption) (*Message, error) {
	return c.sendMedia(ctx, "/sendVoice", chatID, "voice", voice, opts)
}

// SendVideoNote sends a rounded square video message (up to 1 minute).
func (c *Client) SendVideoNote(ctx context.Context, chatID ChatID, videoNote InputFile, opts ...SendOption) (*Message, error) {
	return c.sendMedia(ctx, "/sendVideoNote", chatID, "video_note", videoNote, opts)
}

// sendMedia is the shared body for the single-file send* methods.
func (c *Client) sendMedia(ctx context.Context, method string, chatID ChatID, field string, file InputFile, opts []SendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	var files []fileField
	addFile(req, &files, field, file)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.send(ctx, method, req, files, msg)
	return msg, err
}

// SendMediaGroup sends a group of photos, videos, documents, or audios as an
// album and returns the sent messages. Each item's Media must reference an
// existing file by file_id or URL; uploading album members in a single call is
// not yet supported.
func (c *Client) SendMediaGroup(ctx context.Context, chatID ChatID, media []InputMedia, opts ...SendOption) ([]Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID.String())
	req.Set("media", structString(media))
	for _, opt := range opts {
		opt(req)
	}
	var msgs []Message
	err := c.sendRequest(ctx, "/sendMediaGroup", req, &msgs)
	return msgs, err
}
