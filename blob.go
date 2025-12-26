package blossom

import (
	"fmt"
	"io"
	"mime"
	"net/http"
)

// Blob represents a seekable binary object in Blossom.
//
// Data is the underlying [io.ReadSeekCloser] containing the blob's content. Users must provide a
// seekable reader to support range requests (BUD-01), as well as automatic size and MIME detection.
type Blob struct {
	Data io.ReadSeekCloser
}

// Size returns the total size of the blob in bytes.
func (b Blob) Size() (int64, error) {
	if b.Data == nil {
		return 0, nil
	}

	size, err := b.Data.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, fmt.Errorf("failed to seek to end of blob: %w", err)
	}

	_, err = b.Data.Seek(0, io.SeekStart)
	if err != nil {
		return -1, fmt.Errorf("failed to rewind blob: %w", err)
	}
	return size, nil
}

// Type returns the content type of the blob by inspecting up to the first 512 bytes of its data.
// The returned string is suitable for use as a MIME type in HTTP headers (e.g. Content-Type).
// If the type cannot be determined, it returns the default "application/octet-stream" as specified by BUD-01.
func (b Blob) Type() (string, error) {
	if b.Data == nil {
		return "application/octet-stream", nil
	}

	sniff := make([]byte, 512)
	n, err := io.ReadFull(b.Data, sniff)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", fmt.Errorf("failed to read for MIME sniffing: %w", err)
	}

	_, err = b.Data.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to rewind blob after MIME sniffing: %w", err)
	}
	return http.DetectContentType(sniff[:n]), nil
}

// Extension returns the preferred file extension for the blob's content type.
// The returned extension will begin with a leading dot, as in ".html".
// If the blob's content type cannot be determined, or if no suitable extension is found, it returns ".bin".
func (b Blob) Extension() string {
	ct, err := b.Type()
	if err != nil {
		return ".bin"
	}
	return ExtFromType(ct)
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
