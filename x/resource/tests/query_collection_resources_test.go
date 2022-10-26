package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"

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

			if tc.valid {
				resources := queryResponse.Resources
				expectedResources := tc.response.Resources
				require.Nil(t, err)
				require.Equal(t, len(expectedResources), len(resources))
				for i, r := range resources {
					r.Created = expectedResources[i].Created
					require.Equal(t, expectedResources[i], r)
				}
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
