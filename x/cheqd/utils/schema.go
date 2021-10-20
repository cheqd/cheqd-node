package utils

import "github.com/cheqd/cheqd-node/x/cheqd/utils/strings"

var SchemaSuffix = "?service=CL-Schema"
var SchemaSuffixLen = len(SchemaSuffix)

var AllowedSchemaType = []string{"CL-Schema"}

func IsNotSchemaType(schemaType string) bool {
	return !strings.Include(AllowedSchemaType, schemaType)
}

func IsSchema(did string) bool {
	return len(did) >= SchemaSuffixLen && did[len(did)-SchemaSuffixLen:] == SchemaSuffix
}

func GetDidFromSchema(schema string) string {
	return schema[:len(schema)-SchemaSuffixLen]
}

func GetSchemaFromDid(did string) string {
	return did + SchemaSuffix
}
