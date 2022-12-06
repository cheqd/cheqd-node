package migrations

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateDidSimple(sctx sdk.Context, mctx MigrationContext, apply func(didDocWithMetadata *didtypes.DidDocWithMetadata)) error {
	store := sctx.KVStore(mctx.didStoreKey)

	// Reset counter
	mctx.didKeeperNew.SetDidDocCount(&sctx, 0)

	// Colect all did doc vsersions
	var allDidDocVersions []didtypes.DidDocWithMetadata

	mctx.didKeeperNew.IterateAllDidDocVersions(&sctx, func(metadata didtypes.DidDocWithMetadata) bool {
		allDidDocVersions = append(allDidDocVersions, metadata)
		return true
	})

	// Iterate and migrate did docs. We can use single loop for removing old values, migration
	// and writing new values because there is only one version of each diddoc in the store
	for _, version := range allDidDocVersions {
		// Remove last version pointer
		latestVersionKey := didtypes.GetLatestDidDocVersionKey(version.DidDoc.Id)
		store.Delete(latestVersionKey)

		// Remove version
		versionKey := didtypes.GetDidDocVersionKey(version.DidDoc.Id, version.Metadata.VersionId)
		store.Delete(versionKey)

		// Migrate
		apply(&version)

		// Create as a new did doc
		err := mctx.didKeeperNew.AddNewDidDocVersion(&sctx, &version)
		if err != nil {
			return err
		}
	}

	return nil
}
