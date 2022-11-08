package keeper

import (
	"fmt"

	"github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	Cdc      codec.BinaryCodec
	StoreKey storetypes.StoreKey
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey) *Keeper {
	return &Keeper{
		Cdc:      cdc,
		StoreKey: storeKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetStoreIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.StoreKey), resourcetypes.KeyPrefix(types.DidKey))
	return sdk.KVStorePrefixIterator(store, []byte{})
}
