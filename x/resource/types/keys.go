package types

const (
	// ModuleName defines the module name
	ModuleName = "resource"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// Version support via IBC
	IBCVersion = "cheqd-resource-v3"

	// Port id that this module binds to
	ResourcePortID = "cheqdresource"
)

const (
	ResourceMetadataKey = "resource-metadata:"
	ResourceDataKey     = "resource-data:"
	ResourceCountKey    = "resource-count:"
	ResourcePortIDKey   = "resource-port-id:"
)
