package keeper

import (
	errorsmod "cosmossdk.io/errors"
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func resource(ctx sdk.Context, keeper Keeper, cheqdKeeper didkeeper.Keeper, legacyQuerierCdc *codec.LegacyAmino, collectionID, id string) ([]byte, error) {
	queryServer := NewQueryServer(keeper, cheqdKeeper)

	resp, err := queryServer.Resource(sdk.WrapSDKContext(ctx), &types.QueryResourceRequest{CollectionId: collectionID, Id: id})
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
