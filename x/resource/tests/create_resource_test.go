package tests

import (
	"crypto/sha256"
	"strings"

	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	"github.com/google/uuid"

	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ExpectPayloadToMatchResource(payload *resourcetypes.MsgCreateResourcePayload, resource *resourcetypes.ResourceWithMetadata) {
	// Provided header
	Expect(payload.Id).To(Equal(resource.Metadata.Id))
	Expect(payload.CollectionId).To(Equal(resource.Metadata.CollectionId))
	Expect(payload.Name).To(Equal(resource.Metadata.Name))
	Expect(payload.ResourceType).To(Equal(resource.Metadata.ResourceType))

	// Generated header
	hash := sha256.Sum256(payload.Data)
	Expect(resource.Metadata.Checksum).To(Equal(hash[:]))

	// Provided data
	Expect(payload.Data).To(Equal(resource.Resource.Data))
}

var _ = Describe("Create Resource Tests", func() {
	var setup TestSetup
	var alice didsetup.CreatedDidDocInfo

	BeforeEach(func() {
		setup = Setup()
		alice = setup.CreateSimpleDid()
	})

	Describe("Simple resource", func() {
		var msg *resourcetypes.MsgCreateResourcePayload

		BeforeEach(func() {
			msg = &resourcetypes.MsgCreateResourcePayload{
				CollectionId: alice.CollectionId,
				Id:           uuid.NewString(),
				Name:         "Test Resource Name",
				ResourceType: CLSchemaType,
				Data:         []byte(SchemaData),
			}
		})

		It("Can be created with DIDDoc controller signature", func() {
			_, err := setup.CreateResource(msg, []didsetup.SignInput{alice.SignInput})
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryResource(alice.CollectionId, msg.Id)
			Expect(err).To(BeNil())

			ExpectPayloadToMatchResource(msg, created.Resource)
		})

		It("Can't be created without DIDDoc controller signatures", func() {
			_, err := setup.CreateResource(msg, []didsetup.SignInput{})
			Expect(err.Error()).To(ContainSubstring("signature is required but not found"))
		})

		It("Can't be created with invalid collection id", func() {
			msg.CollectionId = didsetup.GenerateDID(didsetup.Base58_16bytes)

			_, err := setup.CreateResource(msg, []didsetup.SignInput{alice.SignInput})
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Describe("New version", func() {
		var existingResource *resourcetypes.MsgCreateResourceResponse

		BeforeEach(func() {
			existingResource = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Test Resource Name", CLSchemaType, []didsetup.SignInput{alice.SignInput})
		})

		It("Is linked to the previous one when name matches", func() {
			msg := resourcetypes.MsgCreateResourcePayload{
				CollectionId: alice.CollectionId,
				Id:           uuid.NewString(),
				Name:         existingResource.Resource.Name,
				ResourceType: CLSchemaType,
				Data:         []byte(SchemaData),
			}

			_, err := setup.CreateResource(&msg, []didsetup.SignInput{alice.SignInput})
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryResource(alice.CollectionId, msg.Id)
			Expect(err).To(BeNil())

			ExpectPayloadToMatchResource(&msg, created.Resource)
			Expect(created.Resource.Metadata.PreviousVersionId).To(Equal(existingResource.Resource.Id))
		})
	})

	Describe("Resource for deactivated DID", func() {
		var msg *resourcetypes.MsgCreateResourcePayload

		BeforeEach(func() {
			msg = &resourcetypes.MsgCreateResourcePayload{
				CollectionId: alice.CollectionId,
				Id:           uuid.NewString(),
				Name:         "Test Resource Name",
				ResourceType: CLSchemaType,
				Data:         []byte(SchemaData),
			}
		})

		When("DIDDoc is deactivated", func() {
			It("Should fail with error", func() {
				// Deactivate DID
				DeactivateMsg := &didtypes.MsgDeactivateDidDocPayload{
					Id:        alice.Did,
					VersionId: uuid.NewString(),
				}

				signatures := []didsetup.SignInput{alice.DidDocInfo.SignInput}

				res, err := setup.DeactivateDidDoc(DeactivateMsg, signatures)
				Expect(err).To(BeNil())
				Expect(res.Value.Metadata.Deactivated).To(BeTrue())

				// Create resource
				_, err = setup.CreateResource(msg, []didsetup.SignInput{alice.SignInput})
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring(alice.Did + ": DID Doc already deactivated"))
			})
		})
	})

	Describe("UUID with capital letters", func() {
		It("Should work even for UUID with capital letters", func() {
			msg := resourcetypes.MsgCreateResourcePayload{
				CollectionId: alice.CollectionId,
				Id:           UUIDString,
				Name:         "Resource with capital letters in UUID",
				ResourceType: CLSchemaType,
				Data:         []byte(SchemaData),
			}

			_, err := setup.CreateResource(&msg, []didsetup.SignInput{alice.SignInput})
			Expect(err).To(BeNil())

			// check for the same UUID
			created, err := setup.QueryResource(alice.CollectionId, UUIDString)
			Expect(err).To(BeNil())

			Expect(created.Resource.Metadata.Id).To(Equal(strings.ToLower(UUIDString)))

			// check for already normalized UUID
			created, err = setup.QueryResource(
				didutils.NormalizeId(alice.CollectionId),
				didutils.NormalizeId(UUIDString))
			Expect(err).To(BeNil())

			Expect(created.Resource.Metadata.Id).To(Equal(strings.ToLower(UUIDString)))
		})
	})
})
