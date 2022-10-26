package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"

	cheqdsetup "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"
	"github.com/cheqd/cheqd-node/x/resource/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query All Resource Versions", func() {
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

	It("Should return 2 versions for resource 1", func() {
		versions, err := setup.AllResourceVersions(alice.CollectionId, res1v1.Resource.Header.Name, CLSchemaType)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(2))

		ids := []string{versions.Resources[0].Id, versions.Resources[1].Id}

		Expect(ids).To(ContainElement(res1v1.Resource.Header.Id))
		Expect(ids).To(ContainElement(res1v2.Resource.Header.Id))
	})

	It("Should return 1 version for resource 2", func() {
		versions, err := setup.AllResourceVersions(alice.CollectionId, res2v1.Resource.Header.Name, CLSchemaType)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(1))
		Expect(versions.Resources[0].Id).To(Equal(res2v1.Resource.Header.Id))
	})

	It("Should return 0 versions for non-existing resource", func() {
		versions, err := setup.AllResourceVersions(alice.CollectionId, "non-existing", CLSchemaType)
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(0))
	})

	It("Should return 0 versions for existing resource but with unexpected schema", func() {
		versions, err := setup.AllResourceVersions(alice.CollectionId, res1v1.Resource.Header.Name, "non-existing-schema-type")
		Expect(err).To(BeNil())
		Expect(versions.Resources).To(HaveLen(0))
	})
})