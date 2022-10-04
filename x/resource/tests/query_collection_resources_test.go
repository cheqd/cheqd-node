package tests

import (
	"fmt"
	"testing"

	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestQueryGetCollectionResources(t *testing.T) {
	keys := GenerateTestKeys()
	existingResource := ExistingResource()
	cases := []struct {
		valid    bool
		name     string
		msg      *types.QueryGetCollectionResourcesRequest
		response *types.QueryGetCollectionResourcesResponse
		errMsg   string
	}{
		{
			valid: true,
			name:  "Valid: Works",
			msg: &types.QueryGetCollectionResourcesRequest{
				CollectionId: ExistingDIDIdentifier,
			},
			response: &types.QueryGetCollectionResourcesResponse{
				Resources: []*types.ResourceHeader{existingResource.Header},
			},
			errMsg: "",
		},
		{
			valid: false,
			name:  "Not Valid: DID Doc is not found",
			msg: &types.QueryGetCollectionResourcesRequest{
				CollectionId: NotFoundDIDIdentifier,
			},
			response: nil,
			errMsg:   fmt.Sprintf("did:cheqd:test:%s: DID Doc not found", NotFoundDIDIdentifier),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			resourceSetup := InitEnv(t, keys[ExistingDIDKey].PublicKey, keys[ExistingDIDKey].PrivateKey)

			queryResponse, err := resourceSetup.QueryServer.CollectionResources(sdk.WrapSDKContext(resourceSetup.Ctx), msg)

			if tc.valid {
				resources := queryResponse.Resources
				expectedResources := tc.response.Resources
				require.Nil(t, err)
				require.Equal(t, len(expectedResources), len(resources))
				for i, r := range resources {
					expectedResources[i].CollectionId = cheqdtypes.NormalizeIdentifier(expectedResources[i].CollectionId)
					expectedResources[i].Id = cheqdtypes.NormalizeIdentifier(expectedResources[i].Id)
					r.Created = expectedResources[i].Created
					require.Equal(t, expectedResources[i], r)
				}
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
