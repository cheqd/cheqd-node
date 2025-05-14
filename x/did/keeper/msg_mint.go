package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k MsgServer) Mint(goCtx context.Context, req *types.MsgMint) (res *types.MsgMintResponse, err error) {
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}
	if k.authority != req.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err = k.bankkeeper.MintCoins(ctx, types.ModuleName, req.Amount)
	if err != nil {
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(req.ToAddress)
	if err != nil {
		return nil, err
	}

	err = k.bankkeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, req.Amount)
	if err != nil {
		return nil, err
	}

	return
}
