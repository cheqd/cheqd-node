package utils

import (
	"errors"
	"strings"

	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

var AllowedMediaTypes = []string{
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

func IsValidMediaType(rt string) bool {
	return cheqdUtils.Contains(AllowedMediaTypes, rt)
}

func ValidateMediaType(rt string) error {
	if !IsValidMediaType(rt) {
		return errors.New(rt + " mime type is not allowed. Only " + strings.Join(AllowedMediaTypes, ", "))
	}

	return nil
}
