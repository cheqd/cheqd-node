package setup

import "github.com/cheqd/cheqd-node/x/resource/types"

func (s *TestSetup) QueryLatestResourceVersion(collectionID, name, resourceType string) (*types.QueryLatestResourceVersionResponse, error) {
	req := &types.QueryLatestResourceVersionRequest{
		CollectionId: collectionID,
		Name:         name,
		ResourceType: resourceType,
	}

	return s.ResourceQueryServer.LatestResourceVersion(s.StdCtx, req)
}
