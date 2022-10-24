package setup

import "github.com/cheqd/cheqd-node/x/resource/types"

func (s *TestSetup) AllResourceVersions(collectionId, name, resourceType string) (*types.QueryGetAllResourceVersionsResponse, error) {
	req := &types.QueryGetAllResourceVersionsRequest{
		CollectionId: collectionId,
		Name:         name,
		ResourceType: resourceType,
	}

	return s.ResourceQueryServer.AllResourceVersions(s.StdCtx, req)
}
