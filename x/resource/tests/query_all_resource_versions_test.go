package tests_test

import (
	"crypto/ed25519"
	"fmt"

	resourcetests "github.com/cheqd/cheqd-node/x/resource/tests"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryAllResourceVersions", func() {
	Describe("Validate", func() {
		var setup resourcetests.TestSetup
		keys := resourcetests.GenerateTestKeys()
		BeforeEach(func() {
			setup = resourcetests.Setup()
			didDoc := setup.CreateDid(keys[resourcetests.ExistingDIDKey].PublicKey, resourcetests.ExistingDID)
			_, err := setup.SendCreateDid(didDoc, map[string]ed25519.PrivateKey{resourcetests.ExistingDIDKey: keys[resourcetests.ExistingDIDKey].PrivateKey})
			Expect(err).To(BeNil())
			payload := resourcetests.GenerateCreateResourcePayload(resourcetests.ExistingResource())
			_, err = setup.SendCreateResource(payload, map[string]ed25519.PrivateKey{resourcetests.ExistingDIDKey: keys[resourcetests.ExistingDIDKey].PrivateKey})
			Expect(err).To(BeNil())
		})
		DescribeTable("Validate QueryGetAllResourceVersionsRequest",
			func(
				valid bool,
				signerKeys map[string]ed25519.PrivateKey,
				msg *types.QueryGetAllResourceVersionsRequest,
				response *types.QueryGetAllResourceVersionsResponse,
				errMsg string,
			) {
				existingResource := resourcetests.ExistingResource()

				payload := resourcetests.GenerateCreateResourcePayload(existingResource)
				payload.Id = resourcetests.ResourceId

				nextVersionResource, err := setup.SendCreateResource(payload, signerKeys)
				Expect(err).To(BeNil())
				Expect(nextVersionResource).ToNot(Equal(existingResource))

				payload = resourcetests.GenerateCreateResourcePayload(existingResource)
				payload.Id = resourcetests.AnotherResourceId
				payload.Name = "AnotherResourceVersion"
				differentResource, err := setup.SendCreateResource(payload, signerKeys)
				Expect(err).To(BeNil())
				Expect(differentResource).ToNot(Equal(existingResource))
				Expect(differentResource).ToNot(Equal(nextVersionResource))

				queryResponse, err := setup.QueryServer.AllResourceVersions(sdk.WrapSDKContext(setup.Ctx), msg)

				if valid {
					resources := queryResponse.Resources
					existingResource.Header.NextVersionId = nextVersionResource.Header.Id
					expectedResources := map[string]types.Resource{
						existingResource.Header.Id:    existingResource,
						nextVersionResource.Header.Id: *nextVersionResource,
					}
					Expect(err).To(BeNil())
					Expect(len(resources)).To(Equal(len(expectedResources)))
					for _, r := range resources {
						r.Created = expectedResources[r.Id].Header.Created
						Expect(r).To(Equal(expectedResources[r.Id].Header))
					}
				} else {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal(errMsg))
				}
			},
			Entry("Valid: should return all resources",
				true,
				map[string]ed25519.PrivateKey{resourcetests.ExistingDIDKey: keys[resourcetests.ExistingDIDKey].PrivateKey},
				&types.QueryGetAllResourceVersionsRequest{
					CollectionId: resourcetests.ExistingDIDIdentifier,
					Name:         resourcetests.ExistingResource().Header.Name,
				},
				&types.QueryGetAllResourceVersionsResponse{
					Resources: []*types.ResourceHeader{
						resourcetests.ExistingResource().Header,
					},
				},
				"",
			),
			Entry("Invalid: should return an error if the collection id is invalid",
				false,
				map[string]ed25519.PrivateKey{resourcetests.ExistingDIDKey: keys[resourcetests.ExistingDIDKey].PrivateKey},
				&types.QueryGetAllResourceVersionsRequest{
					CollectionId: resourcetests.NotFoundDIDIdentifier,
					Name:         resourcetests.ExistingResource().Header.Name,
				},
				nil,
				fmt.Errorf("did:cheqd:test:%s: DID Doc not found", resourcetests.NotFoundDIDIdentifier).Error(),
			),
		)
	})
})
