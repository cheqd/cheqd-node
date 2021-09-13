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

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	NymKey      = "Nym-value-"
	NymCountKey = "Nym-count-"
)

const (
	DidKey      = "Did-value-"
	DidCountKey = "Did-count-"
)

const (
	AttribKey      = "Attrib-value-"
	AttribCountKey = "Attrib-count-"
)

const (
	SchemaKey      = "Schema-value-"
	SchemaCountKey = "Schema-count-"
)

const (
	Cred_defKey      = "Cred_def-value-"
	Cred_defCountKey = "Cred_def-count-"
)
