package migrations

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This migration should be run after protobuf that's why we use new DidDocWithMetadata
func MigrateDidIndyStyle(sctx sdk.Context, mctx MigrationContext) error {
	store := sctx.KVStore(mctx.didStoreKey)

	var didDocWithMetadata didtypes.DidDocWithMetadata

	didKeys := ReadAllKeys(
		store,
		didutils.StrBytes(didtypes.DidDocVersionKey))

	for _, didKey := range didKeys {
		didDocWithMetadata = didtypes.DidDocWithMetadata{}

		mctx.codec.MustUnmarshal(store.Get(didKey), &didDocWithMetadata)

		// Make all dids indy style
		MoveToIndyStyleIds(&didDocWithMetadata)

		// Remove old DID Doc
		store.Delete(didKey)

		// Set new DID Doc
		err := mctx.didKeeperNew.AddNewDidDocVersion(&sctx, &didDocWithMetadata)
		if err != nil {
			return err
		}
	}

	return nil
}
