package utils

import (
	"github.com/gabriel-vasile/mimetype"
)

func DetectMediaType(data []byte) string {
	mimetype.SetLimit(0) // No limit, whole file content used.

	// The result is always a valid MIME type, with application/octet-stream
	// returned when identification failed.
	return mimetype.Detect(data).String()
}
