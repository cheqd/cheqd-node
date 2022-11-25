package tests

import (
	"fmt"

	. "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
)

var _ = Describe("Create DID tests", func() {
	var setup setup.TestSetup

	BeforeEach(func() {
		setup = Setup()
	})

	It("Valid: Works for simple DIDDoc (Ed25519VerificationKey2020)", func() {
		did := GenerateDID(Base58_16bytes)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []SignInput{
			{
				VerificationMethodId: keyId,
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
		did := GenerateDID(Base58_16bytes)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 types.JsonWebKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildJsonWebKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []SignInput{
			{
				VerificationMethodId: keyId,
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
		bobDid := GenerateDID(Base58_16bytes)
		bobKeypair := GenerateKeyPair()
		bobKeyId := bobDid + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             bobDid,
			Controller:     []string{alice.Did, anna.Did},
			Authentication: []string{bobKeyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   bobKeyId,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           anna.Did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(bobKeypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []SignInput{alice.SignInput, anna.SignInput}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDidDoc(bobDid)
		Expect(err).To(BeNil())
		Expect(msg.ToDidDoc()).To(Equal(*created.Value.DidDoc))
	})

	It("Valid: Works for DIDDoc with all properties", func() {
		did := GenerateDID(Base58_16bytes)

		keypair1 := GenerateKeyPair()
		keyId1 := did + "#key-1"

		keypair2 := GenerateKeyPair()
		keyId2 := did + "#key-2"

		keypair3 := GenerateKeyPair()
		keyId3 := did + "#key-3"

		keypair4 := GenerateKeyPair()
		keyId4 := did + "#key-4"

		msg := &types.MsgCreateDidDocPayload{
			Context:    []string{"abc", "def"},
			Id:         did,
			Controller: []string{did},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId1,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keypair1.Public),
				},
				{
					Id:                   keyId2,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keypair2.Public),
				},
				{
					Id:                   keyId3,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keypair3.Public),
				},
				{
					Id:                   keyId4,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keypair4.Public),
				},
			},
			Authentication:       []string{keyId1, keyId2},
			AssertionMethod:      []string{keyId3},
			CapabilityInvocation: []string{keyId4, keyId1},
			CapabilityDelegation: []string{keyId4, keyId2},
			KeyAgreement:         []string{keyId1, keyId2, keyId3, keyId4},
			Service: []*types.Service{
				{
					Id:              did + "#service-1",
					Type:            "type-1",
					ServiceEndpoint: []string{"endpoint-1"},
				},
			},
			AlsoKnownAs: []string{"alias-1", "alias-2"},
			VersionId:   uuid.NewString(),
		}

		signatures := []SignInput{
			{
				VerificationMethodId: keyId1,
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
		bobDid := GenerateDID(Base58_16bytes)
		bobKeypair := GenerateKeyPair()
		bobKeyId := bobDid + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             bobDid,
			Controller:     []string{alice.Did, bobDid},
			Authentication: []string{bobKeyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   bobKeyId,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           bobDid,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(bobKeypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []SignInput{
			{
				VerificationMethodId: bobKeyId,
				Key:                  bobKeypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("signer: %s: signature is required but not found", alice.Did)))
	})

	It("Not Valid: No signature", func() {
		did := GenerateDID(Base58_16bytes)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Controller:     []string{did},
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []SignInput{}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("signer: %s: signature is required but not found", did)))
	})

	It("Not Valid: Controller not found", func() {
		did := GenerateDID(Base58_16bytes)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		nonExistingDid := GenerateDID(Base58_16bytes)

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Controller:     []string{nonExistingDid},
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []SignInput{}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("%s: DID Doc not found", nonExistingDid)))
	})

	It("Not Valid: Wrong signature", func() {
		did := GenerateDID(Base58_16bytes)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Controller:     []string{did},
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		invalidKey := GenerateKeyPair()

		signatures := []SignInput{
			{
				VerificationMethodId: keyId,
				Key:                  invalidKey.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("method id: %s: invalid signature detected", keyId)))
	})

	It("Not Valid: DID signed by wrong controller", func() {
		// Alice
		alice := setup.CreateSimpleDid()

		// Bob
		bobDid := GenerateDID(Base58_16bytes)
		bobKeypair := GenerateKeyPair()
		bobKeyId := bobDid + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             bobDid,
			Controller:     []string{bobDid},
			Authentication: []string{bobKeyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   bobKeyId,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           bobDid,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(bobKeypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []SignInput{alice.SignInput}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("signer: %s: signature is required but not found", bobDid)))
	})

	It("Not Valid: DID signed by invalid verification method", func() {
		did := GenerateDID(Base58_16bytes)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Controller:     []string{did},
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		invalidKeyId := did + "#key-2"

		signatures := []SignInput{
			{
				VerificationMethodId: invalidKeyId,
				Key:                  keypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("%s: verification method not found", invalidKeyId)))
	})

	It("Not Valid: DIDDoc already exists", func() {
		// Alice
		alice := setup.CreateSimpleDid()

		msg := &types.MsgCreateDidDocPayload{
			Id:             alice.Did,
			Authentication: []string{alice.KeyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   alice.KeyId,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           alice.Did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
				},
			},
		}

		signatures := []SignInput{alice.SignInput}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("%s: DID Doc exists", alice.Did)))
	})
})

var _ = Describe("Check upper/lower case for DID creation", func() {
	var setup setup.TestSetup
	var didPrefix string = "did:cheqd:testnet:"

	type TestCaseUUIDDidStruct struct {
		inputId  string
		resultId string
	}

	DescribeTable("Check upper/lower case for DID creation", func(testCase TestCaseUUIDDidStruct) {
		setup = Setup()
		did := testCase.inputId
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		msg := &types.MsgCreateDidDocPayload{
			Id:             did,
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 types.Ed25519VerificationKey2020{}.Type(),
					Controller:           did,
					VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keypair.Public),
				},
			},
			VersionId: uuid.NewString(),
		}

		signatures := []SignInput{
			{
				VerificationMethodId: keyId,
				Key:                  keypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDidDoc(did)
		Expect(err).To(BeNil())
		Expect(created.Value.DidDoc.Id).To(Equal(testCase.resultId))
	},

		Entry("Lowercase UUIDs", TestCaseUUIDDidStruct{
			inputId:  didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			resultId: didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		}),
		Entry("Uppercase UUIDs", TestCaseUUIDDidStruct{
			inputId:  didPrefix + "A86F9CAE-0902-4A7C-A144-96B60CED2FC9",
			resultId: didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		}),
		Entry("Mixed case UUIDs", TestCaseUUIDDidStruct{
			inputId:  didPrefix + "A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
			resultId: didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		}),
		Entry("Indy-style IDs", TestCaseUUIDDidStruct{
			inputId:  didPrefix + "zABCDEFG123456789abcd",
			resultId: didPrefix + "zABCDEFG123456789abcd",
		}),
	)
})
