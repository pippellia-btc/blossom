package blossom

import (
	"fmt"
	"testing"
)

func TestExtFromType(t *testing.T) {
	tests := []struct {
		contentType string
		ext         string
	}{
		{"image/jpeg", ".jpg"},
		{"image/png", ".png"},
		{"text/html; charset=utf-8", ".html"},
		{"application/json", ".json"},
		{"application/octet-stream", ".bin"},
		{"unknown/type", ".bin"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case%d", i), func(t *testing.T) {
			ext := ExtFromType(test.contentType)
			if ext != test.ext {
				t.Errorf("ExtFromType(%q) = %q, want %q", test.contentType, ext, test.ext)
			}
		})
	}
}
