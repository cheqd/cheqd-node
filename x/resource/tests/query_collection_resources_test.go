package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"

	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Collection Resources", func() {
	var setup TestSetup
	var alice didsetup.CreatedDidDocInfo
	var uuidDID didsetup.CreatedDidDocInfo

	var res1v1 *types.MsgCreateResourceResponse
	var res1v2 *types.MsgCreateResourceResponse
	var res2v1 *types.MsgCreateResourceResponse
	var resUUID *types.MsgCreateResourceResponse

	BeforeEach(func() {
		setup = Setup()

		alice = setup.CreateSimpleDid()
		uuidDID = setup.CreateCustomDidDoc(setup.BuildDidDocWithCustomId(UUIDString))

		res1v1 = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 1", CLSchemaType, []didsetup.SignInput{alice.SignInput})
		res1v2 = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 1", CLSchemaType, []didsetup.SignInput{alice.SignInput})
		res2v1 = setup.CreateSimpleResource(alice.CollectionId, SchemaData, "Resource 2", CLSchemaType, []didsetup.SignInput{alice.SignInput})
		resUUID = setup.CreateSimpleResource(uuidDID.CollectionId, SchemaData, "Resource UUID", CLSchemaType, []didsetup.SignInput{uuidDID.SignInput})
	})

	It("Should return all 3 headers", func() {
		versions, err := setup.CollectionResources(alice.CollectionId)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(3))

		ids := []string{versions.Resources[0].Id, versions.Resources[1].Id, versions.Resources[2].Id}

		Expect(ids).To(ContainElement(res1v1.Resource.Id))
		Expect(ids).To(ContainElement(res1v2.Resource.Id))
		Expect(ids).To(ContainElement(res2v1.Resource.Id))
	})

	It("Should work with capital letters in UUID", func() {
		// Here we are asking for non-normalized UUID
		versions, err := setup.CollectionResources(uuidDID.CollectionId)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(1))

		ids := []string{versions.Resources[0].Id}

		Expect(ids).To(ContainElement(resUUID.Resource.Id))
	})

	It("Should work with capital letters in UUID. Ask with already normalized collectionId", func() {
		// Here we are asking for normalized UUID but it was written with capital letters
		normalizedId := didutils.NormalizeId(uuidDID.CollectionId)
		versions, err := setup.CollectionResources(normalizedId)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(1))

		ids := []string{versions.Resources[0].Id}

		Expect(ids).To(ContainElement(resUUID.Resource.Id))
	})
})
