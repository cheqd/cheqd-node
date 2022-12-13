package helpers

import (
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/google/uuid"
)

func MigrateUUIDDid(did string) string {
	method, namespace, id := didutils.MustSplitDID(did)
	return didutils.JoinDID(method, namespace, MigrateUUIDId(id))
}

func MigrateUUIDId(id string) string {
	// If id is not UUID it should not be changed
	if !didutils.IsValidUUID(id) {
		return id
	}

	// If uuid is already normalized, it should not be changed
	if didutils.NormalizeUUID(id) == id {
		return id
	}

	newId := uuid.NewSHA1(uuid.Nil, []byte(id))
	return didutils.NormalizeUUID(newId.String())
}
