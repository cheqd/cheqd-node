package tests

import (
	"crypto/ed25519"
	"fmt"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryCollectionResources", func() {
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
		DescribeTable("Validate QueryCollectionResources",
			func(
				valid bool,
				msg *resourcetypes.QueryGetCollectionResourcesRequest,
				response *resourcetypes.QueryGetCollectionResourcesResponse,
				errMsg string,
			) {
				queryResponse, err := setup.ResourceQueryServer.CollectionResources(setup.StdCtx, msg)

				if valid {
					resources := queryResponse.Resources
					expectedResources := response.Resources
					Expect(err).To(BeNil())
					Expect(len(expectedResources)).To(Equal(len(resources)))
					for i, r := range resources {
						r.Created = expectedResources[i].Created
						Expect(expectedResources[i]).To(Equal(r))
					}
				} else {
					Expect(err).To(HaveOccurred())
					Expect(errMsg).To(Equal(err.Error()))
				}
			},
			Entry("Valid: Works",
				true,
				&resourcetypes.QueryGetCollectionResourcesRequest{
					CollectionId: ExistingDIDIdentifier,
				},
				&resourcetypes.QueryGetCollectionResourcesResponse{
					Resources: []*resourcetypes.ResourceHeader{ExistingResource().Header},
				},
				"",
			),
			Entry("Invalid: DID Doc is not found",
				false,
				&resourcetypes.QueryGetCollectionResourcesRequest{
					CollectionId: NotFoundDIDIdentifier,
				},
				nil,
				fmt.Errorf("did:cheqd:test:%s: DID Doc not found", NotFoundDIDIdentifier).Error(),
			),
		)
	})
})
