package setup

import "github.com/cheqd/cheqd-node/x/resource/types"

func (s *TestSetup) QueryResource(collectionID, resourceID string) (*types.QueryResourceResponse, error) {
	req := &types.QueryGetResourceRequest{
		CollectionId: collectionID,
		Id:           resourceID,
	}

	return s.ResourceQueryServer.Resource(s.StdCtx, req)
}
