package keeper

import (
	// this line is used by starport scaffolding # 1
	"github.com/cheqd/cheqd-node/x/cheqd/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)

		switch path[0] {
		// this line is used by starport scaffolding # 2
		case types.QueryGetCredDef:
			return getCredDef(ctx, path[1], k, legacyQuerierCdc)

		case types.QueryListCredDef:
			return listCredDef(ctx, k, legacyQuerierCdc)

		case types.QueryGetSchema:
			return getSchema(ctx, path[1], k, legacyQuerierCdc)

		case types.QueryListSchema:
			return listSchema(ctx, k, legacyQuerierCdc)

		case types.QueryGetAttrib:
			return getAttrib(ctx, path[1], k, legacyQuerierCdc)

		case types.QueryListAttrib:
			return listAttrib(ctx, k, legacyQuerierCdc)

		case types.QueryGetDid:
			return getDid(ctx, path[1], k, legacyQuerierCdc)

		case types.QueryListDid:
			return listDid(ctx, k, legacyQuerierCdc)

		case types.QueryGetNym:
			return getNym(ctx, path[1], k, legacyQuerierCdc)

		case types.QueryListNym:
			return listNym(ctx, k, legacyQuerierCdc)

		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}
