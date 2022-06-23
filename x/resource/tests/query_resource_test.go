package tests

import (

	// "crypto/sha256"
	"crypto/sha256"
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
				Id:           existingResource.Header.Id,
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
				Id:           existingResource.Header.Id,
			},
			response: nil,
			errMsg:   fmt.Sprintf("did:cheqd:test:%s: DID Doc not found", NotFoundDIDIdentifier),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			resourceSetup := InitEnv(t, keys[ExistingDIDKey].PublicKey, keys[ExistingDIDKey].PrivateKey)

			queryResponse, err := resourceSetup.QueryServer.Resource(sdk.WrapSDKContext(resourceSetup.Ctx), msg)

			if tc.valid {
				resource := queryResponse.Resource
				require.Nil(t, err)
				require.Equal(t, tc.response.Resource.Header.CollectionId, resource.Header.CollectionId)
				require.Equal(t, tc.response.Resource.Header.Id, resource.Header.Id)
				require.Equal(t, tc.response.Resource.Header.MimeType, resource.Header.MimeType)
				require.Equal(t, tc.response.Resource.Header.ResourceType, resource.Header.ResourceType)
				require.Equal(t, tc.response.Resource.Data, resource.Data)
				require.Equal(t, tc.response.Resource.Header.Name, resource.Header.Name)
				require.Equal(t, sha256.New().Sum(tc.response.Resource.Data), resource.Header.Checksum)
				require.Equal(t, tc.response.Resource.Header.PreviousVersionId, resource.Header.PreviousVersionId)
				require.Equal(t, tc.response.Resource.Header.NextVersionId, resource.Header.NextVersionId)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
