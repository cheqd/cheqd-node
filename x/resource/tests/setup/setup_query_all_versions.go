package setup

import "github.com/cheqd/cheqd-node/x/resource/types"

func (s *TestSetup) AllResourceVersions(collectionId, name string) (*types.QueryGetAllResourceVersionsResponse, error) {
	req := &types.QueryGetAllResourceVersionsRequest{
		CollectionId: collectionId,
		Name:         name,
	}

	return s.ResourceQueryServer.AllResourceVersions(s.StdCtx, req)
}
