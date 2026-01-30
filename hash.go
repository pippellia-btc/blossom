package blossom

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
)

// Hash represents a 32-byte SHA-256 hash.
// It's the default and only hash type supported by Blossom.
type Hash [32]byte

// Hex converts the hash into its hexadecimal representation.
func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

// String converts the hash into its hexadecimal representation.
func (h Hash) String() string {
	return h.Hex()
}

// Value implements the [driver.Valuers] interface so it can serialize itself as a hexadecimal string.
func (h Hash) Value() (driver.Value, error) {
	return h.Hex(), nil
}

// Scan implements the [sql.Scanner] interface so it can deserialize itself.
func (h *Hash) Scan(src any) error {
	switch s := src.(type) {
	case string:
		if len(s) != 64 {
			return fmt.Errorf("invalid hash length: %d", len(s))
		}
		b, err := hex.DecodeString(s)
		if err != nil {
			return err
		}
		copy(h[:], b)
		return nil

	case []byte:
		if len(s) != 32 {
			return fmt.Errorf("invalid hash length: %d", len(s))
		}
		copy(h[:], s)
		return nil

	case nil:
		return fmt.Errorf("NULL cannot be scanned into Hash")

	default:
		return fmt.Errorf("cannot scan %T into Hash", src)
	}
}

// MarshalJSON implements [json.Marshaler].
func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Hex())
}

// UnmarshalJSON implements [json.Unmarshaler].
// It decodes directly into the Hash bytes without intermediate allocations.
func (h *Hash) UnmarshalJSON(data []byte) error {
	// JSON string: 64 hex chars + 2 quotes = 66 bytes
	if len(data) != 66 {
		return fmt.Errorf("invalid hash length: expected 64 hex characters, got %d", len(data)-2)
	}
	if data[0] != '"' || data[65] != '"' {
		return errors.New("invalid JSON string for hash")
	}

	if _, err := hex.Decode(h[:], data[1:65]); err != nil {
		return fmt.Errorf("invalid hash: %w", err)
	}
	return nil
}

// ComputeHash of the provided data, by calling the cryptographically secure
// sha256 implementation of the standard library.
func ComputeHash(data []byte) Hash {
	return sha256.Sum256(data)
}

// NewHasher returns a new hasher that computes a hash incrementally by calling
// the cryptographically secure sha256 implementation of the standard library.
func NewHasher() hash.Hash {
	return sha256.New()
}

// ParseHash from the hexadecimal input string.
func ParseHash(input string) (Hash, error) {
	if len(input) != 64 {
		return Hash{}, errors.New("input lenght must be exactly 64 characters")
	}

	var hash Hash
	b, err := hex.DecodeString(input)
	if err != nil {
		return Hash{}, fmt.Errorf("failed to parsh hash: %w", err)
	}

	copy(hash[:], b)
	return hash, nil
}

// hexValidTable is a lookup table for fast hex character validation.
// It maps ASCII bytes to whether they are valid hex characters.
var hexValidTable = [256]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true, '8': true, '9': true,
	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true,
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true,
}

// ValidateHash checks if the input string is a valid SHA-256 hash, returning an error if it is not.
// It validates that the string is exactly 64 hexadecimal characters without decoding the entire string.
func ValidateHash(input string) error {
	if len(input) != 64 {
		return errors.New("input lenght must be exactly 64 characters")
	}

	for i := 0; i < len(input); i++ {
		if !hexValidTable[input[i]] {
			return fmt.Errorf("invalid hex character at position %d: %c", i, input[i])
		}
	}
	return nil
}
