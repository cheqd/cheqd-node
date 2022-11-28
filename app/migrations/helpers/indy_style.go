package helpers

import (
	"crypto/sha256"

	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/mr-tron/base58"
)

func MigrateIndyStyleDid(did string) string {
	method, namespace, id := didutils.MustSplitDID(did)
	return didutils.JoinDID(method, namespace, MigrateIndyStyleId(id))
}

func MigrateIndyStyleId(id string) string {
	// If id is UUID it should not be changed
	if didutils.IsValidUUID(id) {
		return id
	}

	// Get Hash from current id to make a 32-symbol string
	hash := sha256.Sum256([]byte(id))
	// Indy-style identifier is 16-byte base58 string
	return base58.Encode(hash[:16])
}
