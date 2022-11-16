package types

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

// State:
// did-namespace: -> <did-namespace>
// did-count: -> <did-count>
// did-latest:<did> -> <latest-version>
// did-version:<did>:<version> -> <did-doc>

const (
	LatestDidDocVersionKey = "did-latest:"
	DidDocVersionKey       = "did-version:"
	DidDocCountKey         = "did-count:"
	DidNamespaceKey        = "did-namespace:"
)

func GetLatestDidDocVersionKey(did string) []byte {
	return []byte(LatestDidDocVersionKey + did)
}

func GetDidDocVersionKey(did string, version string) []byte {
	return []byte(DidDocVersionKey + did + ":" + version)
}

func GetLatestDidDocVersionPrefix() []byte {
	return []byte(LatestDidDocVersionKey)
}

func GetDidDocVersionsPrefix(did string) []byte {
	return []byte(DidDocVersionKey + did + ":")
}
