package blossom

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
)

// extToType maps file extensions (without dot) to their MIME types.
// Provides deterministic behavior across all platforms.
var extToType = map[string]string{
	// Images
	"png":  "image/png",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"jpe":  "image/jpeg",
	"gif":  "image/gif",
	"webp": "image/webp",
	"svg":  "image/svg+xml",
	"ico":  "image/x-icon",
	"bmp":  "image/bmp",
	"tiff": "image/tiff",
	"tif":  "image/tiff",
	"avif": "image/avif",
	"heic": "image/heic",
	"heif": "image/heif",
	"apng": "image/apng",
	"jxl":  "image/jxl",

	// Video
	"mp4":  "video/mp4",
	"m4v":  "video/mp4",
	"webm": "video/webm",
	"ogv":  "video/ogg",
	"avi":  "video/x-msvideo",
	"mkv":  "video/x-matroska",
	"mov":  "video/quicktime",
	"wmv":  "video/x-ms-wmv",
	"flv":  "video/x-flv",
	"3gp":  "video/3gpp",
	"3g2":  "video/3gpp2",
	"mpeg": "video/mpeg",
	"mpg":  "video/mpeg",

	// Audio
	"mp3":  "audio/mpeg",
	"ogg":  "audio/ogg",
	"oga":  "audio/ogg",
	"wav":  "audio/wav",
	"weba": "audio/webm",
	"flac": "audio/flac",
	"aac":  "audio/aac",
	"m4a":  "audio/mp4",
	"wma":  "audio/x-ms-wma",
	"opus": "audio/opus",
	"midi": "audio/midi",
	"mid":  "audio/midi",
	"aiff": "audio/x-aiff",
	"aif":  "audio/x-aiff",

	// Documents
	"pdf":  "application/pdf",
	"doc":  "application/msword",
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"xls":  "application/vnd.ms-excel",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"ppt":  "application/vnd.ms-powerpoint",
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"odt":  "application/vnd.oasis.opendocument.text",
	"ods":  "application/vnd.oasis.opendocument.spreadsheet",
	"odp":  "application/vnd.oasis.opendocument.presentation",
	"rtf":  "application/rtf",
	"epub": "application/epub+zip",
	"mobi": "application/x-mobipocket-ebook",

	// Text
	"txt":  "text/plain",
	"html": "text/html",
	"htm":  "text/html",
	"css":  "text/css",
	"js":   "text/javascript",
	"mjs":  "text/javascript",
	"json": "application/json",
	"xml":  "application/xml",
	"csv":  "text/csv",
	"md":   "text/markdown",
	"ics":  "text/calendar",
	"vcf":  "text/vcard",

	// Archives
	"zip": "application/zip",
	"tar": "application/x-tar",
	"gz":  "application/gzip",
	"tgz": "application/gzip",
	"bz2": "application/x-bzip2",
	"xz":  "application/x-xz",
	"7z":  "application/x-7z-compressed",
	"rar": "application/vnd.rar",
	"lz":  "application/x-lzip",
	"zst": "application/zstd",

	// Fonts
	"woff":  "font/woff",
	"woff2": "font/woff2",
	"ttf":   "font/ttf",
	"otf":   "font/otf",
	"eot":   "application/vnd.ms-fontobject",

	// Other
	"bin":    "application/octet-stream",
	"exe":    "application/octet-stream",
	"dll":    "application/octet-stream",
	"so":     "application/octet-stream",
	"dylib":  "application/octet-stream",
	"swf":    "application/x-shockwave-flash",
	"jar":    "application/java-archive",
	"apk":    "application/vnd.android.package-archive",
	"dmg":    "application/x-apple-diskimage",
	"iso":    "application/x-iso9660-image",
	"wasm":   "application/wasm",
	"sqlite": "application/vnd.sqlite3",
	"db":     "application/vnd.sqlite3",
	"glb":    "model/gltf-binary",
	"gltf":   "model/gltf+json",
}

