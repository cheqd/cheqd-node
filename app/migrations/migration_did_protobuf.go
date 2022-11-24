package migrations

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesV1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateDidProtobuf(sctx sdk.Context, mctx MigrationContext) error {
	var didKeys []IteratorKey

	ir := codectypes.NewInterfaceRegistry()

	ir.RegisterInterface("StateValueData", (*didtypesV1.StateValueData)(nil))
	ir.RegisterImplementations((*didtypesV1.StateValueData)(nil), &didtypesV1.Did{})

	CdcV1 := codec.NewProtoCodec(ir)

	didKeys = CollectAllKeys(sctx, mctx.didStoreKey, didutils.StrBytes(didtypesV1.DidKey))

	store := prefix.NewStore(
		sctx.KVStore(mctx.didStoreKey),
		didutils.StrBytes(didtypesV1.DidKey))

	for _, didKey := range didKeys {
		var stateValue didtypesV1.StateValue
		var newDidDocWithMetadata didtypes.DidDocWithMetadata
		CdcV1.MustUnmarshal(store.Get(didKey), &stateValue)

		newDidDocWithMetadata, err := StateValueToDIDDocWithMetadata(&stateValue)

		if err != nil {
			return err
		}

		// Remove old DID Doc
		store.Delete(didKey)

		// Set new DID Doc
		err = mctx.didKeeper.AddNewDidDocVersion(&sctx, &newDidDocWithMetadata)
		if err != nil {
			return err
		}
	}

	return nil
}
