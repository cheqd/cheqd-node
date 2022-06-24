package utils

import (
	"errors"
	"strings"

	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

var AllowedMimeTypes = []string{
	"application/json",
	"application/octet-stream",
	"text/plain",
	"image/apng",
	"image/avif",
	"image/gif",
	"image/jpeg",
	"image/png",
	"image/svg+xml",
	"image/webp",
}

func IsValidMimeType(rt string) bool {
	return cheqdUtils.Contains(AllowedMimeTypes, rt)
}

func ValidateMimeType(rt string) error {
	if !IsValidMimeType(rt) {
		return errors.New(rt + " mime type is not allowed. Only " + strings.Join(AllowedMimeTypes, ", "))
	}

	return nil
}
