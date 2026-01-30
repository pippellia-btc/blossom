package blossom

import (
	"fmt"
	"io"
	"mime"
	"net/http"
)

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
