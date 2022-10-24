package setup

import "github.com/cheqd/cheqd-node/x/resource/types"

func (s *TestSetup) QueryResource(collectionId, resourceId string) (*types.QueryGetResourceResponse, error) {
	req := &types.QueryGetResourceRequest{
		CollectionId: collectionId,
		Id:           resourceId,
	}

	return s.ResourceQueryServer.Resource(s.StdCtx, req)
}
