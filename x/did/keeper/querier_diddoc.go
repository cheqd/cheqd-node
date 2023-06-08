package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func getDidDoc(ctx sdk.Context, id string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	queryServer := NewQueryServer(keeper)

	resp, err := queryServer.DidDoc(sdk.WrapSDKContext(ctx), &types.QueryDidDocRequest{Id: id})
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
