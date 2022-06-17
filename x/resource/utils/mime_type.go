package utils

import (
	"errors"
	"strings"

	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

var AllowedMimeTypes = []string {"application/json", "image/png"}

func ValidateMimeType(rt string) error {
	if ! cheqdUtils.Contains(AllowedMimeTypes, rt) {
		return errors.New(rt + " mime type is not allowed. Only " + strings.Join(AllowedResourceTypes, ",") + " .")
	}

	return nil
}