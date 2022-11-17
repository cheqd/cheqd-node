package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	"github.com/google/uuid"

	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Resource", func() {
	var setup TestSetup
	var alice didsetup.CreatedDidDocInfo
	var resource *types.MsgCreateResourceResponse

	BeforeEach(func() {
		setup = Setup()
		alice = setup.CreateSimpleDid()
		resource = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 1", CLSchemaType, []didsetup.SignInput{alice.SignInput})
	})

	It("Works", func() {
		resp, err := setup.QueryResource(alice.CollectionId, resource.Resource.Id)
		Expect(err).To(BeNil())
		Expect(resp.Resource.Metadata.Id).To(Equal(resource.Resource.Id))
	})

	It("Returns error if resource does not exist", func() {
		nonExistingResource := uuid.NewString()

		_, err := setup.QueryResource(alice.CollectionId, nonExistingResource)
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("Returns error if collection does not exist", func() {
		nonExistingCollection := didsetup.GenerateDID(didsetup.Base58_16bytes)

		_, err := setup.QueryResource(nonExistingCollection, resource.Resource.Id)
		Expect(err.Error()).To(ContainSubstring("DID Doc not found"))
	})
})
