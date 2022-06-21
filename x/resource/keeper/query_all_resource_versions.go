package keeper

import (
	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func allResourceVersions(ctx sdk.Context, keeper Keeper, cheqdKeeper cheqdkeeper.Keeper, legacyQuerierCdc *codec.LegacyAmino, collectionId, name, resourceType, mimeType string) ([]byte, error) {
	queryServer := NewQueryServer(keeper, cheqdKeeper)

	resp, err := queryServer.AllResourceVersions(sdk.WrapSDKContext(ctx), &types.QueryGetAllResourceVersionsRequest{
		CollectionId: collectionId,
		Name:         name,
		ResourceType: resourceType,
		MimeType:     mimeType,
	})
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
