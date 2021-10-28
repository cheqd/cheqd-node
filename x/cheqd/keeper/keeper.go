package keeper

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc      codec.Codec
		storeKey sdk.StoreKey
	}
)

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// IsDidUsed checks if the did is used by DIDDoc
func (k Keeper) IsDidUsed(ctx sdk.Context, did string) bool {
	return k.HasDid(ctx, did)
}

func (k Keeper) EnsureDidIsNotUsed(ctx sdk.Context, did string) error {
	if k.HasDid(ctx, did) {
		return sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID is already used by DIDDoc %s", did))
	}

	return nil
}
