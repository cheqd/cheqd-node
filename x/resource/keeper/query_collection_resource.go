package keeper

//import (
//	"github.com/cheqd/cheqd-node/x/resource/types"
//	"github.com/cosmos/cosmos-sdk/codec"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
//
//	abci "github.com/tendermint/tendermint/abci/types"
//)
//
//func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
//
//
//	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
//		var (
//			res []byte
//			err error
//		)
//
//		switch path[0] {
//		case types.QueryGetResource:
//			return k.Res(ctx, path[1], path[2], k, legacyQuerierCdc)
//		// case types.QueryGetCollectionResources:
//		// 	return getCollectionResources(ctx, path[1], k, legacyQuerierCdc)
//		// case types.QueryGetAllResourceVersions:
//		// 	return getAllResourceVersions(ctx, path[1], k, legacyQuerierCdc)
//
//		default:
//			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
//		}
//
//		return res, err
//	}
//}
