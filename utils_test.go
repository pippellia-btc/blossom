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
		{"application/vnd.android.package-archive", "apk"},
		{"image/jpeg", "jpg"},
		{"image/png", "png"},
		{"text/html; charset=utf-8", "html"},
		{"application/json", "json"},
		{"application/octet-stream", "bin"},
		{"unknown/type", "bin"},
		{"video/mp4", "mp4"},
		{"audio/mpeg", "mp3"},
		{"application/pdf", "pdf"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case=%d", i), func(t *testing.T) {
			ext := ExtFromType(test.contentType)
			if ext != test.ext {
				t.Errorf("ExtFromType(%q) = %q, want %q", test.contentType, ext, test.ext)
			}
		})
	}
}

func TestTypeFromExt(t *testing.T) {
	tests := []struct {
		ext      string
		mimeType string
	}{
		{".jpg", "image/jpeg"},
		{"jpg", "image/jpeg"},
		{".jpeg", "image/jpeg"},
		{"JPEG", "image/jpeg"},
		{".PNG", "image/png"},
		{"png", "image/png"},
		{".mp4", "video/mp4"},
		{".mp3", "audio/mpeg"},
		{".pdf", "application/pdf"},
		{".html", "text/html"},
		{".htm", "text/html"},
		{".json", "application/json"},
		{".unknown", "application/octet-stream"},
		{"", "application/octet-stream"},
		{".bin", "application/octet-stream"},
		{".exe", "application/octet-stream"},
		{".apk", "application/vnd.android.package-archive"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case=%d", i), func(t *testing.T) {
			mimeType := TypeFromExt(test.ext)
			if mimeType != test.mimeType {
				t.Errorf("TypeFromExt(%q) = %q, want %q", test.ext, mimeType, test.mimeType)
			}
		})
	}
}
