package keeper

// func NewQuerier(k Keeper, cheqdKeeper didkeeper.Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
// 	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
// 		var (
// 			res []byte
// 			err error
// 		)

// 		switch path[0] {
// 		case types.QueryGetResource:
// 			return resource(ctx, k, cheqdKeeper, legacyQuerierCdc, path[1], path[2])
// 		case types.QueryGetResourceMetadata:
// 			return resourceMetadata(ctx, k, cheqdKeeper, legacyQuerierCdc, path[1], path[2])
// 		case types.QueryGetCollectionResources:
// 			return collectionResources(ctx, k, cheqdKeeper, legacyQuerierCdc, path[1])

// 		default:
// 			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
// 		}

// 		return res, err
// 	}
// }
