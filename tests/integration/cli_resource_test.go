//go:build integration

package integration

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	cli_types "github.com/cheqd/cheqd-node/x/cheqd/client/cli"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli", func() {
	It("can create diddoc, create resource, query it", func() {
		// Create a new DID Doc
		collectionId := uuid.NewString()
		did := "did:cheqd:" + network.DID_NAMESPACE + ":" + collectionId
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
		Expect(err).To(BeNil())

		payload := types.MsgCreateDidPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 keyId,
					Type:               "Ed25519VerificationKey2020",
					Controller:         did,
					PublicKeyMultibase: string(pubKeyMultibase58),
				},
			},
			Authentication: []string{keyId},
		}

		signInputs := []cli_types.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey,
			},
		}

		res, err := cli.CreateDid(payload, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Create a new Resource
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		res, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Query the Resource
		res2, err := cli.QueryResource(collectionId, resourceId)
		Expect(err).To(BeNil())

		Expect(res2.Resource.Header.CollectionId).To(BeEquivalentTo(collectionId))
		Expect(res2.Resource.Header.Id).To(BeEquivalentTo(resourceId))
		Expect(res2.Resource.Header.Name).To(BeEquivalentTo(resourceName))
		Expect(res2.Resource.Header.ResourceType).To(BeEquivalentTo(resourceType))
		Expect(res2.Resource.Header.MediaType).To(Equal("application/json"))
		Expect(res2.Resource.Data).To(BeEquivalentTo(testdata.JSON_FILE_CONTENT))
	})
})
