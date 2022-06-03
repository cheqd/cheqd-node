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

	ResourceMethod = ModuleName
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	ResourceKey          = "resource:"
	ResourceCountKey     = "resource-count:"
	ResourceNamespaceKey = "resource-namespace:"
)
