package keeper

import (
	"fmt"

	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeKey     storetypes.StoreKey
	paramSpace   types.ParamSubspace
	portKeeper   types.PortKeeper
	scopedKeeper exported.ScopedKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, paramSpace types.ParamSubspace, portKeeper types.PortKeeper, scopedKeeper exported.ScopedKeeper) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		paramSpace:   paramSpace,
		portKeeper:   portKeeper,
		scopedKeeper: scopedKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
