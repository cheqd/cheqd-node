package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"

	cheqdsetup "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"
	"github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Collection Resources", func() {
	var setup TestSetup
	var alice cheqdsetup.CreatedDidInfo

	var res1v1 *types.MsgCreateResourceResponse
	var res1v2 *types.MsgCreateResourceResponse
	var res2v1 *types.MsgCreateResourceResponse

	BeforeEach(func() {
		setup = Setup()

		alice = setup.CreateSimpleDid()

		res1v1 = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 1", CLSchemaType, []cheqdsetup.SignInput{alice.SignInput})
		res1v2 = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 1", CLSchemaType, []cheqdsetup.SignInput{alice.SignInput})
		res2v1 = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 2", CLSchemaType, []cheqdsetup.SignInput{alice.SignInput})
	})

	It("Should return all 3 headerrs", func() {
		versions, err := setup.CollectionResources(alice.CollectionId)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(3))

		ids := []string{versions.Resources[0].Id, versions.Resources[1].Id, versions.Resources[2].Id}

		Expect(ids).To(ContainElement(res1v1.Resource.Header.Id))
		Expect(ids).To(ContainElement(res1v2.Resource.Header.Id))
		Expect(ids).To(ContainElement(res2v1.Resource.Header.Id))
	})
})
