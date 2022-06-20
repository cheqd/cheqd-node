package tests

import (
	"crypto/sha256"

	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/stretchr/testify/require"
)

func CompareResources(t require.TestingT, expectedResource *types.Resource, resource types.Resource) {
	require.Equal(t, expectedResource.CollectionId, resource.CollectionId)
	require.Equal(t, expectedResource.Id, resource.Id)
	require.Equal(t, expectedResource.MimeType, resource.MimeType)
	require.Equal(t, expectedResource.ResourceType, resource.ResourceType)
	require.Equal(t, expectedResource.Data, resource.Data)
	require.Equal(t, expectedResource.Name, resource.Name)
	require.Equal(t, sha256.New().Sum(expectedResource.Data), resource.Checksum)
	require.Equal(t, expectedResource.PreviousVersionId, resource.PreviousVersionId)
	require.Equal(t, expectedResource.NextVersionId, resource.NextVersionId)
}
