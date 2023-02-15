package types

const (
	// ModuleName defines the module name
	ModuleName = "resource"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

const (
	ResourceMetadataKey = "resource-metadata:"
	ResourceDataKey     = "resource-data:"
	ResourceCountKey    = "resource-count:"
)

// GetResourceDataKey returns the byte representation of resource key
func GetResourceDataKey(collectionID string, id string) []byte {
	return []byte(ResourceDataKey + collectionID + ":" + id)
}

// GetResourceMetadataKey returns the byte representation of resource key
func GetResourceMetadataKey(collectionID string, id string) []byte {
	return []byte(ResourceMetadataKey + collectionID + ":" + id)
}

// GetResourceMetadataCollectionPrefix used to iterate over all resource metadatas in a collection
func GetResourceMetadataCollectionPrefix(collectionID string) []byte {
	return []byte(ResourceMetadataKey + collectionID + ":")
}
