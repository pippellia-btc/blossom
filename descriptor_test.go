package blossom

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBlobDescriptorJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  BlobDescriptor
		err   bool
	}{
		// error cases
		{name: "invalid json", input: `{`, err: true},
		{name: "invalid hash", input: `{"sha256": "invalid", "size": 100, "url": "https://example.com/blob"}`, err: true},
		{name: "invalid size", input: `{"sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "size": "not a number", "url": "https://example.com/blob"}`, err: true},

		// success cases
		{
			name:  "minimal",
			input: `{"sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "size": 1024, "url": "https://example.com/blob"}`,
			want: BlobDescriptor{
				Hash: mustParseHash("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"),
				Size: 1024,
				URL:  "https://example.com/blob",
			},
		},
		{
			name:  "full",
			input: `{"sha256": "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789", "size": 2048, "url": "https://example.com/blob", "type": "image/png", "uploaded": 1234567890}`,
			want: BlobDescriptor{
				Hash:     mustParseHash("abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"),
				Size:     2048,
				URL:      "https://example.com/blob",
				Type:     "image/png",
				Uploaded: 1234567890,
			},
		},
		{
			name:  "with extra fields",
			input: `{"sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "size": 512, "url": "https://example.com/blob", "custom": "value"}`,
			want: BlobDescriptor{
				Hash:  mustParseHash("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"),
				Size:  512,
				URL:   "https://example.com/blob",
				Extra: map[string]json.RawMessage{"custom": json.RawMessage(`"value"`)},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var d BlobDescriptor
			err := json.Unmarshal([]byte(test.input), &d)

			if test.err {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(d, test.want) {
				t.Errorf("got %+v, want %+v", d, test.want)
			}

			// roundtrip
			data, err := json.Marshal(d)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}

			var d2 BlobDescriptor
			if err := json.Unmarshal(data, &d2); err != nil {
				t.Fatalf("Unmarshal roundtrip error: %v", err)
			}
			if !reflect.DeepEqual(d, d2) {
				t.Errorf("roundtrip mismatch: %+v != %+v", d, d2)
			}
		})
	}
}

func mustParseHash(s string) Hash {
	h, err := ParseHash(s)
	if err != nil {
		panic(err)
	}
	return h
}
