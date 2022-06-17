package tests

import (

	// "crypto/sha256"
	"fmt"
	"testing"

	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestQueryGetResource(t *testing.T) {
	keys := GenerateTestKeys()
	existingResource := ExistingResource()
	cases := []struct {
		valid    bool
		name     string
		msg      *types.QueryGetResourceRequest
		response *types.QueryGetResourceResponse
		errMsg   string
	}{
		{
			valid: true,
			name:  "Valid: Works",
			msg: &types.QueryGetResourceRequest{
				CollectionId: ExistingDIDIdentifier,
				Id:           existingResource.Id,
			},
			response: &types.QueryGetResourceResponse{
				Resource: &existingResource,
			},
			errMsg: "",
		},
		{
			valid: false,
			name:  "Not Valid: Resource is not found",
			msg: &types.QueryGetResourceRequest{
				CollectionId: ExistingDIDIdentifier,
				Id:           ResourceId,
			},
			response: nil,
			errMsg:   fmt.Sprintf("resource %s:%s: not found", ExistingDIDIdentifier, ResourceId),
		},
		{
			valid: false,
			name:  "Not Valid: DID Doc is not found",
			msg: &types.QueryGetResourceRequest{
				CollectionId: NotFoundDIDIdentifier,
				Id:           existingResource.Id,
			},
			response: nil,
			errMsg:   fmt.Sprintf("resource %s:%s: not found", NotFoundDIDIdentifier, existingResource.Id),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			resourceSetup := InitEnv(t, keys[ExistingDIDKey].PublicKey, keys[ExistingDIDKey].PrivateKey)

			queryResponse, err := resourceSetup.ResourceKeeper.Resource(sdk.WrapSDKContext(resourceSetup.Ctx), msg)

			if tc.valid {
				resource := queryResponse.Resource
				require.Nil(t, err)
				require.Equal(t, tc.response.Resource.CollectionId, resource.CollectionId)
				require.Equal(t, tc.response.Resource.Id, resource.Id)
				require.Equal(t, tc.response.Resource.MimeType, resource.MimeType)
				require.Equal(t, tc.response.Resource.ResourceType, resource.ResourceType)
				require.Equal(t, tc.response.Resource.Data, resource.Data)
				require.Equal(t, tc.response.Resource.Name, resource.Name)
				// require.Equal(t, string(sha256.New().Sum(response.Resource.Data)), resource.Checksum)
				require.Equal(t, tc.response.Resource.PreviousVersionId, resource.PreviousVersionId)
				require.Equal(t, tc.response.Resource.NextVersionId, resource.NextVersionId)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
