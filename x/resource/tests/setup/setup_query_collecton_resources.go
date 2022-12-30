package setup

import "github.com/cheqd/cheqd-node/x/resource/types"

func (s *TestSetup) CollectionResources(collectionID string) (*types.QueryCollectionResourcesResponse, error) {
	req := &types.QueryGetCollectionResourcesRequest{
		CollectionId: collectionID,
	}

	return s.ResourceQueryServer.CollectionResources(s.StdCtx, req)
}
