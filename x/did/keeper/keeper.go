package keeper

import (
	"fmt"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      storetypes.StoreKey
		paramSpace    types.ParamSubspace
		accountKeeper types.AccountKeeper
		bankkeeper    types.BankKeeper
	}
)

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, paramSpace types.ParamSubspace, ak types.AccountKeeper, bk types.BankKeeper) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    paramSpace,
		accountKeeper: ak,
		bankkeeper:    bk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
