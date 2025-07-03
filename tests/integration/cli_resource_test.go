//go:build integration

package integration

import (
	"crypto/ed25519"
	"fmt"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive resource", func() {
	var tmpDir string
	var didFeeParams didtypes.FeeParams
	var resourceFeeParams resourcetypes.FeeParams

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		// Query did fee params
		didRes, err := cli.QueryDidParams()
		Expect(err).To(BeNil())

		didFeeParams = didRes.Params

		// Query resource fee params
		resourceRes, err := cli.QueryResourceParams()
		Expect(err).To(BeNil())

		resourceFeeParams = resourceRes.Params
	})

	It("can create diddoc, create resource, query it, query all resource versions of the same resource name, query resource collection", func() {
		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can create diddoc"))
		// Create a new DID Doc
		collectionID := uuid.NewString()
		did := "did:cheqd:" + network.DidNamespace + ":" + collectionID
		keyId := did + "#key1"

		publicKey, privateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(publicKey)

		payload := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]interface{}{
					"id":                 keyId,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": publicKeyMultibase,
				},
			},
			Authentication: []string{keyId},
		}

		signInputs := []didcli.SignInput{
			{
				VerificationMethodID: keyId,
				PrivKey:              privateKey,
			},
		}

		versionID := uuid.NewString()
		useMin := false
		tax, err := cli.ResolveFeeFromParams(didFeeParams.CreateDid, useMin)
		Expect(err).To(BeNil())
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(tax.Amount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Create a new Resource
		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can create resource"))

		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		useMin = false
		tax, err = cli.ResolveFeeFromParams(resourceFeeParams.Json, useMin)
		Expect(err).To(BeNil())

		res, err = cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(tax.Amount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query resource"))
		// Query the Resource
		res2, err := cli.QueryResource(collectionID, resourceID)
		Expect(err).To(BeNil())

		Expect(res2.Resource.Metadata.CollectionId).To(BeEquivalentTo(collectionID))
		Expect(res2.Resource.Metadata.Id).To(BeEquivalentTo(resourceID))
		Expect(res2.Resource.Metadata.Name).To(BeEquivalentTo(resourceName))
		Expect(res2.Resource.Metadata.Version).To(BeEquivalentTo(resourceVersion))
		Expect(res2.Resource.Metadata.ResourceType).To(BeEquivalentTo(resourceType))
		Expect(res2.Resource.Metadata.MediaType).To(Equal("application/json"))
		Expect(res2.Resource.Resource.Data).To(BeEquivalentTo(testdata.JSON_FILE_CONTENT))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query resource metadata"))
		// Query the Resource Metadata
		res3, err := cli.QueryResourceMetadata(collectionID, resourceID)
		Expect(err).To(BeNil())

		Expect(res3.Resource.CollectionId).To(BeEquivalentTo(collectionID))
		Expect(res3.Resource.Id).To(BeEquivalentTo(resourceID))
		Expect(res3.Resource.Name).To(BeEquivalentTo(resourceName))
		Expect(res3.Resource.Version).To(BeEquivalentTo(resourceVersion))
		Expect(res3.Resource.ResourceType).To(BeEquivalentTo(resourceType))
		Expect(res3.Resource.MediaType).To(Equal("application/json"))

		// Create Resource next version
		nextResourceId := uuid.NewString()
		nextResourceName := resourceName
		nextResourceVersion := "2.0"
		nextResourceType := resourceType
		nextResourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		res, err = cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           nextResourceId,
			Name:         nextResourceName,
			Version:      nextResourceVersion,
			ResourceType: nextResourceType,
		}, signInputs, nextResourceFile, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(tax.Amount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query resource next version"))
		// Query the Resource's next version
		res2, err = cli.QueryResource(collectionID, nextResourceId)
		Expect(err).To(BeNil())

		Expect(res2.Resource.Metadata.CollectionId).To(BeEquivalentTo(collectionID))
		Expect(res2.Resource.Metadata.Id).To(BeEquivalentTo(nextResourceId))
		Expect(res2.Resource.Metadata.Name).To(BeEquivalentTo(nextResourceName))
		Expect(res2.Resource.Metadata.Version).To(BeEquivalentTo(nextResourceVersion))
		Expect(res2.Resource.Metadata.ResourceType).To(BeEquivalentTo(nextResourceType))
		Expect(res2.Resource.Metadata.MediaType).To(Equal("application/json"))
		Expect(res2.Resource.Resource.Data).To(BeEquivalentTo(testdata.JSON_FILE_CONTENT))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query resource metadata"))
		// Query the Resource's next version Metadata
		res3, err = cli.QueryResourceMetadata(collectionID, nextResourceId)
		Expect(err).To(BeNil())

		Expect(res3.Resource.CollectionId).To(BeEquivalentTo(collectionID))
		Expect(res3.Resource.Id).To(BeEquivalentTo(nextResourceId))
		Expect(res3.Resource.Name).To(BeEquivalentTo(nextResourceName))
		Expect(res3.Resource.Version).To(BeEquivalentTo(nextResourceVersion))
		Expect(res3.Resource.ResourceType).To(BeEquivalentTo(nextResourceType))
		Expect(res3.Resource.MediaType).To(Equal("application/json"))

		// Create a second DID Doc
		secondCollectionId := uuid.NewString()
		secondDid := "did:cheqd:" + network.DidNamespace + ":" + secondCollectionId
		secondKeyId := secondDid + "#key1"

		secondPublicKey, secondPrivateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		secondPublicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(secondPublicKey)

		secondPayload := didcli.DIDDocument{
			ID: secondDid,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 secondKeyId,
					"type":               "Ed25519VerificationKey2020",
					"controller":         secondDid,
					"publicKeyMultibase": secondPublicKeyMultibase,
				},
			},
			Authentication: []string{secondKeyId},
		}

		secondSignInputs := []didcli.SignInput{
			{
				VerificationMethodID: secondKeyId,
				PrivKey:              secondPrivateKey,
			},
		}

		versionID = uuid.NewString()

		useMin = false
		tax, err = cli.ResolveFeeFromParams(didFeeParams.CreateDid, useMin)
		Expect(err).To(BeNil())
		res, err = cli.CreateDidDoc(tmpDir, secondPayload, secondSignInputs, versionID, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(tax.Amount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Create a second Resource
		secondResourceId := uuid.NewString()
		secondResourceName := "TestResource2"
		secondResourceVersion := "1.0"
		secondResourceType := "TestType2"
		secondResourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		useMin = true
		tax, err = cli.ResolveFeeFromParams(resourceFeeParams.Json, useMin)
		Expect(err).To(BeNil())

		res, err = cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: secondCollectionId,
			Id:           secondResourceId,
			Name:         secondResourceName,
			Version:      secondResourceVersion,
			ResourceType: secondResourceType,
		}, secondSignInputs, secondResourceFile, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(tax.Amount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query resource"))
		// Query the Resource
		res2, err = cli.QueryResource(secondCollectionId, secondResourceId)
		Expect(err).To(BeNil())

		Expect(res2.Resource.Metadata.CollectionId).To(BeEquivalentTo(secondCollectionId))
		Expect(res2.Resource.Metadata.Id).To(BeEquivalentTo(secondResourceId))
		Expect(res2.Resource.Metadata.Name).To(BeEquivalentTo(secondResourceName))
		Expect(res2.Resource.Metadata.Version).To(BeEquivalentTo(secondResourceVersion))
		Expect(res2.Resource.Metadata.ResourceType).To(BeEquivalentTo(secondResourceType))
		Expect(res2.Resource.Metadata.MediaType).To(Equal("application/json"))
		Expect(res2.Resource.Resource.Data).To(BeEquivalentTo(testdata.JSON_FILE_CONTENT))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query resource metadata"))
		// Query the Resource Metadata
		res3, err = cli.QueryResourceMetadata(secondCollectionId, secondResourceId)
		Expect(err).To(BeNil())

		Expect(res3.Resource.CollectionId).To(BeEquivalentTo(secondCollectionId))
		Expect(res3.Resource.Id).To(BeEquivalentTo(secondResourceId))
		Expect(res3.Resource.Name).To(BeEquivalentTo(secondResourceName))
		Expect(res3.Resource.Version).To(BeEquivalentTo(secondResourceVersion))
		Expect(res3.Resource.ResourceType).To(BeEquivalentTo(secondResourceType))
		Expect(res3.Resource.MediaType).To(Equal("application/json"))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query resource collection"))
		// Query Resource Collection
		res4, err := cli.QueryResourceCollection(collectionID)
		Expect(err).To(BeNil())
		Expect(len(res4.Resources)).To(Equal(2))
		Expect(res4.Resources[0].CollectionId).To(Equal(collectionID))
		Expect(res4.Resources[1].CollectionId).To(Equal(collectionID))
		Expect([]string{res4.Resources[0].Id, res4.Resources[1].Id}).To(ContainElements(resourceID, nextResourceId))

		// Query second Resource Collection
		res5, err := cli.QueryResourceCollection(secondCollectionId)
		Expect(err).To(BeNil())
		Expect(len(res5.Resources)).To(Equal(1))
		Expect(res5.Resources[0].CollectionId).To(Equal(secondCollectionId))
		Expect(res5.Resources[0].Id).To(Equal(secondResourceId))
	})
})
