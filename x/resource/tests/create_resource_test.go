package tests_test

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"

	cheqdutils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	resourcetests "github.com/cheqd/cheqd-node/x/resource/tests"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateResource", func() {
	Describe("Validate", func() {
		var setup resourcetests.TestSetup
		var err error
		keys := resourcetests.GenerateTestKeys()
		BeforeEach(func() {
			setup = resourcetests.Setup()
			didDoc := setup.CreateDid(keys[resourcetests.ExistingDIDKey].PublicKey, resourcetests.ExistingDID)
			_, err = setup.SendCreateDid(didDoc, map[string]ed25519.PrivateKey{resourcetests.ExistingDIDKey: keys[resourcetests.ExistingDIDKey].PrivateKey})
			Expect(err).To(BeNil())
			resourcePayload := resourcetests.GenerateCreateResourcePayload(resourcetests.ExistingResource())
			_, err = setup.SendCreateResource(resourcePayload, map[string]ed25519.PrivateKey{resourcetests.ExistingDIDKey: keys[resourcetests.ExistingDIDKey].PrivateKey})
			Expect(err).To(BeNil())
		})
		DescribeTable("Validate MsgCreateResource",
			func(
				valid bool,
				signerKeys map[string]ed25519.PrivateKey,
				msg *resourcetypes.MsgCreateResourcePayload,
				mediaType string,
				previousVersionId string,
				errMsg string,
			) {
				resource, err := setup.SendCreateResource(msg, signerKeys)
				if valid {
					Expect(err).To(BeNil())

					did := cheqdutils.JoinDID("cheqd", "test", resource.Header.CollectionId)
					didStateValue, err := setup.Keeper.GetDid(&setup.Ctx, did)
					Expect(err).To(BeNil())
					Expect(didStateValue.Metadata.Resources).Should(ContainElement(resource.Header.Id))

					Expect(resource.Header.CollectionId).To(Equal(msg.CollectionId))
					Expect(resource.Header.Id).To(Equal(msg.Id))
					Expect(resource.Header.Name).To(Equal(msg.Name))
					Expect(resource.Header.ResourceType).To(Equal(msg.ResourceType))
					Expect(resource.Header.MediaType).To(Equal(mediaType))
					Expect(resource.Header.PreviousVersionId).To(Equal(previousVersionId))
					expectedChecksum := sha256.Sum256(msg.Data)
					Expect(resource.Header.Checksum).To(Equal(expectedChecksum[:]))
				} else {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal(errMsg))
				}
			},
			Entry("Valid: Works",
				true,
				map[string]ed25519.PrivateKey{resourcetests.ExistingDIDKey: keys[resourcetests.ExistingDIDKey].PrivateKey},
				&resourcetypes.MsgCreateResourcePayload{
					CollectionId: resourcetests.ExistingDIDIdentifier,
					Id:           resourcetests.ResourceId,
					Name:         "Test Resource Name",
					ResourceType: resourcetests.CLSchemaType,
					Data:         []byte(resourcetests.SchemaData),
				},
				resourcetests.JsonResourceType,
				"",
				"",
			),
			Entry("Valid: Add new resource version",
				true,
				map[string]ed25519.PrivateKey{resourcetests.ExistingDIDKey: keys[resourcetests.ExistingDIDKey].PrivateKey},
				&resourcetypes.MsgCreateResourcePayload{
					CollectionId: resourcetests.ExistingResource().Header.CollectionId,
					Id:           resourcetests.ResourceId,
					Name:         resourcetests.ExistingResource().Header.Name,
					ResourceType: resourcetests.ExistingResource().Header.ResourceType,
					Data:         resourcetests.ExistingResource().Data,
				},
				resourcetests.ExistingResource().Header.MediaType,
				resourcetests.ExistingResource().Header.Id,
				"",
			),
			Entry("Invalid: No signature",
				false,
				map[string]ed25519.PrivateKey{},
				&resourcetypes.MsgCreateResourcePayload{
					CollectionId: resourcetests.ExistingDIDIdentifier,
					Id:           resourcetests.ResourceId,
					Name:         "Test Resource Name",
					ResourceType: resourcetests.CLSchemaType,
					Data:         []byte(resourcetests.SchemaData),
				},
				resourcetests.JsonResourceType,
				"",
				fmt.Errorf("signer: %s: signature is required but not found", resourcetests.ExistingDID).Error(),
			),
			Entry("Invalid: Resource Id is not an acceptable format",
				false,
				map[string]ed25519.PrivateKey{},
				&resourcetypes.MsgCreateResourcePayload{
					CollectionId: resourcetests.ExistingDIDIdentifier,
					Id:           resourcetests.IncorrectResourceId,
					Name:         "Test Resource Name",
					ResourceType: resourcetests.CLSchemaType,
					Data:         []byte(resourcetests.SchemaData),
				},
				resourcetests.JsonResourceType,
				"",
				fmt.Errorf("signer: %s: signature is required but not found", resourcetests.ExistingDID).Error(),
			),
			Entry("Invalid: DidDoc not found",
				false,
				map[string]ed25519.PrivateKey{},
				&resourcetypes.MsgCreateResourcePayload{
					CollectionId: resourcetests.NotFoundDIDIdentifier,
					Id:           resourcetests.IncorrectResourceId,
					Name:         "Test Resource Name",
					ResourceType: resourcetests.CLSchemaType,
					Data:         []byte(resourcetests.SchemaData),
				},
				resourcetests.JsonResourceType,
				"",
				fmt.Errorf("did:cheqd:test:%s: not found", resourcetests.NotFoundDIDIdentifier).Error(),
			),
		)
	})
})
