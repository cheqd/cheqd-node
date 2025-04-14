package v5

import (
	"cosmossdk.io/core/store"

	"github.com/cheqd/cheqd-node/x/did/exported"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "did"
)

var ParamsKey = []byte{0x09}

// MigrateStore migrates the x/distribution module state from the consensus version 2 to
// version 3. Specifically, it takes the parameters that are currently stored
// and managed by the x/params module and stores them directly into the x/distribution
// module state.
func MigrateStore(ctx sdk.Context, storeService store.KVStoreService, legacySubspace exported.Subspace, cdc codec.BinaryCodec) error {
	store := storeService.OpenKVStore(ctx)
	var currParams types.FeeParams
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.ValidateBasic(); err != nil {
		return err
	}

	bz, err := cdc.Marshal(&currParams)
	if err != nil {
		return err
	}

	return store.Set(ParamsKey, bz)
}
