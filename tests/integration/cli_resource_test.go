//go:build integration

package integration

import (
	"crypto/ed25519"
	"fmt"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	cli_types "github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive resource", func() {
	It("can create diddoc, create resource, query it, query all resource versions of the same resource name, query resource collection", func() {
		// Create a new DID Doc
		collectionId := uuid.NewString()
		did := "did:cheqd:" + network.DID_NAMESPACE + ":" + collectionId
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
		Expect(err).To(BeNil())

		payload := types.MsgCreateDidDocPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 "Ed25519VerificationKey2020",
					Controller:           did,
					VerificationMaterial: "{\"publicKeyMultibase\": \"" + string(pubKeyMultibase58) + "\"}",
				},
			},
			Authentication: []string{keyId},
			VersionId:      uuid.NewString(),
		}

		signInputs := []cli_types.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey,
			},
		}

		res, err := cli.CreateDidDoc(payload, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.GREEN, "can create resource"))
		// Create a new Resource
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		res, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceVersion, resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.GREEN, "can query resource"))
		// Query the Resource
		res2, err := cli.QueryResource(collectionId, resourceId)
		Expect(err).To(BeNil())

		Expect(res2.Resource.Metadata.CollectionId).To(BeEquivalentTo(collectionId))
		Expect(res2.Resource.Metadata.Id).To(BeEquivalentTo(resourceId))
		Expect(res2.Resource.Metadata.Name).To(BeEquivalentTo(resourceName))
		Expect(res2.Resource.Metadata.ResourceType).To(BeEquivalentTo(resourceType))
		Expect(res2.Resource.Metadata.MediaType).To(Equal("application/json"))
		Expect(res2.Resource.Resource.Data).To(BeEquivalentTo(testdata.JSON_FILE_CONTENT))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.GREEN, "can query resource metadata"))
		// Query the Resource Metadata
		res3, err := cli.QueryResourceMetadata(collectionId, resourceId)
		Expect(err).To(BeNil())

		Expect(res3.Resource.CollectionId).To(BeEquivalentTo(collectionId))
		Expect(res3.Resource.Id).To(BeEquivalentTo(resourceId))
		Expect(res3.Resource.Name).To(BeEquivalentTo(resourceName))
		Expect(res3.Resource.ResourceType).To(BeEquivalentTo(resourceType))
		Expect(res3.Resource.MediaType).To(Equal("application/json"))

		// Create Resource next version
		nextResourceId := uuid.NewString()
		nextResourceName := resourceName
		nextResourceVersion := "2.0"
		nextResourceType := resourceType
		nextResourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		res, err = cli.CreateResource(collectionId, nextResourceId, nextResourceName, nextResourceVersion, nextResourceType, nextResourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Create a second DID Doc
		secondCollectionId := uuid.NewString()
		secondDid := "did:cheqd:" + network.DID_NAMESPACE + ":" + secondCollectionId
		secondKeyId := secondDid + "#key1"

		secondPubKey, secondPrivKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		secondPubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, secondPubKey)
		Expect(err).To(BeNil())

		secondPayload := types.MsgCreateDidDocPayload{
			Id: secondDid,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   secondKeyId,
					Type:                 "Ed25519VerificationKey2020",
					Controller:           secondDid,
					VerificationMaterial: "{\"publicKeyMultibase\": \"" + string(secondPubKeyMultibase58) + "\"}",
				},
			},
			Authentication: []string{secondKeyId},
			VersionId:      uuid.NewString(),
		}

		secondSignInputs := []cli_types.SignInput{
			{
				VerificationMethodId: secondKeyId,
				PrivKey:              secondPrivKey,
			},
		}

		res, err = cli.CreateDidDoc(secondPayload, secondSignInputs, testdata.BASE_ACCOUNT_2)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Create a second Resource
		secondResourceId := uuid.NewString()
		secondResourceName := "TestResource2"
		secondResourceVersion := "1.0"
		secondResourceType := "TestType2"
		secondResourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		res, err = cli.CreateResource(secondCollectionId, secondResourceId, secondResourceName, secondResourceVersion, secondResourceType, secondResourceFile, secondSignInputs, testdata.BASE_ACCOUNT_2)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.GREEN, "can query resource collection"))
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
