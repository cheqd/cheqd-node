package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"

	cheqdsetup "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"
	cheqdutils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Collection Resources", func() {
	var setup TestSetup
	var alice cheqdsetup.CreatedDidInfo
	var uuidDID cheqdsetup.CreatedDidInfo

	var res1v1 *types.MsgCreateResourceResponse
	var res1v2 *types.MsgCreateResourceResponse
	var res2v1 *types.MsgCreateResourceResponse
	var resUUID *types.MsgCreateResourceResponse

	BeforeEach(func() {
		setup = Setup()

		alice = setup.CreateSimpleDid()
		uuidDID = setup.CreateUUIDDid(UUIDString)

		res1v1 = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 1", CLSchemaType, []cheqdsetup.SignInput{alice.SignInput})
		res1v2 = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 1", CLSchemaType, []cheqdsetup.SignInput{alice.SignInput})
		res2v1 = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 2", CLSchemaType, []cheqdsetup.SignInput{alice.SignInput})
		resUUID = setup.CreateSimpleResource(uuidDID.CollectionId, SchemaData, "Resource UUID", CLSchemaType, []cheqdsetup.SignInput{uuidDID.SignInput})
	})

	It("Should return all 3 headers", func() {
		versions, err := setup.CollectionResources(alice.CollectionId)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(3))

		ids := []string{versions.Resources[0].Id, versions.Resources[1].Id, versions.Resources[2].Id}

		Expect(ids).To(ContainElement(res1v1.Resource.Header.Id))
		Expect(ids).To(ContainElement(res1v2.Resource.Header.Id))
		Expect(ids).To(ContainElement(res2v1.Resource.Header.Id))
	})

	It("Should work with capital letters in UUID", func() {
		// Here we are asking for non-normalized UUID
		versions, err := setup.CollectionResources(uuidDID.CollectionId)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(1))

		ids := []string{versions.Resources[0].Id}

		Expect(ids).To(ContainElement(resUUID.Resource.Header.Id))
	})

	It("Should work with capital letters in UUID. Ask with already normalized collectionId", func() {
		// Here we are asking for normalized UUID but it was written with capital letters
		normalizedId := cheqdutils.NormalizeId(uuidDID.CollectionId)
		versions, err := setup.CollectionResources(normalizedId)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(1))

		ids := []string{versions.Resources[0].Id}

		Expect(ids).To(ContainElement(resUUID.Resource.Header.Id))
	})
})
