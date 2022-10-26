package tests

import (
	"fmt"

	. "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/multiformats/go-multibase"

	"github.com/stretchr/testify/require"
)

var _ = Describe("Create DID tests", func() {
	var setup setup.TestSetup

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
		alice := setup.CreateSimpleDid()
		anna := setup.CreateSimpleDid()

		// Bob
		bobDid := GenerateDID(Base58_16chars)
		bobKeypair := GenerateKeyPair()
		bobKeyId := bobDid + "#key-1"

		msg := &types.MsgCreateDidPayload{
			Id:             bobDid,
			Controller:     []string{alice.Did, anna.Did},
			Authentication: []string{bobKeyId},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 bobKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         anna.Did,
					PublicKeyMultibase: MustEncodeBase58(bobKeypair.Public),
				},
			},
		}

		signatures := []SignInput{alice.SignInput, anna.SignInput}

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
			signers: []string{
				"did:cheqd:test:yyyyyyyyyyyyyyyy#key-1",
				"did:cheqd:test:yyyyyyyyyyyyyyyy#key-5",
				AliceKey1,
				BobKey1,
				BobKey2,
				BobKey3,
				CharlieKey1,
				CharlieKey2,
				CharlieKey3,
			},
			msg: &types.MsgCreateDidPayload{
				Id: "did:cheqd:test:yyyyyyyyyyyyyyyy",
				Authentication: []string{
					"did:cheqd:test:yyyyyyyyyyyyyyyy#key-1",
					"did:cheqd:test:yyyyyyyyyyyyyyyy#key-5",
				},
				Context:              []string{"abc", "de"},
				CapabilityInvocation: []string{"did:cheqd:test:yyyyyyyyyyyyyyyy#key-2"},
				CapabilityDelegation: []string{"did:cheqd:test:yyyyyyyyyyyyyyyy#key-3"},
				KeyAgreement:         []string{"did:cheqd:test:yyyyyyyyyyyyyyyy#key-4"},
				AlsoKnownAs:          []string{"SomeUri"},
				Service: []*types.Service{
					{
						Id:              "did:cheqd:test:yyyyyyyyyyyyyyyy#service-1",
						Type:            "DIDCommMessaging",
						ServiceEndpoint: "ServiceEndpoint",
					},
				},
				Controller: []string{"did:cheqd:test:yyyyyyyyyyyyyyyy", AliceDID, BobDID, CharlieDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:yyyyyyyyyyyyyyyy#key-1",
						Type:       Ed25519VerificationKey2020,
						Controller: "did:cheqd:test:yyyyyyyyyyyyyyyy",
					},
					{
						Id:         "did:cheqd:test:yyyyyyyyyyyyyyyy#key-2",
						Type:       Ed25519VerificationKey2020,
						Controller: "did:cheqd:test:yyyyyyyyyyyyyyyy",
					},
					{
						Id:         "did:cheqd:test:yyyyyyyyyyyyyyyy#key-3",
						Type:       Ed25519VerificationKey2020,
						Controller: "did:cheqd:test:yyyyyyyyyyyyyyyy",
					},
					{
						Id:         "did:cheqd:test:yyyyyyyyyyyyyyyy#key-4",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:yyyyyyyyyyyyyyyy",
					},
					{
						Id:         "did:cheqd:test:yyyyyyyyyyyyyyyy#key-5",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:yyyyyyyyyyyyyyyy",
					},
				},
			},
		},
		{
			valid: false,
			name:  "Not Valid: Second controller did not sign request",
			msg: &types.MsgCreateDidPayload{
				Id:         ImposterDID,
				Controller: []string{AliceDID, BobDID},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[AliceKey1],
			},
			errMsg: fmt.Sprintf("signer: %s: signature is required but not found", BobDID),
		},
		{
			valid: false,
			name:  "Not Valid: No signature",
			msg: &types.MsgCreateDidPayload{
				Id:         ImposterDID,
				Controller: []string{AliceDID, BobDID},
			},
			errMsg: fmt.Sprintf("signer: %s: signature is required but not found", AliceDID),
		},
		{
			valid: false,
			name:  "Not Valid: Controller not found",
			msg: &types.MsgCreateDidPayload{
				Id:         ImposterDID,
				Controller: []string{AliceDID, NotFounDID},
			},
			signers: []string{AliceKey1, ImposterKey1},
			keys: map[string]KeyPair{
				AliceKey1:    keys[AliceKey1],
				ImposterKey1: GenerateKeyPair(),
			},
			errMsg: fmt.Sprintf("%s: DID Doc not found", NotFounDID),
		},
		{
			valid: false,
			name:  "Not Valid: Wrong signature",
			msg: &types.MsgCreateDidPayload{
				Id:         ImposterDID,
				Controller: []string{AliceDID},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[BobKey1],
			},
			errMsg: fmt.Sprintf("method id: %s: invalid signature detected", AliceKey1),
		},
		{
			valid: false,
			name:  "Not Valid: DID signed by wrong controller",
			msg: &types.MsgCreateDidPayload{
				Id:             ImposterDID,
				Authentication: []string{ImposterKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 ImposterKey1,
						Type:               Ed25519VerificationKey2020,
						Controller:         ImposterDID,
						PublicKeyMultibase: "z" + base58.Encode(keys[ImposterKey1].PublicKey),
					},
				},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[AliceKey1],
			},
			errMsg: fmt.Sprintf("signer: %s: signature is required but not found", ImposterDID),
		},
		{
			valid: false,
			name:  "Not Valid: DID self-signed by not existing verification method",
			msg: &types.MsgCreateDidPayload{
				Id:             ImposterDID,
				Authentication: []string{ImposterKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 ImposterKey1,
						Type:               Ed25519VerificationKey2020,
						Controller:         ImposterDID,
						PublicKeyMultibase: "z" + base58.Encode(keys[ImposterKey1].PublicKey),
					},
				},
			},
			signers: []string{ImposterKey2},
			keys: map[string]KeyPair{
				ImposterKey2: GenerateKeyPair(),
			},
			errMsg: fmt.Sprintf("%s: verification method not found", ImposterKey2),
		},
		{
			valid: false,
			name:  "Not Valid: Self-signature not found",
			msg: &types.MsgCreateDidPayload{
				Id:             ImposterDID,
				Controller:     []string{AliceDID, ImposterDID},
				Authentication: []string{ImposterKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 ImposterKey1,
						Type:               Ed25519VerificationKey2020,
						Controller:         ImposterDID,
						PublicKeyMultibase: "z" + base58.Encode(keys[ImposterKey1].PublicKey),
					},
				},
			},
			signers: []string{AliceKey1, ImposterKey2},
			keys: map[string]KeyPair{
				AliceKey1:    keys[AliceKey1],
				ImposterKey2: GenerateKeyPair(),
			},
			errMsg: fmt.Sprintf("%s: verification method not found", ImposterKey2),
		},
		{
			valid: false,
			name:  "Not Valid: DID Doc already exists",
			keys: map[string]KeyPair{
				CharlieKey1: GenerateKeyPair(),
			},
			signers: []string{CharlieKey1},
			msg: &types.MsgCreateDidPayload{
				Id:             CharlieDID,
				Authentication: []string{CharlieKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         CharlieKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: CharlieDID,
					},
				},
			},
			errMsg: fmt.Sprintf("%s: DID Doc exists", CharlieDID),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			setup := InitEnv(t, keys)

			for _, vm := range msg.VerificationMethod {
				if vm.PublicKeyMultibase == "" {
					vm.PublicKeyMultibase, err = multibase.Encode(multibase.Base58BTC, tc.keys[vm.Id].PublicKey)
				}
				require.NoError(t, err)
			}

			signerKeys := map[string]ed25519.PrivateKey{}
			for _, signer := range tc.signers {
				signerKeys[signer] = tc.keys[signer].PrivateKey
			}

			did, err := setup.SendCreateDid(msg, signerKeys)

			if tc.valid {
				require.Nil(t, err)
				require.Equal(t, tc.msg.Id, did.Id)
				require.Equal(t, utils.NormalizeIdentifiersList(tc.msg.Controller), did.Controller)
				require.Equal(t, tc.msg.VerificationMethod, did.VerificationMethod)
				require.Equal(t, tc.msg.Authentication, did.Authentication)
				require.Equal(t, tc.msg.AssertionMethod, did.AssertionMethod)
				require.Equal(t, tc.msg.CapabilityInvocation, did.CapabilityInvocation)
				require.Equal(t, tc.msg.CapabilityDelegation, did.CapabilityDelegation)
				require.Equal(t, tc.msg.KeyAgreement, did.KeyAgreement)
				require.Equal(t, tc.msg.AlsoKnownAs, did.AlsoKnownAs)
				require.Equal(t, tc.msg.Service, did.Service)
				require.Equal(t, tc.msg.Context, did.Context)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}

func TestHandler_DidDocAlreadyExists(t *testing.T) {
	setup := Setup()

	_, _, _ = setup.InitDid(AliceDID)
	_, _, err := setup.InitDid(AliceDID)

	require.Error(t, err)
	require.Equal(t, fmt.Sprintf("%s: DID Doc exists", AliceDID), err.Error())
}
