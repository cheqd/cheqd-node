package keeper

import (
	cheqdkeeper "github.com/canow-co/cheqd-node/x/cheqd/keeper"
	"github.com/canow-co/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func resource(ctx sdk.Context, keeper Keeper, cheqdKeeper cheqdkeeper.Keeper, legacyQuerierCdc *codec.LegacyAmino, collectionId, id string) ([]byte, error) {
	queryServer := NewQueryServer(keeper, cheqdKeeper)

	resp, err := queryServer.Resource(sdk.WrapSDKContext(ctx), &types.QueryGetResourceRequest{CollectionId: collectionId, Id: id})
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
