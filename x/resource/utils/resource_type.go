package utils

import (
	"errors"
	"strings"

	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

var AllowedResourceTypes = []string {"CL-Schema", "JSONSchema2020"}

func ValidateResourceType(rt string) error {
	if ! cheqdUtils.Contains(AllowedResourceTypes, rt) {
		return errors.New(rt + " resource type is not allowed. Only " + strings.Join(AllowedResourceTypes, ",") + " .")
	}

	return nil
}