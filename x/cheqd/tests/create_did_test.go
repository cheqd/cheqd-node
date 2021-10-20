package tests

import (
	"crypto/ed25519"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateDID(t *testing.T) {
	prefilledDids := []struct {
		keys    map[string]KeyPair
		signers []string
		msg     *types.MsgCreateDid
	}{
		{
			keys: map[string]KeyPair{
				AliceKey1: GenerateKeyPair(),
				AliceKey2: GenerateKeyPair(),
			},
			signers: []string{AliceKey1},
			msg: &types.MsgCreateDid{
				Id:             AliceDID,
				Authentication: []string{AliceKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			keys: map[string]KeyPair{
				BobKey1: GenerateKeyPair(),
				BobKey2: GenerateKeyPair(),
				BobKey3: GenerateKeyPair(),
				BobKey4: GenerateKeyPair(),
			},
			signers: []string{BobKey2},
			msg: &types.MsgCreateDid{
				Id: BobDID,
				Authentication: []string{
					BobKey1,
					BobKey2,
					BobKey3,
				},
				CapabilityDelegation: []string{
					BobKey4,
				},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         BobKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey3,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey4,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
		{
			keys: map[string]KeyPair{
				CharlieKey1: GenerateKeyPair(),
				CharlieKey2: GenerateKeyPair(),
				CharlieKey3: GenerateKeyPair(),
			},
			signers: []string{CharlieKey2},
			msg: &types.MsgCreateDid{
				Id: CharlieDID,
				Authentication: []string{
					CharlieKey1,
					CharlieKey2,
					CharlieKey3,
				},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         CharlieKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         CharlieKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         CharlieKey3,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
	}

	setup := Setup()
	keys := map[string]KeyPair{}

	for _, prefilled := range prefilledDids {
		msg := prefilled.msg

		for _, vm := range msg.VerificationMethod {
			vm.PublicKeyMultibase = "z" + base58.Encode(prefilled.keys[vm.Id].PublicKey)
		}

		signerKeys := map[string]ed25519.PrivateKey{}
		for _, signer := range prefilled.signers {
			signerKeys[signer] = prefilled.keys[signer].PrivateKey
		}

		for keyId, key := range prefilled.keys {
			keys[keyId] = key
		}

		_, _ = setup.SendCreateDid(msg, signerKeys)
	}

	cases := []struct {
		valid   bool
		name    string
		keys    map[string]KeyPair
		signers []string
		msg     *types.MsgCreateDid
		errMsg  string
	}{
		{
			valid: true,
			name:  "Works",
			keys: map[string]KeyPair{
				"did:cheqd:test:123456qwertyui2#key-1": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:123456qwertyui2#key-1"},
			msg: &types.MsgCreateDid{
				Id:             "did:cheqd:test:123456qwertyui2",
				Authentication: []string{"did:cheqd:test:123456qwertyui2#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui2#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui2",
					},
				},
			},
		},
		{
			valid: true,
			name:  "With controller works",
			msg: &types.MsgCreateDid{
				Id:         "did:cheqd:test:controller1",
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
			name:  "Full message works",
			keys: map[string]KeyPair{
				"did:cheqd:test:123456qwertyui#key-1": GenerateKeyPair(),
				"did:cheqd:test:123456qwertyui#key-2": GenerateKeyPair(),
				"did:cheqd:test:123456qwertyui#key-3": GenerateKeyPair(),
				"did:cheqd:test:123456qwertyui#key-4": GenerateKeyPair(),
				"did:cheqd:test:123456qwertyui#key-5": GenerateKeyPair(),
				AliceKey1:                             keys[AliceKey1],
				BobKey1:                               keys[BobKey1],
				BobKey2:                               keys[BobKey2],
				BobKey3:                               keys[BobKey3],
				CharlieKey1:                           keys[CharlieKey1],
				CharlieKey2:                           keys[CharlieKey2],
				CharlieKey3:                           keys[CharlieKey3],
			},
			signers: []string{
				"did:cheqd:test:123456qwertyui#key-1",
				"did:cheqd:test:123456qwertyui#key-5",
				AliceKey1,
				BobKey1,
				BobKey2,
				BobKey3,
				CharlieKey1,
				CharlieKey2,
				CharlieKey3,
			},
			msg: &types.MsgCreateDid{
				Id: "did:cheqd:test:123456qwertyui",
				Authentication: []string{
					"did:cheqd:test:123456qwertyui#key-1",
					"did:cheqd:test:123456qwertyui#key-5",
				},
				Context:              []string{"abc", "de"},
				CapabilityInvocation: []string{"did:cheqd:test:123456qwertyui#key-2"},
				CapabilityDelegation: []string{"did:cheqd:test:123456qwertyui#key-3"},
				KeyAgreement:         []string{"did:cheqd:test:123456qwertyui#key-4"},
				AlsoKnownAs:          []string{"did:cheqd:test:123456eqweqwe"},
				Service: []*types.DidService{
					{
						Id:              "did:cheqd:test:123456qwertyui#service-1",
						Type:            "DIDCommMessaging",
						ServiceEndpoint: "ServiceEndpoint",
					},
				},
				Controller: []string{"did:cheqd:test:123456qwertyui", AliceDID, BobDID, CharlieDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
					{
						Id:         "did:cheqd:test:123456qwertyui#key-2",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
					{
						Id:         "did:cheqd:test:123456qwertyui#key-3",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
					{
						Id:         "did:cheqd:test:123456qwertyui#key-4",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
					{
						Id:         "did:cheqd:test:123456qwertyui#key-5",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
				},
			},
		},
		{
			valid: false,
			name:  "Second controller did not sign request",
			msg: &types.MsgCreateDid{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{AliceDID, BobDID},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[AliceKey1],
			},
			errMsg: "signature did:cheqd:test:bob not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Bad request",
			msg: &types.MsgCreateDid{
				Id: "did:cheqd:test:controller1",
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[AliceKey1],
			},
			errMsg: "The message must contain either a Controller or a Authentication: bad request",
		},
		{
			valid: false,
			name:  "No signature",
			msg: &types.MsgCreateDid{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{AliceDID, BobDID},
			},
			errMsg: "Signatures: is required",
		},
		{
			valid: false,
			name:  "Empty request",
			msg: &types.MsgCreateDid{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{AliceDID, BobDID},
			},
			errMsg: "Signatures: is required",
		},
		{
			valid: false,
			name:  "Controller not found",
			msg: &types.MsgCreateDid{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{AliceDID, "did:cheqd:test:notfound"},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[AliceKey1],
			},
			errMsg: "did:cheqd:test:notfound: DID Doc not found",
		},
		{
			valid: false,
			name:  "Wrong signature",
			msg: &types.MsgCreateDid{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{AliceDID},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[BobKey1],
			},
			errMsg: "did:cheqd:test:alice: invalid signature detected",
		},
		{
			valid: false,
			name:  "Controller verification method not found",
			msg: &types.MsgCreateDid{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{BobDID},
			},
			signers: []string{BobKey4},
			keys: map[string]KeyPair{
				BobKey4: keys[BobKey4],
			},
			errMsg: "Authentication did:cheqd:test:bob#key-4 not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Second controller verification method not found",
			msg: &types.MsgCreateDid{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{AliceDID, BobDID, CharlieDID},
			},
			signers: []string{AliceKey1, BobKey4, CharlieKey3},
			keys: map[string]KeyPair{
				AliceKey1:   keys[AliceKey1],
				BobKey4:     keys[BobKey4],
				CharlieKey3: keys[CharlieKey3],
			},
			errMsg: "Authentication did:cheqd:test:bob#key-4 not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "DID signed by wrong controller",
			msg: &types.MsgCreateDid{
				Id:             "did:cheqd:test:123456qwertyui",
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
				},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[AliceKey1],
			},
			errMsg: "signature did:cheqd:test:123456qwertyui not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "DID self-signed by not existing verification method",
			msg: &types.MsgCreateDid{
				Id:             "did:cheqd:test:123456qwertyui",
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
				},
			},
			signers: []string{"did:cheqd:test:123456qwertyui#key-2"},
			keys: map[string]KeyPair{
				"did:cheqd:test:123456qwertyui#key-2": GenerateKeyPair(),
			},
			errMsg: "Authentication did:cheqd:test:123456qwertyui#key-2 not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Self-signature not found",
			msg: &types.MsgCreateDid{
				Id:             "did:cheqd:test:123456qwertyui",
				Controller:     []string{AliceDID, "did:cheqd:test:123456qwertyui"},
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
				},
			},
			signers: []string{AliceKey1, "did:cheqd:test:123456qwertyui#key-2"},
			keys: map[string]KeyPair{
				AliceKey1:                             keys[AliceKey1],
				"did:cheqd:test:123456qwertyui#key-2": GenerateKeyPair(),
			},
			errMsg: "Authentication did:cheqd:test:123456qwertyui#key-2 not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "DID Doc already exists",
			keys: map[string]KeyPair{
				"did:cheqd:test:123456qwertyui#key-1": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:123456qwertyui#key-1"},
			msg: &types.MsgCreateDid{
				Id:             "did:cheqd:test:123456qwertyui",
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
				},
			},
			errMsg: "DID DOC already exists for DID did:cheqd:test:123456qwertyui: DID Doc exists",
		},
		{
			valid: false,
			name:  "Verification Method ID doesnt match",
			msg: &types.MsgCreateDid{
				Id:             "did:cheqd:test:controller1",
				Controller:     []string{AliceDID, CharlieDID},
				Authentication: []string{"#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
				},
			},
			signers: []string{AliceKey1, CharlieKey3},
			keys: map[string]KeyPair{
				AliceKey1:   keys[AliceKey1],
				CharlieKey3: keys[CharlieKey3],
			},
			errMsg: "did:cheqd:test:123456qwertyui#key-1 not belong did:cheqd:test:controller1 DID Doc: invalid verification method",
		},
		{
			valid: false,
			name:  "Full Verification Method ID doesnt match",
			msg: &types.MsgCreateDid{
				Id:             "did:cheqd:test:controller1",
				Controller:     []string{AliceDID, CharlieDID},
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
				},
			},
			signers: []string{AliceKey1, CharlieKey3},
			keys: map[string]KeyPair{
				AliceKey1:   keys[AliceKey1],
				CharlieKey3: keys[CharlieKey3],
			},
			errMsg: "did:cheqd:test:123456qwertyui#key-1 not belong did:cheqd:test:controller1 DID Doc: invalid verification method",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg

			for _, vm := range msg.VerificationMethod {
				vm.PublicKeyMultibase = "z" + base58.Encode(tc.keys[vm.Id].PublicKey)
			}

			signerKeys := map[string]ed25519.PrivateKey{}
			for _, signer := range tc.signers {
				signerKeys[signer] = tc.keys[signer].PrivateKey
			}

			did, err := setup.SendCreateDid(msg, signerKeys)

			if tc.valid {
				require.Nil(t, err)
				require.Equal(t, did.Id, did.Id)
				require.Equal(t, did.Controller, did.Controller)
				require.Equal(t, did.VerificationMethod, did.VerificationMethod)
				require.Equal(t, did.Authentication, did.Authentication)
				require.Equal(t, did.AssertionMethod, did.AssertionMethod)
				require.Equal(t, did.CapabilityInvocation, did.CapabilityInvocation)
				require.Equal(t, did.CapabilityDelegation, did.CapabilityDelegation)
				require.Equal(t, did.KeyAgreement, did.KeyAgreement)
				require.Equal(t, did.AlsoKnownAs, did.AlsoKnownAs)
				require.Equal(t, did.Service, did.Service)
				require.Equal(t, did.Context, did.Context)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
