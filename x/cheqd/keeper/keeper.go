package keeper

import (
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
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
		memKey   sdk.StoreKey
	}
)

func NewKeeper(cdc codec.Codec, storeKey, memKey sdk.StoreKey) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// IsDidUsed checks if the did is used by DIDDoc, Schema or CredDef
func (k Keeper) IsDidUsed(ctx sdk.Context, did string) bool {
	if k.HasDid(ctx, did) || k.HasSchema(ctx, utils.GetSchemaFromDid(did)) ||k.HasCredDef(ctx, utils.GetCredDefFromDid(did)) {
		return true
	}

	return false
}

func (k Keeper) EnsureDidIsNotUsed(ctx sdk.Context, did string) error {
	if k.HasDid(ctx, did) {
		return sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID is already used by DIDDoc %s", did))
	}

	if k.HasSchema(ctx, utils.GetSchemaFromDid(did)) {
		return sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID is already used by Schema %s", did))
	}

	if k.HasCredDef(ctx, utils.GetCredDefFromDid(did)) {
		return sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID is already used by CredDef %s", did))
	}

	return nil
}
