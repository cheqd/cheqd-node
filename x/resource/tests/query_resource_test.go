package tests

import (
	. "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	"github.com/google/uuid"

	// "crypto/sha256"
	"crypto/sha256"
	"fmt"
	"testing"

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

			if tc.valid {
				resource := queryResponse.Resource
				require.Nil(t, err)
				require.Equal(t, tc.response.Resource.Header.CollectionId, resource.Header.CollectionId)
				require.Equal(t, tc.response.Resource.Header.Id, resource.Header.Id)
				require.Equal(t, tc.response.Resource.Header.MediaType, resource.Header.MediaType)
				require.Equal(t, tc.response.Resource.Header.ResourceType, resource.Header.ResourceType)
				require.Equal(t, tc.response.Resource.Data, resource.Data)
				require.Equal(t, tc.response.Resource.Header.Name, resource.Header.Name)
				require.Equal(t, sha256.New().Sum(tc.response.Resource.Data), resource.Header.Checksum)
				require.Equal(t, tc.response.Resource.Header.PreviousVersionId, resource.Header.PreviousVersionId)
				require.Equal(t, tc.response.Resource.Header.NextVersionId, resource.Header.NextVersionId)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
