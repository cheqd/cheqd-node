package tests

import (
	"crypto/ed25519"
	"fmt"

	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryAllResourceVersions", func() {
	Describe("Validate", func() {
		var setup TestSetup
		keys := GenerateTestKeys()
		BeforeEach(func() {
			setup = Setup()
			didDoc := setup.CreateDid(keys[ExistingDIDKey].Public, ExistingDID)
			_, err := setup.SendCreateDid(didDoc, map[string]ed25519.PrivateKey{ExistingDIDKey: keys[ExistingDIDKey].Private})
			Expect(err).To(BeNil())
			payload := GenerateCreateResourcePayload(ExistingResource())
			_, err = setup.SendCreateResource(payload, map[string]ed25519.PrivateKey{ExistingDIDKey: keys[ExistingDIDKey].Private})
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
				existingResource := ExistingResource()

				payload := GenerateCreateResourcePayload(existingResource)
				payload.Id = ResourceId

				nextVersionResource, err := setup.SendCreateResource(payload, signerKeys)
				Expect(err).To(BeNil())
				Expect(nextVersionResource).ToNot(Equal(existingResource))

				payload = GenerateCreateResourcePayload(existingResource)
				payload.Id = AnotherResourceId
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
				map[string]ed25519.PrivateKey{ExistingDIDKey: keys[ExistingDIDKey].Private},
				&types.QueryGetAllResourceVersionsRequest{
					CollectionId: ExistingDIDIdentifier,
					Name:         ExistingResource().Header.Name,
				},
				&types.QueryGetAllResourceVersionsResponse{
					Resources: []*types.ResourceHeader{
						ExistingResource().Header,
					},
				},
				"",
			),
			Entry("Invalid: should return an error if the collection id is invalid",
				false,
				map[string]ed25519.PrivateKey{ExistingDIDKey: keys[ExistingDIDKey].Private},
				&types.QueryGetAllResourceVersionsRequest{
					CollectionId: NotFoundDIDIdentifier,
					Name:         ExistingResource().Header.Name,
				},
				nil,
				fmt.Errorf("did:cheqd:test:%s: DID Doc not found", NotFoundDIDIdentifier).Error(),
			),
		)
	})
})
