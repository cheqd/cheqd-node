package tests

import (
	"fmt"

	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
)

var _ = Describe("Create DID tests", func() {
	var setup testsetup.TestSetup

	BeforeEach(func() {
		setup = testsetup.Setup()
	})

	It("Valid: Works for simple DIDDoc (Ed25519VerificationKey2020)", func() {
		did := testsetup.GenerateDID(testsetup.Base58_16bytes)
		keypair := testsetup.GenerateKeyPair()
		keyID := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Authentication: []string{keyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyID,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []testsetup.SignInput{
			{
				VerificationMethodID: keyID,
				Key:                  keypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDidDoc(did)
		Expect(err).To(BeNil())
		Expect(msg.ToDidDoc()).To(Equal(*created.Value.DidDoc))
	})

	It("Valid: Works for simple DIDDoc (JsonWebKey2020)", func() {
		did := testsetup.GenerateDID(testsetup.Base58_16bytes)
		keypair := testsetup.GenerateKeyPair()
		keyID := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Authentication: []string{keyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyID,
					VerificationMethodType: types.JSONWebKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildJSONWebKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []testsetup.SignInput{
			{
				VerificationMethodID: keyID,
				Key:                  keypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDidDoc(did)
		Expect(err).To(BeNil())
		Expect(msg.ToDidDoc()).To(Equal(*created.Value.DidDoc))
	})

	It("Valid: DID with external controllers", func() {
		// Alice
		alice := setup.CreateSimpleDid()
		anna := setup.CreateSimpleDid()

		// Bob
		bobDid := testsetup.GenerateDID(testsetup.Base58_16bytes)
		bobKeypair := testsetup.GenerateKeyPair()
		bobKeyID := bobDid + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             bobDid,
			Controller:     []string{alice.Did, anna.Did},
			Authentication: []string{bobKeyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     bobKeyID,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             anna.Did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(bobKeypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []testsetup.SignInput{alice.SignInput, anna.SignInput}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDidDoc(bobDid)
		Expect(err).To(BeNil())
		Expect(msg.ToDidDoc()).To(Equal(*created.Value.DidDoc))
	})

	It("Valid: Works for DIDDoc with all properties", func() {
		did := testsetup.GenerateDID(testsetup.Base58_16bytes)

		keypair1 := testsetup.GenerateKeyPair()
		keyID1 := did + "#key-1"

		keypair2 := testsetup.GenerateKeyPair()
		keyID2 := did + "#key-2"

		keypair3 := testsetup.GenerateKeyPair()
		keyID3 := did + "#key-3"

		keypair4 := testsetup.GenerateKeyPair()
		keyID4 := did + "#key-4"

		msg := &types.MsgCreateDidDocPayload{
			Context:    []string{"abc", "def"},
			Id:         did,
			Controller: []string{did},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyID1,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(keypair1.Public),
				},
				{
					Id:                     keyID2,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(keypair2.Public),
				},
				{
					Id:                     keyID3,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(keypair3.Public),
				},
				{
					Id:                     keyID4,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(keypair4.Public),
				},
			},
			Authentication:       []string{keyID1, keyID2},
			AssertionMethod:      []string{keyID3},
			CapabilityInvocation: []string{keyID4, keyID1},
			CapabilityDelegation: []string{keyID4, keyID2},
			KeyAgreement:         []string{keyID1, keyID2, keyID3, keyID4},
			Service: []*types.Service{
				{
					Id:              did + "#service-1",
					ServiceType:     "type-1",
					ServiceEndpoint: []string{"endpoint-1"},
				},
			},
			AlsoKnownAs: []string{"alias-1", "alias-2"},
			VersionId:   uuid.NewString(),
		}

		signatures := []testsetup.SignInput{
			{
				VerificationMethodID: keyID1,
				Key:                  keypair1.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDidDoc(did)
		Expect(err).To(BeNil())
		Expect(msg.ToDidDoc()).To(Equal(*created.Value.DidDoc))
	})

	// **************************
	// ***** Negative cases *****
	// **************************

	It("Not Valid: Second controller did not sign request", func() {
		// Alice
		alice := setup.CreateSimpleDid()

		// Bob
		bobDid := testsetup.GenerateDID(testsetup.Base58_16bytes)
		bobKeypair := testsetup.GenerateKeyPair()
		bobKeyID := bobDid + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             bobDid,
			Controller:     []string{alice.Did, bobDid},
			Authentication: []string{bobKeyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     bobKeyID,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             bobDid,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(bobKeypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []testsetup.SignInput{
			{
				VerificationMethodID: bobKeyID,
				Key:                  bobKeypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("signer: %s: signature is required but not found", alice.Did)))
	})

	It("Not Valid: No signature", func() {
		did := testsetup.GenerateDID(testsetup.Base58_16bytes)
		keypair := testsetup.GenerateKeyPair()
		keyID := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Controller:     []string{did},
			Authentication: []string{keyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyID,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []testsetup.SignInput{}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("signer: %s: signature is required but not found", did)))
	})

	It("Not Valid: Controller not found", func() {
		did := testsetup.GenerateDID(testsetup.Base58_16bytes)
		keypair := testsetup.GenerateKeyPair()
		keyID := did + "#key-1"

		nonExistingDid := testsetup.GenerateDID(testsetup.Base58_16bytes)

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Controller:     []string{nonExistingDid},
			Authentication: []string{keyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyID,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []testsetup.SignInput{}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("%s: DID Doc not found", nonExistingDid)))
	})

	It("Not Valid: Wrong signature", func() {
		did := testsetup.GenerateDID(testsetup.Base58_16bytes)
		keypair := testsetup.GenerateKeyPair()
		keyID := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Controller:     []string{did},
			Authentication: []string{keyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyID,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		invalidKey := testsetup.GenerateKeyPair()

		signatures := []testsetup.SignInput{
			{
				VerificationMethodID: keyID,
				Key:                  invalidKey.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("method id: %s: invalid signature detected", keyID)))
	})

	It("Not Valid: DID signed by wrong controller", func() {
		// Alice
		alice := setup.CreateSimpleDid()

		// Bob
		bobDid := testsetup.GenerateDID(testsetup.Base58_16bytes)
		bobKeypair := testsetup.GenerateKeyPair()
		bobKeyID := bobDid + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             bobDid,
			Controller:     []string{bobDid},
			Authentication: []string{bobKeyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     bobKeyID,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             bobDid,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(bobKeypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []testsetup.SignInput{alice.SignInput}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("signer: %s: signature is required but not found", bobDid)))
	})

	It("Not Valid: DID signed by invalid verification method", func() {
		did := testsetup.GenerateDID(testsetup.Base58_16bytes)
		keypair := testsetup.GenerateKeyPair()
		keyID := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Controller:     []string{did},
			Authentication: []string{keyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyID,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		invalidKeyID := did + "#key-2"

		signatures := []testsetup.SignInput{
			{
				VerificationMethodID: invalidKeyID,
				Key:                  keypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("%s: verification method not found", invalidKeyID)))
	})

	It("Not Valid: DIDDoc already exists", func() {
		// Alice
		alice := setup.CreateSimpleDid()

		msg := &types.MsgCreateDidDocPayload{
			Id:             alice.Did,
			Authentication: []string{alice.KeyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     alice.KeyID,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             alice.Did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
				},
			},
		}

		signatures := []testsetup.SignInput{alice.SignInput}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("%s: DID Doc exists", alice.Did)))
	})
})

var _ = Describe("Check upper/lower case for DID creation", func() {
	var setup testsetup.TestSetup
	didPrefix := "did:cheqd:testnet:"

	type TestCaseUUIDDidStruct struct {
		inputID  string
		resultID string
	}

	DescribeTable("Check upper/lower case for DID creation", func(testCase TestCaseUUIDDidStruct) {
		setup = testsetup.Setup()
		did := testCase.inputID
		keypair := testsetup.GenerateKeyPair()
		keyID := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Authentication: []string{keyID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyID,
					VerificationMethodType: types.Ed25519VerificationKey2020Type,
					Controller:             did,
					VerificationMaterial:   testsetup.BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []testsetup.SignInput{
			{
				VerificationMethodID: keyID,
				Key:                  keypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDidDoc(did)
		Expect(err).To(BeNil())
		Expect(created.Value.DidDoc.Id).To(Equal(testCase.resultID))
	},

		Entry("Lowercase UUIDs", TestCaseUUIDDidStruct{
			inputID:  didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			resultID: didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		}),
		Entry("Uppercase UUIDs", TestCaseUUIDDidStruct{
			inputID:  didPrefix + "A86F9CAE-0902-4A7C-A144-96B60CED2FC9",
			resultID: didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		}),
		Entry("Mixed case UUIDs", TestCaseUUIDDidStruct{
			inputID:  didPrefix + "A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
			resultID: didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		}),
		Entry("Indy-style IDs", TestCaseUUIDDidStruct{
			inputID:  didPrefix + "zABCDEFG123456789abcd",
			resultID: didPrefix + "zABCDEFG123456789abcd",
		}),
	)
})
