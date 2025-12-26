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

// BlobMeta groups metadata of a [Blob].
type BlobMeta struct {
	Hash      Hash
	MIME      string
	Size      int64
	CreatedAt int64
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

// ContentType returns the content type of the blob by inspecting up to the first 512 bytes of its data.
// The returned string is suitable for use as a MIME type in HTTP headers (e.g. Content-Type).
// If the type cannot be determined, it returns the default "application/octet-stream" as specified by BUD-01.
func (b Blob) ContentType() (string, error) {
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

// MediaType returns the media type of the blob by parsing its content type and stripping any parameters.
// For example, if the content type is "text/html; charset=utf-8", the media type will be "text/html".
func (b Blob) MediaType() (string, error) {
	ct, err := b.ContentType()
	if err != nil {
		return "", err
	}

	media, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return "", fmt.Errorf("failed to parse media type %q: %w", ct, err)
	}
	return media, nil
}
