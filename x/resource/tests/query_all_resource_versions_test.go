package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"

	cheqdsetup "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"
	"github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

	It("Should return 2 versions for resource 1", func() {
		versions, err := setup.AllResourceVersions(alice.CollectionId, res1v1.Resource.Header.Name)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(2))

		ids := []string{versions.Resources[0].Id, versions.Resources[1].Id}

		Expect(ids).To(ContainElement(res1v1.Resource.Header.Id))
		Expect(ids).To(ContainElement(res1v2.Resource.Header.Id))
	})

	It("Should return 1 version for resource 2", func() {
		versions, err := setup.AllResourceVersions(alice.CollectionId, res2v1.Resource.Header.Name)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(1))
		Expect(versions.Resources[0].Id).To(Equal(res2v1.Resource.Header.Id))
	})

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
