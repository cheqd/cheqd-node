package tests_test

import (
	"crypto/ed25519"
	"fmt"

	resourcetests "github.com/cheqd/cheqd-node/x/resource/tests"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryCollectionResources", func() {
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
		DescribeTable("Validate QueryCollectionResources",
			func(
				valid bool,
				msg *resourcetypes.QueryGetCollectionResourcesRequest,
				response *resourcetypes.QueryGetCollectionResourcesResponse,
				errMsg string,
			) {
				queryResponse, err := setup.QueryServer.CollectionResources(sdk.WrapSDKContext(setup.Ctx), msg)

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
					CollectionId: resourcetests.ExistingDIDIdentifier,
				},
				&resourcetypes.QueryGetCollectionResourcesResponse{
					Resources: []*resourcetypes.ResourceHeader{resourcetests.ExistingResource().Header},
				},
				"",
			),
			Entry("Invalid: DID Doc is not found",
				false,
				&resourcetypes.QueryGetCollectionResourcesRequest{
					CollectionId: resourcetests.NotFoundDIDIdentifier,
				},
				nil,
				fmt.Errorf("did:cheqd:test:%s: DID Doc not found", resourcetests.NotFoundDIDIdentifier).Error(),
			),
		)
	})
})
