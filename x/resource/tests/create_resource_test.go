package tests

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"

	cheqdutils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateResource", func() {
	Describe("Validate", func() {
		var setup TestSetup
		var err error
		keys := GenerateTestKeys()
		BeforeEach(func() {
			setup = Setup()
			didDoc := setup.CreateDid(keys[ExistingDIDKey].Public, ExistingDID)
			_, err = setup.SendCreateDid(didDoc, map[string]ed25519.PrivateKey{ExistingDIDKey: keys[ExistingDIDKey].Private})
			Expect(err).To(BeNil())
			resourcePayload := GenerateCreateResourcePayload(ExistingResource())
			_, err = setup.SendCreateResource(resourcePayload, map[string]ed25519.PrivateKey{ExistingDIDKey: keys[ExistingDIDKey].Private})
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
				map[string]ed25519.PrivateKey{ExistingDIDKey: keys[ExistingDIDKey].Private},
				&resourcetypes.MsgCreateResourcePayload{
					CollectionId: ExistingDIDIdentifier,
					Id:           ResourceId,
					Name:         "Test Resource Name",
					ResourceType: CLSchemaType,
					Data:         []byte(SchemaData),
				},
				JsonResourceType,
				"",
				"",
			),
			Entry("Valid: Add new resource version",
				true,
				map[string]ed25519.PrivateKey{ExistingDIDKey: keys[ExistingDIDKey].Private},
				&resourcetypes.MsgCreateResourcePayload{
					CollectionId: ExistingResource().Header.CollectionId,
					Id:           ResourceId,
					Name:         ExistingResource().Header.Name,
					ResourceType: ExistingResource().Header.ResourceType,
					Data:         ExistingResource().Data,
				},
				ExistingResource().Header.MediaType,
				ExistingResource().Header.Id,
				"",
			),
			Entry("Invalid: No signature",
				false,
				map[string]ed25519.PrivateKey{},
				&resourcetypes.MsgCreateResourcePayload{
					CollectionId: ExistingDIDIdentifier,
					Id:           ResourceId,
					Name:         "Test Resource Name",
					ResourceType: CLSchemaType,
					Data:         []byte(SchemaData),
				},
				JsonResourceType,
				"",
				fmt.Errorf("signer: %s: signature is required but not found", ExistingDID).Error(),
			),
			Entry("Invalid: Resource Id is not an acceptable format",
				false,
				map[string]ed25519.PrivateKey{},
				&resourcetypes.MsgCreateResourcePayload{
					CollectionId: ExistingDIDIdentifier,
					Id:           IncorrectResourceId,
					Name:         "Test Resource Name",
					ResourceType: CLSchemaType,
					Data:         []byte(SchemaData),
				},
				JsonResourceType,
				"",
				fmt.Errorf("signer: %s: signature is required but not found", ExistingDID).Error(),
			),
			Entry("Invalid: DidDoc not found",
				false,
				map[string]ed25519.PrivateKey{},
				&resourcetypes.MsgCreateResourcePayload{
					CollectionId: NotFoundDIDIdentifier,
					Id:           IncorrectResourceId,
					Name:         "Test Resource Name",
					ResourceType: CLSchemaType,
					Data:         []byte(SchemaData),
				},
				JsonResourceType,
				"",
				fmt.Errorf("did:cheqd:test:%s: not found", NotFoundDIDIdentifier).Error(),
			),
		)
	})
})
