package utils

var AllowedSchemaType = []string{"CL-Schema"}

func IsNotSchemaType(schemaType string) bool {
	return !StringArrayContains(AllowedSchemaType, schemaType)
}
