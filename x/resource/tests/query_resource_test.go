package tests

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryGetResource", func() {
	Describe("Validate", func() {
		var setup TestSetup
		keys := GenerateTestKeys()
		existingResource := ExistingResource()
		BeforeEach(func() {
			setup = Setup()
			didDoc := setup.CreateDid(keys[ExistingDIDKey].Public, ExistingDID)
			_, err := setup.SendCreateDid(didDoc, map[string]ed25519.PrivateKey{ExistingDIDKey: keys[ExistingDIDKey].Private})
			Expect(err).To(BeNil())
			payload := GenerateCreateResourcePayload(ExistingResource())
			_, err = setup.SendCreateResource(payload, map[string]ed25519.PrivateKey{ExistingDIDKey: keys[ExistingDIDKey].Private})
			Expect(err).To(BeNil())
		})
		DescribeTable("Validate QueryGetResourceRequest",
			func(
				valid bool,
				msg *resourcetypes.QueryGetResourceRequest,
				response *resourcetypes.QueryGetResourceResponse,
				errMsg string,
			) {
				queryResponse, err := setup.ResourceQueryServer.Resource(setup.StdCtx, msg)

				if valid {
					resource := queryResponse.Resource
					Expect(err).To(BeNil())
					Expect(response.Resource.Header.CollectionId).To(Equal(resource.Header.CollectionId))
					Expect(response.Resource.Header.Id).To(Equal(resource.Header.Id))
					Expect(response.Resource.Header.MediaType).To(Equal(resource.Header.MediaType))
					Expect(response.Resource.Header.ResourceType).To(Equal(resource.Header.ResourceType))
					Expect(response.Resource.Data).To(Equal(resource.Data))
					Expect(response.Resource.Header.Name).To(Equal(resource.Header.Name))
					checksum := sha256.Sum256(response.Resource.Data)
					Expect(checksum[:]).To(Equal(resource.Header.Checksum))
					Expect(response.Resource.Header.PreviousVersionId).To(Equal(resource.Header.PreviousVersionId))
					Expect(response.Resource.Header.NextVersionId).To(Equal(resource.Header.NextVersionId))
				} else {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(errMsg))
				}
			},
			Entry("Valid: Works",
				true,
				&resourcetypes.QueryGetResourceRequest{
					CollectionId: ExistingDIDIdentifier,
					Id:           existingResource.Header.Id,
				},
				&resourcetypes.QueryGetResourceResponse{
					Resource: &existingResource,
				},
				"",
			),
			Entry("Invalid: Resource not found",
				false,
				&resourcetypes.QueryGetResourceRequest{
					CollectionId: ExistingDIDIdentifier,
					Id:           AnotherResourceId,
				},
				nil,
				fmt.Errorf("resource %s:%s: not found", ExistingDIDIdentifier, AnotherResourceId).Error(),
			),
			Entry("Invalid: DIDDoc not found",
				false,
				&resourcetypes.QueryGetResourceRequest{
					CollectionId: NotFoundDIDIdentifier,
					Id:           existingResource.Header.Id,
				},
				nil,
				fmt.Errorf("did:cheqd:test:%s: DID Doc not found", NotFoundDIDIdentifier).Error(),
			),
		)
	})
})
