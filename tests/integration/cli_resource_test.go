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
	It("can create diddoc, create resource, query it, query all resource versions of the same resource name, query resource collection", func() {
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

		// Create Resource next version
		nextResourceId := uuid.NewString()
		nextResourceName := resourceName
		nextResourceType := resourceType
		nextResourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		res, err = cli.CreateResource(collectionId, nextResourceId, nextResourceName, nextResourceType, nextResourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Query all Resource versions
		res3, err := cli.QueryAllResourceVersions(collectionId, resourceName)
		Expect(err).To(BeNil())
		Expect(len(res3.Resources)).To(Equal(2))
		Expect(res3.Resources[0].CollectionId).To(Equal(collectionId))
		Expect(res3.Resources[1].CollectionId).To(Equal(collectionId))
		Expect([]string{res3.Resources[0].Id, res3.Resources[1].Id}).To(ContainElements(resourceId, nextResourceId))

		// Create a second DID Doc
		secondCollectionId := uuid.NewString()
		secondDid := "did:cheqd:" + network.DID_NAMESPACE + ":" + secondCollectionId
		secondKeyId := secondDid + "#key1"

		secondPubKey, secondPrivKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		secondPubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, secondPubKey)
		Expect(err).To(BeNil())

		secondPayload := types.MsgCreateDidPayload{
			Id: secondDid,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 secondKeyId,
					Type:               "Ed25519VerificationKey2020",
					Controller:         secondDid,
					PublicKeyMultibase: string(secondPubKeyMultibase58),
				},
			},
			Authentication: []string{secondKeyId},
		}

		secondSignInputs := []cli_types.SignInput{
			{
				VerificationMethodId: secondKeyId,
				PrivKey:              secondPrivKey,
			},
		}

		res, err = cli.CreateDid(secondPayload, secondSignInputs, testdata.BASE_ACCOUNT_2)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Create a second Resource
		secondResourceId := uuid.NewString()
		secondResourceName := "TestResource2"
		secondResourceType := "TestType2"
		secondResourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		res, err = cli.CreateResource(secondCollectionId, secondResourceId, secondResourceName, secondResourceType, secondResourceFile, secondSignInputs, testdata.BASE_ACCOUNT_2)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Query Resource Collection
		res4, err := cli.QueryResourceCollection(collectionId)
		Expect(err).To(BeNil())
		Expect(len(res4.Resources)).To(Equal(2))
		Expect(res4.Resources[0].CollectionId).To(Equal(collectionId))
		Expect(res4.Resources[1].CollectionId).To(Equal(collectionId))
		Expect([]string{res4.Resources[0].Id, res4.Resources[1].Id}).To(ContainElements(resourceId, nextResourceId))

		// Query second Resource Collection
		res5, err := cli.QueryResourceCollection(secondCollectionId)
		Expect(err).To(BeNil())
		Expect(len(res5.Resources)).To(Equal(1))
		Expect(res5.Resources[0].CollectionId).To(Equal(secondCollectionId))
		Expect(res5.Resources[0].Id).To(Equal(secondResourceId))
	})
})
