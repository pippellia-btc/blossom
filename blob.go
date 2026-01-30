package blossom

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
