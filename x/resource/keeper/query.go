package keeper

import (
	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, cheqdKeeper cheqdkeeper.Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)

		switch path[0] {
		case types.QueryGetResource:
			return resource(ctx, k, cheqdKeeper, legacyQuerierCdc, path[1], path[2])
		case types.QueryGetCollectionResources:
			return collectionResources(ctx, k, cheqdKeeper, legacyQuerierCdc, path[1])
		case types.QueryGetAllResourceVersions:
			return allResourceVersions(ctx, k, cheqdKeeper, legacyQuerierCdc, path[1], path[2], path[3])

		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}
