package setup

import "github.com/cheqd/cheqd-node/x/resource/types"

func (s *TestSetup) CollectionResources(collectionId string) (*types.QueryGetCollectionResourcesResponse, error) {
	req := &types.QueryGetCollectionResourcesRequest{
		CollectionId: collectionId,
	}

	return s.ResourceQueryServer.CollectionResources(s.StdCtx, req)
}
