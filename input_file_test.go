package tbot

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInputFile_References(t *testing.T) {
	for _, f := range []InputFile{FileID("abc"), FileURL("https://x/y.jpg")} {
		if f.needsUpload() {
			t.Errorf("%T should not need upload", f)
		}
	}
	if got := FileID("abc").value(); got != "abc" {
		t.Errorf("FileID value = %q", got)
	}
	if got := FileURL("https://x/y.jpg").value(); got != "https://x/y.jpg" {
		t.Errorf("FileURL value = %q", got)
	}
}

func TestInputFile_BytesUpload(t *testing.T) {
	f := FileBytes("note.txt", []byte("hello"))
	if !f.needsUpload() {
		t.Fatal("FileBytes should need upload")
	}
	rc, name, err := f.open()
	if err != nil {
		t.Fatalf("open() error: %v", err)
	}
	defer func() { _ = rc.Close() }()
	if name != "note.txt" {
		t.Errorf("name = %q, want note.txt", name)
	}
	b, _ := io.ReadAll(rc)
	if string(b) != "hello" {
		t.Errorf("contents = %q, want hello", b)
	}
}

func TestInputFile_ReaderUpload(t *testing.T) {
	f := FileReader("r.bin", strings.NewReader("data"))
	rc, name, err := f.open()
	if err != nil {
		t.Fatalf("open() error: %v", err)
	}
	defer func() { _ = rc.Close() }()
	if name != "r.bin" {
		t.Errorf("name = %q", name)
	}
	b, _ := io.ReadAll(rc)
	if string(b) != "data" {
		t.Errorf("contents = %q", b)
	}
}

func TestInputFile_PathUpload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "photo.jpg")
	if err := os.WriteFile(path, []byte("JPEGBYTES"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f := FilePath(path)
	if !f.needsUpload() {
		t.Fatal("FilePath should need upload")
	}
	rc, name, err := f.open()
	if err != nil {
		t.Fatalf("open() error: %v", err)
	}
	defer func() { _ = rc.Close() }()
	if name != "photo.jpg" {
		t.Errorf("name = %q, want base name photo.jpg", name)
	}
	b, _ := io.ReadAll(rc)
	if string(b) != "JPEGBYTES" {
		t.Errorf("contents = %q", b)
	}
}

func TestInputFile_PathOpenError(t *testing.T) {
	_, _, err := FilePath(filepath.Join(t.TempDir(), "missing.bin")).open()
	if err == nil {
		t.Fatal("expected error opening a missing file")
	}
}
