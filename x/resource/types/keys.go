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

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	ResourceMetadataKey = "resource-metadata:"
	ResourceDataKey     = "resource-data:"
	ResourceCountKey    = "resource-count:"
)
