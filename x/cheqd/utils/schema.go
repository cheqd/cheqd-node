package utils

import "github.com/cheqd/cheqd-node/x/cheqd/utils/strings"

var AllowedSchemaType = []string{"CL-Schema"}

func IsNotSchemaType(schemaType string) bool {
	return !strings.Include(AllowedSchemaType, schemaType)
}
