package blossom

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// WriteBlob writes the blob to the response writer.
// It automatically sets the Content-Type and Content-Length headers according to BUD-01.
func WriteBlob(w http.ResponseWriter, b Blob) error {
	ct := b.Type()
	size := b.Size()

	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))

	written, err := io.Copy(w, b)
	if err != nil {
		return err
	}
	if written != size {
		return fmt.Errorf("copied size mismatch: expected %d, wrote %d", size, written)
	}
	return nil
}

// WriteError writes the error to the http response. If the reason is non-empty,
// it writes it to the "X-Reason" header as per BUD-01.
func WriteError(w http.ResponseWriter, e Error) {
	if e.Reason != "" {
		w.Header().Set("X-Reason", e.Reason)
	}
	http.Error(w, "", e.Code)
}