// typeToExt maps MIME types to their preferred file extensions.
// For types with multiple possible extensions, the most common one is used.
var typeToExt = map[string]string{
	// Images
	"image/png":     "png",
	"image/jpeg":    "jpg",
	"image/gif":     "gif",
	"image/webp":    "webp",
	"image/svg+xml": "svg",
	"image/x-icon":  "ico",
	"image/bmp":     "bmp",
	"image/tiff":    "tiff",
	"image/avif":    "avif",
	"image/heic":    "heic",
	"image/heif":    "heif",
	"image/apng":    "apng",
	"image/jxl":     "jxl",

	// Video
	"video/mp4":        "mp4",
	"video/webm":       "webm",
	"video/ogg":        "ogv",
	"video/x-msvideo":  "avi",
	"video/x-matroska": "mkv",
	"video/quicktime":  "mov",
	"video/x-ms-wmv":   "wmv",
	"video/x-flv":      "flv",
	"video/3gpp":       "3gp",
	"video/3gpp2":      "3g2",
	"video/mpeg":       "mpeg",

	// Audio
	"audio/mpeg":     "mp3",
	"audio/ogg":      "ogg",
	"audio/wav":      "wav",
	"audio/webm":     "weba",
	"audio/flac":     "flac",
	"audio/aac":      "aac",
	"audio/mp4":      "m4a",
	"audio/x-ms-wma": "wma",
	"audio/opus":     "opus",
	"audio/midi":     "midi",
	"audio/x-aiff":   "aiff",

	// Documents
	"application/pdf":    "pdf",
	"application/msword": "doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "docx",
	"application/vnd.ms-excel": "xls",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         "xlsx",
	"application/vnd.ms-powerpoint":                                             "ppt",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
	"application/vnd.oasis.opendocument.text":                                   "odt",
	"application/vnd.oasis.opendocument.spreadsheet":                            "ods",
	"application/vnd.oasis.opendocument.presentation":                           "odp",
	"application/rtf":                "rtf",
	"application/epub+zip":           "epub",
	"application/x-mobipocket-ebook": "mobi",

	// Text
	"text/plain":             "txt",
	"text/html":              "html",
	"text/css":               "css",
	"text/javascript":        "js",
	"application/javascript": "js",
	"application/json":       "json",
	"application/xml":        "xml",
	"text/xml":               "xml",
	"text/csv":               "csv",
	"text/markdown":          "md",
	"text/calendar":          "ics",
	"text/vcard":             "vcf",

	// Archives
	"application/zip":              "zip",
	"application/x-tar":            "tar",
	"application/gzip":             "gz",
	"application/x-gzip":           "gz",
	"application/x-bzip2":          "bz2",
	"application/x-xz":             "xz",
	"application/x-7z-compressed":  "7z",
	"application/vnd.rar":          "rar",
	"application/x-rar-compressed": "rar",
	"application/x-lzip":           "lz",
	"application/zstd":             "zst",

	// Fonts
	"font/woff":                     "woff",
	"font/woff2":                    "woff2",
	"font/ttf":                      "ttf",
	"font/otf":                      "otf",
	"application/vnd.ms-fontobject": "eot",

	// Other
	"application/octet-stream":                "bin",
	"application/x-shockwave-flash":           "swf",
	"application/java-archive":                "jar",
	"application/vnd.android.package-archive": "apk",
	"application/x-apple-diskimage":           "dmg",
	"application/x-iso9660-image":             "iso",
	"application/wasm":                        "wasm",
	"application/vnd.sqlite3":                 "sqlite",
	"application/x-sqlite3":                   "sqlite",
	"model/gltf-binary":                       "glb",
	"model/gltf+json":                         "gltf",
}

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
// The returned extension does not include a leading dot, e.g. "html" not ".html".
// If no suitable extension is found, it returns "bin".
func ExtFromType(contentType string) string {
	m, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "bin"
	}

	if ext, ok := typeToExt[m]; ok {
		return ext
	}
	return "bin"
}

// TypeFromExt returns the MIME type for the given file extension.
// The extension can be provided with or without a leading dot.
// If no suitable type is found, it returns "application/octet-stream".
func TypeFromExt(ext string) string {
	ext = strings.TrimPrefix(ext, ".")
	ext = strings.ToLower(ext)
	if m, ok := extToType[ext]; ok {
		return m
	}
	return "application/octet-stream"
}
