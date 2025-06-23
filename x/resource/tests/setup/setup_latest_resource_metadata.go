package setup

import "github.com/cheqd/cheqd-node/x/resource/types"

func (s *TestSetup) QueryLatestResourceVersionMetadata(collectionID, name, resourceType string) (*types.QueryLatestResourceVersionMetadataResponse, error) {
	req := &types.QueryLatestResourceVersionMetadataRequest{
		CollectionId: collectionID,
		Name:         name,
		ResourceType: resourceType,
	}

	return s.ResourceQueryServer.LatestResourceVersionMetadata(s.StdCtx, req)
}
