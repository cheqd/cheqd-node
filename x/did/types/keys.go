package types

import (
	"cosmossdk.io/collections"
)

const (
	// ModuleName defines the module name
	ModuleName = "cheqd"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	DidMethod = ModuleName
)

// State:
// did-namespace: -> <did-namespace>
// did-count: -> <did-count>
// did-latest:<did> -> <latest-version>
// did-version:<did>:<version> -> <did-doc>

var (
	DidDocCountKeyPrefix         = collections.NewPrefix("did-count:")
	DidNamespaceKeyPrefix        = collections.NewPrefix("did-namespace:")
	LatestDidDocVersionKeyPrefix = collections.NewPrefix("did-latest:")
	DidDocVersionKeyPrefix       = collections.NewPrefix("did-version:")
	ParamStoreKeyFeeParams       = collections.NewPrefix("feeparams")
)
