package v1

const (
	// ModuleName defines the module name
	ModuleName = "cheqd"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	DidMethod = ModuleName
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	DidKey          = "did:"
	DidCountKey     = "did-count:"
	DidNamespaceKey = "did-namespace:"
)
