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
	"github.com/tendermint/tendermint/libs/rand"
)

var _ = Describe("cheqd cli negative", func() {
	var collectionId = ""

	BeforeEach( func ()  {
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
	})

	It("Create Resource. Cannot create resource missing arguments, wrong collectionId", func() {
		// *********************** Negative cases ***********************

		println("*********************** Create Resource negative cases start ***********************")

		collectionId := uuid.NewString()
		did := "did:cheqd:" + network.DID_NAMESPACE + ":" + collectionId
		keyId := did + "#key1"

		_, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		signInputs := []cli_types.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey,
			},
		}



		// Create resource Resource with invalid DID
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		// Wrong CollectionId
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// Missing 
		// "--collection-id", collectionId,
		// "--resource-id", resourceId,
		// "--resource-name", resourceName,
		// "--resource-type", resourceType,
		// "--resource-file", resourceFile,

		// Missing collectionId
		_, err = cli.CreateResource("", resourceId, resourceName, resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// Missing resourceId
		_, err = cli.CreateResource(collectionId, "", resourceName, resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// Missing resourceName
		_, err = cli.CreateResource(collectionId, resourceId, "", resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// Missing resourceType
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, "", resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// Missing resourceFile
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, "", signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// Empty signInputs
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, []cli_types.SignInput{}, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// Missing from
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, signInputs, "")
		Expect(err).To(HaveOccurred())

		println("*********************** Create Resource negative cases finish ***********************")
	})


	It("Query Resource. Missing/wrong arguments", func() {
		// *********************** Negative cases ***********************

		// Create a new DID Doc

		println("*********************** Query Resource negative cases start ***********************")

		resourceId := uuid.NewString()

		// Unknown resourceId
		_, err := cli.QueryResource(collectionId, resourceId)
		Expect(err).To(HaveOccurred())

		collectionId = uuid.NewString()

		// Unknown collectionId
		_, err = cli.QueryResource(collectionId, resourceId)
		Expect(err).To(HaveOccurred())

		println("*********************** Query Resource negative cases finish ***********************")
	})

	It("QueryAllResourceVersions. Missing/wrong arguments", func() {
		// *********************** Negative cases ***********************

		// Create a new DID Doc

		println("*********************** QueryAllResourceVersions negative cases start ***********************")

		resourceId := uuid.NewString()
		resourceName := rand.Str(10)

		// Unknown resourceName
		_, err := cli.QueryAllResourceVersions(collectionId, resourceName)
		Expect(err).To(HaveOccurred())

		collectionId = uuid.NewString()

		// Unknown collectionId
		_, err = cli.QueryResource(collectionId, resourceId)
		Expect(err).To(HaveOccurred())

		println("*********************** QueryAllResourceVersions negative cases finish ***********************")
	})

	It("QueryResourceCollection. Missing/wrong arguments", func() {
		// *********************** Negative cases ***********************

		// Create a new DID Doc

		println("*********************** QueryResourceCollection negative cases start ***********************")

		collectionId = uuid.NewString()

		// Unknown collectionId
		_, err := cli.QueryResourceCollection(collectionId)
		Expect(err).To(HaveOccurred())

		println("*********************** QueryResourceCollection negative cases finish ***********************")
	})
})