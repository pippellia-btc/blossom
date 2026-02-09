package blossom

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Blob is anything that can be read (and closed) and has a known size and MIME type.
// This makes it possible to always write a Blob in an HTTP response with the
// correct 'Content-Type' and 'Content-Length' headers.
//
// For typical scenarios, Blobs can be created with already provided constructors
// such as [BlobFromFile], [BlobFromBytes], and [BlobFromStream].
// For more advanced use cases (e.g. encrypted blobs, thread-safe blobs),
// you can implement the [Blob] interface directly.
type Blob interface {
	io.ReadCloser

	// Type returns the content type (MIME) of the blob.
	Type() string

	// Size return the total size of the blob in bytes.
	Size() int64
}

// blob is a generic [Blob] implementation with a reader as the underlying data.
type blob struct {
	io.ReadCloser
	size int64
	typ  string
}

func (b blob) Size() int64  { return b.size }
func (b blob) Type() string { return b.typ }

// BlobFromStream creates a Blob from a reader with known size and content type.
// It's typically used when dealing with HTTP requests or responses,
// where the size and content type are known from the headers.
func BlobFromStream(r io.ReadCloser, size int64, contentType string) Blob {
	return blob{
		ReadCloser: r,
		size:       size,
		typ:        contentType,
	}
}

// seekableBlob is a [Blob] implementation with a seekable reader as the underlying data.
type seekableBlob struct {
	io.ReadSeekCloser
	size int64
	typ  string
}

func (b seekableBlob) Size() int64  { return b.size }
func (b seekableBlob) Type() string { return b.typ }

// BlobFromPath creates a Blob from a file path.
// The file size and content type are detected automatically.
// The caller is responsible for closing the returned Blob.
func BlobFromPath(path string) (Blob, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	blob, err := BlobFromFile(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	return blob, nil
}

// BlobFromFile creates a Blob from the given file.
// The file size and content type are detected automatically.
// The caller is responsible for closing the returned Blob.
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
	return seekableBlob{ReadSeekCloser: f, size: info.Size(), typ: typ}, nil
}

type nopCloser struct {
	io.ReadSeeker
}

func (_ nopCloser) Close() error {
	return nil
}

// BlobFromBytes creates a Blob from the given byte slice.
// The size and content type are detected automatically.
func BlobFromBytes(data []byte) Blob {
	return seekableBlob{
		ReadSeekCloser: nopCloser{bytes.NewReader(data)},
		size:           int64(len(data)),
		typ:            http.DetectContentType(data),
	}
}

// WriteBlob writes the entire blob to the response writer.
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

// ServeBlob serves the blob to the response writer.
// If the blob is seekable (implements [io.ReadSeeker]), ServeBlob automatically handles Range requests.
// Otherwise, it falls back to streaming the full blob with [WriteBlob].
func ServeBlob(w http.ResponseWriter, r *http.Request, b Blob) error {
	if seeker, ok := b.(io.ReadSeeker); ok {
		// If seekable, use http.ServeContent for full Range support.
		// We set the Content-Type to avoid any potential MIME type sniffing issues.
		// Unfortunately http.ServeContent doesn't return errors that we can report to the caller.
		w.Header().Set("Content-Type", b.Type())
		http.ServeContent(w, r, "", time.Time{}, seeker)
		return nil
	}
	return WriteBlob(w, b)
}
