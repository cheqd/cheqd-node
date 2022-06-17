package keeper

//func getResource(ctx sdk.Context, collectionId string, id string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
//	queryServer := NewQueryServer(keeper)
//
//	resp, err := queryServer.Resource(sdk.WrapSDKContext(ctx), &types.QueryGetResourceRequest{CollectionId: collectionId, Id: id})
//	if err != nil {
//		return nil, err
//	}
//
//	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
//	if err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
//	}
//
//	return bz, nil
//}
