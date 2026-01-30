package blossom

import (
	"encoding/json"
	"errors"
	"fmt"
)

// BlobDescriptor is a description of a blossom blob.
// It contains the URL, hash, size, type, uploaded timestamp, and optional extra metadata.
// Learn more here: https://github.com/hzrd149/blossom/blob/master/buds/02.md#blob-descriptor
type BlobDescriptor struct {
	URL      string
	Hash     Hash
	Size     int64
	Type     string
	Uploaded int64

	// Extra fields to store arbitrary optional metadata.
	Extra map[string]json.RawMessage
}

// MarshalJSON implements [json.Marshaler].
// It serializes the descriptor to JSON with Extra fields flattened into the top level.
func (d BlobDescriptor) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"url":      d.URL,
		"sha256":   d.Hash.String(),
		"size":     d.Size,
		"type":     d.Type,
		"uploaded": d.Uploaded,
	}
	for k, v := range d.Extra {
		m[k] = v
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements [json.Unmarshaler].
// It deserializes JSON into the descriptor, placing unknown fields into Extra.
func (d *BlobDescriptor) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if raw, ok := m["url"]; ok {
		if err := json.Unmarshal(raw, &d.URL); err != nil {
			return fmt.Errorf("failed to parse url: %w", err)
		}
		delete(m, "url")
	}

	if raw, ok := m["sha256"]; ok {
		if err := json.Unmarshal(raw, &d.Hash); err != nil {
			return fmt.Errorf("invalid sha256: %w", err)
		}
		delete(m, "sha256")
	}

	if raw, ok := m["size"]; ok {
		if err := json.Unmarshal(raw, &d.Size); err != nil {
			return fmt.Errorf("invalid size: %w", err)
		}
		delete(m, "size")
	}

	if raw, ok := m["type"]; ok {
		if err := json.Unmarshal(raw, &d.Type); err != nil {
			return fmt.Errorf("invalid type: %w", err)
		}
		delete(m, "type")
	}

	if raw, ok := m["uploaded"]; ok {
		if err := json.Unmarshal(raw, &d.Uploaded); err != nil {
			return fmt.Errorf("invalid uploaded: %w", err)
		}
		delete(m, "uploaded")
	}

	// Remaining fields go into Extra
	if len(m) > 0 {
		d.Extra = m
	}
	return nil
}

// String returns a human-readable representation of the descriptor for debugging.
func (d BlobDescriptor) String() string {
	return fmt.Sprintf("BlobDescriptor{Hash: %s, Size: %d, Type: %s, URL: %s}",
		d.Hash, d.Size, d.Type, d.URL)
}

// Validate checks that the descriptor has valid required fields.
// It returns an error if the hash is zero, size is negative, or URL is empty.
func (d BlobDescriptor) Validate() error {
	if d.Hash == (Hash{}) {
		return errors.New("hash is required")
	}
	if d.Size < 0 {
		return errors.New("size cannot be negative")
	}
	if d.URL == "" {
		return errors.New("URL is required")
	}
	return nil
}
