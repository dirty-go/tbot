package tbot

import (
	"context"
	"net/url"
)

// GetFile returns a File describing where to download a file by its id. The
// returned File.FilePath is valid for at least one hour; download it from
// https://api.telegram.org/file/bot<token>/<file_path>.
func (c *Client) GetFile(ctx context.Context, fileID string) (*File, error) {
	req := url.Values{}
	req.Set("file_id", fileID)
	file := &File{}
	err := c.sendRequest(ctx, "/getFile", req, file)
	return file, err
}
