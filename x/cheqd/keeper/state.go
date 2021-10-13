package keeper

import (
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasDidDoc checks if the did exist in the store
func (k Keeper) HasDidDoc(ctx sdk.Context, id string) error {
	if k.HasDid(ctx, id) {
		return sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID DOC already exists for DID %s", id))
	}

	if k.HasSchema(ctx, id) {
		return sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID DOC already exists for Schema %s", id))
	}

	if k.HasCredDef(ctx, id) {
		return sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID DOC already exists for CredDef %s", id))
	}

	return nil
}
