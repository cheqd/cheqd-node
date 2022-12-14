package migrations

import (
	"fmt"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateDidSimple(sctx sdk.Context, mctx MigrationContext, apply func(didDocWithMetadata *didtypes.DidDocWithMetadata)) error {
	sctx.Logger().Debug("MigrateDidSimple: Starting migration")

	store := sctx.KVStore(mctx.didStoreKey)

	sctx.Logger().Debug("MigrateDidSimple: Reset counter")
	// Reset counter
	mctx.didKeeperNew.SetDidDocCount(&sctx, 0)

	// Collect all DIDDoc versions
	var allDidDocVersions []didtypes.DidDocWithMetadata

	sctx.Logger().Debug("MigrateDidSimple: Iterate through all DIDDocs")
	mctx.didKeeperNew.IterateAllDidDocVersions(&sctx, func(metadata didtypes.DidDocWithMetadata) bool {
		allDidDocVersions = append(allDidDocVersions, metadata)
		return true
	})

	// Iterate and migrate did docs. We can use single loop for removing old values, migration
	// and writing new values because there is only one version of each diddoc in the store
	for _, version := range allDidDocVersions {
		sctx.Logger().Debug(fmt.Sprintf("Old version of DIDDoc: %s", string(mctx.codec.MustMarshalJSON(&version))))

		// Remove last version pointer
		latestVersionKey := didtypes.GetLatestDidDocVersionKey(version.DidDoc.Id)
		store.Delete(latestVersionKey)

		// Remove version
		versionKey := didtypes.GetDidDocVersionKey(version.DidDoc.Id, version.Metadata.VersionId)
		store.Delete(versionKey)

		// Migrate
		apply(&version)

		sctx.Logger().Debug(fmt.Sprintf("New version of DIDDoc: %s", string(mctx.codec.MustMarshalJSON(&version))))

		// Create as a new did doc
		err := mctx.didKeeperNew.AddNewDidDocVersion(&sctx, &version)
		if err != nil {
			return err
		}
	}
	sctx.Logger().Debug("MigrateDidSimple: Execution completed")

	return nil
}
