package migrations

import (
	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This migration should be run after protobuf that's why we use new DidDocWithMetadata
func MigrateDidIndyStyle(sctx sdk.Context, mctx MigrationContext) error {
	store := sctx.KVStore(mctx.didStoreKey)

	// Reset counter
	mctx.didKeeperNew.SetDidDocCount(&sctx, 0)

	// Colect all did doc vsersions
	var allDidDocVersions []didtypes.DidDocWithMetadata

	mctx.didKeeperNew.IterateAllDidDocVersions(&sctx, func(metadata didtypes.DidDocWithMetadata) bool {
		allDidDocVersions = append(allDidDocVersions, metadata)
		return true
	})

	// Iterate and migrate did docs. We can use single loop for removing old values
	// and writing new values because there is only one version of did doc in the store
	for _, version := range allDidDocVersions {
		// Remove last version pointer
		latestVersionKey := didtypes.GetLatestDidDocVersionKey(version.DidDoc.Id)
		store.Delete(latestVersionKey)

		// Remove version
		versionKey := didtypes.GetDidDocVersionKey(version.DidDoc.Id, version.Metadata.VersionId)
		store.Delete(versionKey)

		// Migrate all dids, make them indy style
		newDid := helpers.MigrateIndyStyleDid(version.DidDoc.Id)
		version.ReplaceDids(version.DidDoc.Id, newDid)

		// Create as a new did doc
		err := mctx.didKeeperNew.AddNewDidDocVersion(&sctx, &version)
		if err != nil {
			return err
		}
	}

	return nil
}
