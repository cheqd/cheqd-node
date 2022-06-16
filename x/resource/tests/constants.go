package tests

import (
	"crypto/sha256"

	"github.com/cheqd/cheqd-node/x/resource/types"
)

const (
	CLSchemaType        = "CL-Schema"
	SchemaData          = "{\"attr\":[\"name\",\"age\"]}"
	TestResourceName    = "Test Resource Name"
	JsonResourceType    = "application/json"
	ResourceId          = "988b0ab3-6a39-4598-83ec-b84c6cf8da15"
	IncorrectResourceId = "1234"

	NotFoundDIDIdentifier = "nfdnfdnfdnfdnfdd"
	ExistingDIDIdentifier = "eeeeeeeeeeeeeeee"
	ExistingDID           = "did:cheqd:test:" + ExistingDIDIdentifier
	ExistingDIDKey        = ExistingDID + "#key-1"
)

func ExistingResource() types.Resource {
	data := []byte(SchemaData)
	checksum := string(sha256.New().Sum(data))
	return types.Resource{
		CollectionId: ExistingDIDIdentifier,
		Id:           "a09abea0-22e0-4b35-8f70-9cc3a6d0b5fd",
		Name:         "Existing Resource Name",
		ResourceType: CLSchemaType,
		MimeType:     JsonResourceType,
		Data:         data,
		Checksum:     checksum,
	}
}
