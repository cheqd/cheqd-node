package keeper

import (
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasDidDoc checks if the did exist in the store
func (k Keeper) HasDidDoc(ctx sdk.Context, did string) error {
	if k.HasDid(ctx, did) {
		return sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID DOC already exists for DID %s", did))
	}

	if k.HasSchema(ctx, did+"/schema") {
		return sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID DOC already exists for Schema %s", did))
	}

	if k.HasCredDef(ctx, did+"/credDef") {
		return sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID DOC already exists for CredDef %s", did))
	}

	return nil
}
