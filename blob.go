package blossom

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
)

// Blob is anything that can be read, and has a known size and MIME type.
// This makes it possible to always write a Blob in an HTTP response with the
// correct 'Content-Type' and 'Content-Length' headers.
//
// For typical scenarios, Blobs can be created with already provided constructors
// such as [BlobFromFile], [BlobFromBytes], and [BlobFromStream].
// For more advanced use cases (e.g. encrypted blobs, thread-safe blobs),
// you can implement the [Blob] interface directly.
type Blob interface {
	io.Reader

	// Size return the total size of the blob in bytes.
	Size() int64

	// Type returns the MIME type of the blob.
	Type() string
}

// blob is a generic [Blob] implementation with a reader as the underlying data.
type blob struct {
	io.Reader
	size int64
	typ  string
}

func (b blob) Size() int64  { return b.size }
func (b blob) Type() string { return b.typ }

// BlobFromFile creates a Blob from the given file.
// The file size and content type are detected automatically.
func BlobFromFile(f *os.File) (Blob, error) {
	if f == nil {
		return nil, fmt.Errorf("file is nil")
	}

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	typ, err := DetectType(f)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob from file: %w", err)
	}
	return blob{Reader: f, size: info.Size(), typ: typ}, nil
}

// BlobFromBytes creates a Blob from the given byte slice.
// The size and content type are detected automatically.
func BlobFromBytes(data []byte) Blob {
	return blob{
		Reader: bytes.NewReader(data),
		size:   int64(len(data)),
		typ:    http.DetectContentType(data),
	}
}

// BlobFromStream creates a Blob from a reader with known size and content type.
// It's typically used when dealing with HTTP requests or responses,
// where the size and content type are known from the headers.
func BlobFromStream(r io.Reader, size int64, contentType string) Blob {
	return blob{
		Reader: r,
		size:   size,
		typ:    contentType,
	}
}

// DetectType returns the content type by inspecting up to the first 512 bytes of its data.
// The returned string is suitable for use as a MIME type in HTTP headers (e.g. Content-Type).
// If the type cannot be determined, it returns the default "application/octet-stream" as specified by BUD-01.
func DetectType(r io.ReadSeeker) (string, error) {
	if r == nil {
		return "application/octet-stream", nil
	}

	sniff := make([]byte, 512)
	n, err := io.ReadFull(r, sniff)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", fmt.Errorf("failed to read for MIME sniffing: %w", err)
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to rewind reader after MIME sniffing: %w", err)
	}
	return http.DetectContentType(sniff[:n]), nil
}

// DetectSize returns the total size of the reader in bytes.
func DetectSize(r io.ReadSeeker) (int64, error) {
	if r == nil {
		return 0, nil
	}

	size, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, fmt.Errorf("failed to seek to end of reader: %w", err)
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return -1, fmt.Errorf("failed to rewind reader: %w", err)
	}
	return size, nil
}

// BlobMeta groups metadata of a [Blob].
type BlobMeta struct {
	Hash      Hash
	Type      string // matches [Blob.Type]
	Size      int64
	CreatedAt int64
}

// Extension returns the preferred file extension for the blob's content type.
// The returned extension will begin with a leading dot, as in ".html".
// If no suitable extension is found, it returns ".bin".
func (b BlobMeta) Extension() string {
	return ExtFromType(b.Type)
}

// ExtFromType returns the preferred file extension for the given content type.
// The returned extension will begin with a leading dot, as in ".html".
// If no suitable extension is found, it returns ".bin".
func ExtFromType(contentType string) string {
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || len(exts) == 0 {
		return ".bin"
	}
	// return the last extension, which is the longest, as that's
	// usually the most specific e.g. ".html" instead of ".htm"
	return exts[len(exts)-1]
}
