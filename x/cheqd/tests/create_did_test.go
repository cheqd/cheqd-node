package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("Create DID tests new", func() {
	var setup TestSetup

	BeforeEach(func() {
		setup = Setup()
	})

	It("Valid: Works for simple did doc", func() {
		did := GenerateDID(Base58_16chars)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		msg := &types.MsgCreateDidPayload{
			Id:             did,
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 keyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         did,
					PublicKeyMultibase: MustEncodeBase58(keypair.Public),
				},
			},
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
		created, err := setup.QueryDid(did)
		Expect(err).To(BeNil())
		Expect(msg.ToDid()).To(Equal(*created.Did))
	})

	It("Valid: DID with external controllers", func() {
		// Alice
		aliceDid, aliceKeypair, aliceKeyId := setup.CreateTestDid()
		annaDid, annaKeypair, annaKeyId := setup.CreateTestDid()

		// Bob
		bobDid := GenerateDID(Base58_16chars)
		bobKeypair := GenerateKeyPair()
		bobKeyId := bobDid + "#key-1"

		msg := &types.MsgCreateDidPayload{
			Id:             bobDid,
			Controller:     []string{aliceDid, annaDid},
			Authentication: []string{bobKeyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 bobKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         annaDid,
					PublicKeyMultibase: MustEncodeBase58(bobKeypair.Public),
				},
			},
		}

		signatures := []SignInput{
			{
				VerificationMethodId: aliceKeyId,
				Key:                  aliceKeypair.Private,
			},
			{
				VerificationMethodId: annaKeyId,
				Key:                  annaKeypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDid(bobDid)
		Expect(err).To(BeNil())
		Expect(msg.ToDid()).To(Equal(*created.Did))
	})

	It("Valid: Works for the did doc witha all properties", func() {
		did := GenerateDID(Base58_16chars)

		keypair1 := GenerateKeyPair()
		keyId1 := did + "#key-1"

		keypair2 := GenerateKeyPair()
		keyId2 := did + "#key-2"

		keypair3 := GenerateKeyPair()
		keyId3 := did + "#key-3"

		keypair4 := GenerateKeyPair()
		keyId4 := did + "#key-4"

		msg := &types.MsgCreateDidPayload{
			Context:    []string{"abc", "def"},
			Id:         did,
			Controller: []string{did},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 keyId1,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         did,
					PublicKeyMultibase: MustEncodeBase58(keypair1.Public),
				},
				{
					Id:                 keyId2,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         did,
					PublicKeyMultibase: MustEncodeBase58(keypair2.Public),
				},
				{
					Id:                 keyId3,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         did,
					PublicKeyMultibase: MustEncodeBase58(keypair3.Public),
				},
				{
					Id:                 keyId4,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         did,
					PublicKeyMultibase: MustEncodeBase58(keypair4.Public),
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
					ServiceEndpoint: "endpoint-1",
				},
			},
			AlsoKnownAs: []string{"alias-1", "alias-2"},
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
		created, err := setup.QueryDid(did)
		Expect(err).To(BeNil())
		Expect(msg.ToDid()).To(Equal(*created.Did))
	})

	// **************************
	// ***** Negative cases *****
	// **************************

	It("Not Valid: Second controller did not sign request", func() {
		// Alice
		aliceDid, _, _ := setup.CreateTestDid()

		// Bob
		bobDid := GenerateDID(Base58_16chars)
		bobKeypair := GenerateKeyPair()
		bobKeyId := bobDid + "#key-1"

		msg := &types.MsgCreateDidPayload{
			Id:             bobDid,
			Controller:     []string{aliceDid, bobDid},
			Authentication: []string{bobKeyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 bobKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         bobDid,
					PublicKeyMultibase: MustEncodeBase58(bobKeypair.Public),
				},
			},
		}

		signatures := []SignInput{
			{
				VerificationMethodId: bobKeyId,
				Key:                  bobKeypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("signer: %s: signature is required but not found", aliceDid)))
	})

	It("Not Valid: No signature", func() {
		did := GenerateDID(Base58_16chars)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		msg := &types.MsgCreateDidPayload{
			Id:             did,
			Controller:     []string{did},
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 keyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         did,
					PublicKeyMultibase: MustEncodeBase58(keypair.Public),
				},
			},
		}

		signatures := []SignInput{}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("signer: %s: signature is required but not found", did)))
	})

	It("Not Valid: Controller not found", func() {
		did := GenerateDID(Base58_16chars)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		nonExistingDid := GenerateDID(Base58_16chars)

		msg := &types.MsgCreateDidPayload{
			Id:             did,
			Controller:     []string{nonExistingDid},
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 keyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         did,
					PublicKeyMultibase: MustEncodeBase58(keypair.Public),
				},
			},
		}

		signatures := []SignInput{}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("%s: DID Doc not found", nonExistingDid)))
	})

	It("Not Valid: Wrong signature", func() {
		did := GenerateDID(Base58_16chars)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		msg := &types.MsgCreateDidPayload{
			Id:             did,
			Controller:     []string{did},
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 keyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         did,
					PublicKeyMultibase: MustEncodeBase58(keypair.Public),
				},
			},
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
		_, aliceKeypair, aliceKeyId := setup.CreateTestDid()

		// Bob
		bobDid := GenerateDID(Base58_16chars)
		bobKeypair := GenerateKeyPair()
		bobKeyId := bobDid + "#key-1"

		msg := &types.MsgCreateDidPayload{
			Id:             bobDid,
			Controller:     []string{bobDid},
			Authentication: []string{bobKeyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 bobKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         bobDid,
					PublicKeyMultibase: MustEncodeBase58(bobKeypair.Public),
				},
			},
		}

		signatures := []SignInput{
			{
				VerificationMethodId: aliceKeyId,
				Key:                  aliceKeypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("signer: %s: signature is required but not found", bobDid)))
	})

	It("Not Valid: DID self-signed by not existing verification method", func() {
		did := GenerateDID(Base58_16chars)
		keypair := GenerateKeyPair()
		keyId := did + "#key-1"

		msg := &types.MsgCreateDidPayload{
			Id:             did,
			Controller:     []string{did},
			Authentication: []string{keyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 keyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         did,
					PublicKeyMultibase: MustEncodeBase58(keypair.Public),
				},
			},
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

	It("Not Valid: DID Doc already exists", func() {
		// Alice
		aliceDid, aliceKeypair, aliceKeyId := setup.CreateTestDid()

		msg := &types.MsgCreateDidPayload{
			Id:             aliceDid,
			Authentication: []string{aliceKeyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 aliceKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         aliceDid,
					PublicKeyMultibase: MustEncodeBase58(aliceKeypair.Public),
				},
			},
		}

		signatures := []SignInput{
			{
				VerificationMethodId: aliceKeyId,
				Key:                  aliceKeypair.Private,
			},
		}

		_, err := setup.CreateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("%s: DID Doc exists", aliceDid)))
	})
})
