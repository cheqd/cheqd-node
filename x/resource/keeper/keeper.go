package keeper

import (
	"fmt"

	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
}

func NewKeeper(cdc codec.BinaryCodec, storeKey sdk.StoreKey) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
