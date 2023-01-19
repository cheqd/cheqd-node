//go:build integration

package integration

import (
	"crypto/ed25519"
	"fmt"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	clitypes "github.com/cheqd/cheqd-node/x/did/client/cli"
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
	resourcecli "github.com/cheqd/cheqd-node/x/resource/client/cli"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive resource", func() {
	var tmpDir string

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()
	})

	It("can create diddoc, create resource, query it, query all resource versions of the same resource name, query resource collection", func() {
		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can create diddoc"))
		// Create a new DID Doc
		collectionID := uuid.NewString()
		did := "did:cheqd:" + network.DidNamespace + ":" + collectionID
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(pubKey)

		payload := types.MsgCreateDidDocPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyId,
					VerificationMethodType: "Ed25519VerificationKey2020",
					Controller:             did,
					VerificationMaterial:   publicKeyMultibase,
				},
			},
			Authentication: []string{keyId},
			VersionId:      uuid.NewString(),
		}

		signInputs := []clitypes.SignInput{
			{
				VerificationMethodID: keyId,
				PrivKey:              privKey,
			},
		}

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_1, cli.CliGasParams)
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

		res, err = cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionID:    collectionID,
			ResourceID:      resourceID,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_1, cli.CliGasParams)
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

		res, err = cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionID:    collectionID,
			ResourceID:      nextResourceId,
			ResourceName:    nextResourceName,
			ResourceVersion: nextResourceVersion,
			ResourceType:    nextResourceType,
			ResourceFile:    nextResourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_1, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query resource"))
		// Query the Resource
		res2, err = cli.QueryResource(collectionID, resourceID)
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
		res3, err = cli.QueryResourceMetadata(collectionID, resourceID)
		Expect(err).To(BeNil())

		Expect(res3.Resource.CollectionId).To(BeEquivalentTo(collectionID))
		Expect(res3.Resource.Id).To(BeEquivalentTo(resourceID))
		Expect(res3.Resource.Name).To(BeEquivalentTo(resourceName))
		Expect(res3.Resource.Version).To(BeEquivalentTo(resourceVersion))
		Expect(res3.Resource.ResourceType).To(BeEquivalentTo(resourceType))
		Expect(res3.Resource.MediaType).To(Equal("application/json"))

		// Create a second DID Doc
		secondCollectionId := uuid.NewString()
		secondDid := "did:cheqd:" + network.DidNamespace + ":" + secondCollectionId
		secondKeyId := secondDid + "#key1"

		secondPubKey, secondPrivKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		secondpubKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(secondPubKey)

		secondPayload := types.MsgCreateDidDocPayload{
			Id: secondDid,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     secondKeyId,
					VerificationMethodType: "Ed25519VerificationKey2020",
					Controller:             secondDid,
					VerificationMaterial:   secondpubKeyMultibase,
				},
			},
			Authentication: []string{secondKeyId},
			VersionId:      uuid.NewString(),
		}

		secondSignInputs := []clitypes.SignInput{
			{
				VerificationMethodID: secondKeyId,
				PrivKey:              secondPrivKey,
			},
		}

		res, err = cli.CreateDidDoc(tmpDir, secondPayload, secondSignInputs, testdata.BASE_ACCOUNT_2, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Create a second Resource
		secondResourceId := uuid.NewString()
		secondResourceName := "TestResource2"
		secondResourceVersion := "1.0"
		secondResourceType := "TestType2"
		secondResourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		res, err = cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionID:    secondCollectionId,
			ResourceID:      secondResourceId,
			ResourceName:    secondResourceName,
			ResourceVersion: secondResourceVersion,
			ResourceType:    secondResourceType,
			ResourceFile:    secondResourceFile,
		}, secondSignInputs, testdata.BASE_ACCOUNT_1, cli.CliGasParams)
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
