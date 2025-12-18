package blossom

import "io"

type Blob io.ReadSeekCloser

// BlobMeta groups metadata of a Blob.
type BlobMeta struct {
	Hash      Hash
	MIME      string
	Size      int64
	CreatedAt int64
}
