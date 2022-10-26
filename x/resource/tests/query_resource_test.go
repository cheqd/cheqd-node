package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	"github.com/google/uuid"

	cheqdsetup "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"
	"github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Collection Resources", func() {
	var setup TestSetup
	var alice cheqdsetup.CreatedDidInfo
	var resource *types.MsgCreateResourceResponse

	BeforeEach(func() {
		setup = Setup()
		alice = setup.CreateSimpleDid()
		resource = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 1", CLSchemaType, []cheqdsetup.SignInput{alice.SignInput})
	})

	It("Works", func() {
		versions, err := setup.QueryResource(alice.CollectionId, resource.Resource.Header.Id)
		Expect(err).To(BeNil())
		Expect(versions.Resource.Header.Id).To(Equal(resource.Resource.Header.Id))
	})

	It("Returns error if resource does not exist", func() {
		nonExistingResource := uuid.NewString()

		_, err := setup.QueryResource(alice.CollectionId, nonExistingResource)
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("Returns error if collection does not exist", func() {
		nonExistingCollection := cheqdsetup.GenerateDID(cheqdsetup.Base58_16chars)

		_, err := setup.QueryResource(nonExistingCollection, resource.Resource.Header.Id)
		Expect(err.Error()).To(ContainSubstring("DID Doc not found"))
	})
})