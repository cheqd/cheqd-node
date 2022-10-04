package tests

import (
	"crypto/ed25519"
	"fmt"
	"testing"

	cheqdtests "github.com/cheqd/cheqd-node/x/cheqd/tests"
	cheqdutils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func SendAnotherResourceVersion(t require.TestingT, resourceSetup TestSetup, keys map[string]cheqdtests.KeyPair) types.Resource {
	newResourcePayload := GenerateCreateResourcePayload(ExistingResource())
	newResourcePayload.Id = AnotherResourceId
	didKey := map[string]ed25519.PrivateKey{
		ExistingDIDKey: keys[ExistingDIDKey].PrivateKey,
	}
	newResourcePayload.Name = "AnotherResourceVersion"
	createdResource, err := resourceSetup.SendCreateResource(newResourcePayload, didKey)
	require.Nil(t, err)

	return *createdResource
}

func TestQueryGetAllResourceVersions(t *testing.T) {
	keys := GenerateTestKeys()
	existingResource := ExistingResource()
	cases := []struct {
		valid    bool
		name     string
		msg      *types.QueryGetAllResourceVersionsRequest
		response *types.QueryGetAllResourceVersionsResponse
		errMsg   string
	}{
		{
			valid: true,
			name:  "Valid: Works",
			msg: &types.QueryGetAllResourceVersionsRequest{
				CollectionId: ExistingDIDIdentifier,
				Name:         existingResource.Header.Name,
			},
			response: &types.QueryGetAllResourceVersionsResponse{
				Resources: []*types.ResourceHeader{existingResource.Header},
			},
			errMsg: "",
		},
		{
			valid: false,
			name:  "Not Valid: DID Doc is not found",
			msg: &types.QueryGetAllResourceVersionsRequest{
				CollectionId: NotFoundDIDIdentifier,
				Name:         existingResource.Header.Name,
			},
			response: nil,
			errMsg:   fmt.Sprintf("did:cheqd:test:%s: DID Doc not found", NotFoundDIDIdentifier),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			resourceSetup := InitEnv(t, keys[ExistingDIDKey].PublicKey, keys[ExistingDIDKey].PrivateKey)

			newResourcePayload := GenerateCreateResourcePayload(ExistingResource())
			newResourcePayload.Id = ResourceId
			didKey := map[string]ed25519.PrivateKey{
				ExistingDIDKey: keys[ExistingDIDKey].PrivateKey,
			}
			// Resource with the same version but another Id
			createdResource, err := resourceSetup.SendCreateResource(newResourcePayload, didKey)
			require.Nil(t, err)

			// Send another Resource but with another Name (should affect the version choosing)
			SendAnotherResourceVersion(t, resourceSetup, keys)

			queryResponse, err := resourceSetup.QueryServer.AllResourceVersions(sdk.WrapSDKContext(resourceSetup.Ctx), msg)

			if tc.valid {
				resources := queryResponse.Resources
				existingResource.Header.NextVersionId = createdResource.Header.Id
				expectedResources := map[string]types.Resource{
					existingResource.Header.Id: existingResource,
					createdResource.Header.Id:  *createdResource,
				}
				require.Nil(t, err)
				require.Equal(t, len(expectedResources), len(resources))
				for _, r := range resources {
					expectedResources[r.Id].Header.CollectionId = cheqdutils.NormalizeIdentifier(expectedResources[r.Id].Header.CollectionId)
					expectedResources[r.Id].Header.Id = cheqdutils.NormalizeIdentifier(expectedResources[r.Id].Header.Id)
					r.Created = expectedResources[r.Id].Header.Created
					require.Equal(t, r, expectedResources[r.Id].Header)
				}
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
