package setup

import "github.com/cheqd/cheqd-node/x/resource/types"

func (s *TestSetup) QueryResourceMetadata(collectionId, resourceId string) (*types.QueryGetResourceMetadataResponse, error) {
	req := &types.QueryGetResourceMetadataRequest{
		CollectionId: collectionId,
		Id:           resourceId,
	}

	return s.ResourceQueryServer.ResourceMetadata(s.StdCtx, req)
}
