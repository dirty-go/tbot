package tbot

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// InputFile is a file to send to the Bot API. The API accepts three shapes:
//
//   - a previously-uploaded file referenced by its file_id (FileID),
//   - an HTTPS URL Telegram fetches itself (FileURL),
//   - a fresh upload streamed in a multipart/form-data body (FilePath,
//     FileReader, FileBytes).
//
// The request layer inspects the value and chooses form-encoding for reference
// files or a multipart body for uploads, so callers never assemble the
// multipart parts themselves.
//
// The interface is sealed: only constructors in this package implement it, so
// the set of shapes is closed and exhaustive.
type InputFile interface {
	// value is the wire reference (file_id or URL) for reference files, or ""
	// when the file must be uploaded.
	value() string
	// needsUpload reports whether the file must be sent as multipart.
	needsUpload() bool
	// open returns the upload contents and the filename for the multipart
	// part. The caller owns the returned ReadCloser and must close it. It is
	// only called when needsUpload reports true.
	open() (io.ReadCloser, string, error)
}

// FileID references a file already stored on Telegram's servers.
func FileID(id string) InputFile { return refFile{ref: id} }

// FileURL references a file by HTTPS URL for Telegram to fetch.
func FileURL(url string) InputFile { return refFile{ref: url} }

// FilePath uploads a file read from the local filesystem. The base name of
// path is used as the multipart filename.
func FilePath(path string) InputFile { return pathFile{path: path} }

// FileReader uploads a file streamed from r, labelled with name. If r also
// implements io.Closer it is closed after the body is assembled.
func FileReader(name string, r io.Reader) InputFile {
	return readerFile{name: name, reader: r}
}

// FileBytes uploads an in-memory file labelled with name.
func FileBytes(name string, data []byte) InputFile {
	return bytesFile{name: name, data: data}
}

// refFile is a file_id or URL passed by reference, never uploaded.
type refFile struct{ ref string }

func (f refFile) value() string                        { return f.ref }
func (f refFile) needsUpload() bool                    { return false }
func (f refFile) open() (io.ReadCloser, string, error) { return nil, "", nil }

// pathFile uploads a file opened lazily from disk.
type pathFile struct{ path string }

func (f pathFile) value() string     { return "" }
func (f pathFile) needsUpload() bool { return true }
func (f pathFile) open() (io.ReadCloser, string, error) {
	// The path is supplied deliberately by the caller to choose what to
	// upload; opening it is the documented purpose of FilePath.
	file, err := os.Open(f.path) //nolint:gosec // caller-chosen upload path
	if err != nil {
		return nil, "", err
	}
	return file, filepath.Base(f.path), nil
}

// readerFile uploads a file streamed from an arbitrary reader.
type readerFile struct {
	name   string
	reader io.Reader
}

func (f readerFile) value() string     { return "" }
func (f readerFile) needsUpload() bool { return true }
func (f readerFile) open() (io.ReadCloser, string, error) {
	if rc, ok := f.reader.(io.ReadCloser); ok {
		return rc, f.name, nil
	}
	return io.NopCloser(f.reader), f.name, nil
}

// bytesFile uploads an in-memory buffer.
type bytesFile struct {
	name string
	data []byte
}

func (f bytesFile) value() string     { return "" }
func (f bytesFile) needsUpload() bool { return true }
func (f bytesFile) open() (io.ReadCloser, string, error) {
	return io.NopCloser(bytes.NewReader(f.data)), f.name, nil
}
