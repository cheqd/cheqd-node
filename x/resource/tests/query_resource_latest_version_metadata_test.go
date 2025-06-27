package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	"github.com/google/uuid"

	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query LatestResourceVersion Metadata", func() {
	var setup TestSetup
	var alice didsetup.CreatedDidDocInfo
	var resource *types.MsgCreateResourceResponse

	BeforeEach(func() {
		setup = Setup()
		alice = setup.CreateSimpleDid()
		_ = setup.CreateSimpleResource(alice.CollectionID, SchemaData, "Resource 1%", CLSchemaType, []didsetup.SignInput{alice.SignInput})
		// update resource providing same name and resourceType
		resource = setup.CreateSimpleResource(alice.CollectionID, SchemaData, "Resource 1%", CLSchemaType, []didsetup.SignInput{alice.SignInput})
	})

	It("Works", func() {
		metadata, err := setup.QueryLatestResourceVersionMetadata(alice.CollectionID, "Resource 1%", CLSchemaType)
		Expect(err).To(BeNil())
		Expect(metadata.Resource).To(Equal(resource.Resource))
	})

	It("Works with different composed forms", func() {
		resource := setup.CreateSimpleResource(alice.CollectionID, SchemaData, "Résourcℯ", CLSchemaType, []didsetup.SignInput{alice.SignInput})
		metadata, err := setup.QueryLatestResourceVersionMetadata(alice.CollectionID, "Résourcℯ", CLSchemaType)
		Expect(err).To(BeNil())
		Expect(metadata.Resource).To(Equal(resource.Resource))

		_, err = setup.QueryLatestResourceVersionMetadata(alice.CollectionID, "Resource", CLSchemaType)
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("Returns error (not found) if resource does not exist", func() {
		nonExistingResource := uuid.NewString()

		_, err := setup.QueryLatestResourceVersionMetadata(alice.CollectionID, nonExistingResource, CLSchemaType)
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("Returns error (not found) if resource index keys have trailing spaces", func() {
		_, err := setup.QueryLatestResourceVersionMetadata(alice.CollectionID, "  Resource 1%  ", CLSchemaType)
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("Returns error (not found) if resource index keys have different composed forms", func() {
		_, err := setup.QueryLatestResourceVersionMetadata(alice.CollectionID, "Résourcℯ 1%", CLSchemaType)
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("Returns error (not found) if resource index keys have no spaces", func() {
		_, err := setup.QueryLatestResourceVersionMetadata(alice.CollectionID, "Resource1%", CLSchemaType)
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("Returns error (did doc not found) if collection does not exist", func() {
		nonExistingCollection := didsetup.GenerateDID(didsetup.Base58_16bytes)

		_, err := setup.QueryLatestResourceVersionMetadata(nonExistingCollection, "Resource 1%", CLSchemaType)
		Expect(err.Error()).To(ContainSubstring("DID Doc not found"))
	})
})
