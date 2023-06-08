package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func getDidDocVersion(ctx sdk.Context, id, version string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	queryServer := NewQueryServer(keeper)

	resp, err := queryServer.DidDocVersion(sdk.WrapSDKContext(ctx), &types.QueryDidDocVersionRequest{Id: id, Version: version})
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
