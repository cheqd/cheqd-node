package tests

import (
	"crypto/ed25519"
	"fmt"
	"testing"

	"github.com/btcsuite/btcutil/base58"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/multiformats/go-multibase"

	"github.com/stretchr/testify/require"
)

func TestCreateDID(t *testing.T) {
	var err error
	keys := GenerateTestKeys()
	cases := []struct {
		valid   bool
		name    string
		keys    map[string]KeyPair
		signers []string
		msg     *types.MsgCreateDidPayload
		errMsg  string
	}{
		{
			valid: true,
			name:  "Valid: Works",
			keys: map[string]KeyPair{
				ImposterKey1: GenerateKeyPair(),
			},
			signers: []string{ImposterKey1},
			msg: &types.MsgCreateDidPayload{
				Id:             ImposterDID,
				Authentication: []string{ImposterKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         ImposterKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: ImposterDID,
					},
				},
			},
		},
		{
			valid: true,
			name:  "Valid: Works with Key Agreement",
			keys: map[string]KeyPair{
				ImposterKey1: GenerateKeyPair(),
				AliceKey1:    keys[AliceKey1],
			},
			signers: []string{ImposterKey1, AliceKey1},
			msg: &types.MsgCreateDidPayload{
				Id:           ImposterDID,
				KeyAgreement: []string{ImposterKey1},
				Controller:   []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         ImposterKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: ImposterDID,
					},
				},
			},
		},
		{
			valid: true,
			name:  "Valid: Works with Assertion Method",
			keys: map[string]KeyPair{
				ImposterKey1: GenerateKeyPair(),
				AliceKey1:    keys[AliceKey1],
			},
			signers: []string{AliceKey1, ImposterKey1},
			msg: &types.MsgCreateDidPayload{
				Id:              ImposterDID,
				AssertionMethod: []string{ImposterKey1},
				Controller:      []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         ImposterKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: ImposterDID,
					},
				},
			},
		},
		{
			valid: true,
			name:  "Valid: Works with Capability Delegation",
			keys: map[string]KeyPair{
				ImposterKey1: GenerateKeyPair(),
				AliceKey1:    keys[AliceKey1],
			},
			signers: []string{AliceKey1, ImposterKey1},
			msg: &types.MsgCreateDidPayload{
				Id:                   ImposterDID,
				CapabilityDelegation: []string{ImposterKey1},
				Controller:           []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         ImposterKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: ImposterDID,
					},
				},
			},
		},
		{
			valid: true,
			name:  "Valid: Works with Capability Invocation",
			keys: map[string]KeyPair{
				ImposterKey1: GenerateKeyPair(),
				AliceKey1:    keys[AliceKey1],
			},
			signers: []string{AliceKey1, ImposterKey1},
			msg: &types.MsgCreateDidPayload{
				Id:                   ImposterDID,
				CapabilityInvocation: []string{ImposterKey1},
				Controller:           []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         ImposterKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: ImposterDID,
					},
				},
			},
		},
		{
			valid: true,
			name:  "Valid: With controller works",
			msg: &types.MsgCreateDidPayload{
				Id:         ImposterDID,
				Controller: []string{AliceDID, BobDID},
			},
			signers: []string{AliceKey1, BobKey3},
			keys: map[string]KeyPair{
				AliceKey1: keys[AliceKey1],
				BobKey3:   keys[BobKey3],
			},
		},
		{
			valid: true,
			name:  "Valid: Full message works",
			keys: map[string]KeyPair{
				"did:cheqd:test:yyyyyyyyyyyyyyyy#key-1": GenerateKeyPair(),
				"did:cheqd:test:yyyyyyyyyyyyyyyy#key-2": GenerateKeyPair(),
				"did:cheqd:test:yyyyyyyyyyyyyyyy#key-3": GenerateKeyPair(),
				"did:cheqd:test:yyyyyyyyyyyyyyyy#key-4": GenerateKeyPair(),
				"did:cheqd:test:yyyyyyyyyyyyyyyy#key-5": GenerateKeyPair(),
				AliceKey1:                               keys[AliceKey1],
				BobKey1:                                 keys[BobKey1],
				BobKey2:                                 keys[BobKey2],
				BobKey3:                                 keys[BobKey3],
				CharlieKey1:                             keys[CharlieKey1],
				CharlieKey2:                             keys[CharlieKey2],
				CharlieKey3:                             keys[CharlieKey3],
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
				require.Equal(t, types.NormalizeIdentifiersList(tc.msg.Controller), did.Controller)
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

func TestHandler_Identifiers(t *testing.T) {
	keys := GenerateTestKeys()
	didPrefix := "did:cheqd:test:"
	cases := []struct {
		name     string
		id       string
		resultId string
	}{
		{
			name:     "Low Case UUID",
			id:       didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			resultId: didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		},
		{
			name:     "Upper Case UUID",
			id:       didPrefix + "A86F9CAE-0902-4A7C-A144-96B60CED2FC9",
			resultId: didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		},
		{
			name:     "Mixed Case UUID",
			id:       didPrefix + "A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
			resultId: didPrefix + "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		},
		{
			name:     "Indy-style id",
			id:       didPrefix + "MjYxNzYKMjYxNzYK",
			resultId: didPrefix + "MjYxNzYKMjYxNzYK",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setup := InitEnv(t, keys)

			key := GenerateKeyPair()
			keyAlias := tc.id + "#key1"
			publicKeyMultibase, _ := multibase.Encode(multibase.Base58BTC, key.PublicKey)
			signerKeys := map[string]ed25519.PrivateKey{keyAlias: key.PrivateKey}
			msg := &types.MsgCreateDidPayload{
				Id:             tc.id,
				Authentication: []string{keyAlias},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 keyAlias,
						Type:               Ed25519VerificationKey2020,
						Controller:         tc.id,
						PublicKeyMultibase: publicKeyMultibase,
					},
				},
			}

			did, err := setup.SendCreateDid(msg, signerKeys)

			require.Nil(t, err)
			require.Equal(t, tc.resultId, did.Id)
		})
	}
}
