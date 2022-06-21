package tests

import (

	// "crypto/sha256"
	// "crypto/ed25519"
	"fmt"
	"testing"

	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestQueryGetCollectionResources(t *testing.T) {
	keys := GenerateTestKeys()
	// existingResource := ExistingResource()
	cases := []struct {
		valid    bool
		name     string
		msg      *types.QueryGetCollectionResourcesRequest
		response *types.QueryGetCollectionResourcesResponse
		errMsg   string
	}{
		// {
		// 	valid: true,
		// 	name:  "Valid: Works",
		// 	msg: &types.QueryGetCollectionResourcesRequest{
		// 		CollectionId: ExistingDIDIdentifier,
		// 	},
		// 	response: &types.QueryGetCollectionResourcesResponse{
		// 		Resources: []*types.Resource{&existingResource},
		// 	},
		// 	errMsg: "",
		// },
		// {
		// 	valid: false,
		// 	name:  "Not Valid: Resource is not found",
		// 	msg: &types.QueryGetCollectionResources{
		// 		CollectionId: ExistingDIDIdentifier,
		// 		Id:           ResourceId,
		// 	},
		// 	response: nil,
		// 	errMsg:   fmt.Sprintf("resource %s:%s: not found", ExistingDIDIdentifier, ResourceId),
		// },
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

			// newResourcePayload := GenerateCreateResourcePayload(ExistingResource())
			// newResourcePayload.Id = ResourceId
			// didKey := map[string]ed25519.PrivateKey{
			// 	ExistingDIDKey: keys[ExistingDIDKey].PrivateKey,
			// }
			// createdResource, err := resourceSetup.SendCreateResource(newResourcePayload, didKey)
			// require.Nil(t, err)

			queryResponse, err := resourceSetup.QueryServer.CollectionResources(sdk.WrapSDKContext(resourceSetup.Ctx), msg)

			if tc.valid {
				resources := queryResponse.Resources
				expectedResources := tc.response.Resources
				// expectedResources := map[string]types.Resource {
				// 	existingResource.Id: existingResource,
				// 	createdResource.Id: *createdResource,
				// }
				require.Nil(t, err)
				require.Equal(t, len(expectedResources), len(resources))
				for i, r := range resources {
					CompareResources(t, expectedResources[i], *r)
				}
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
