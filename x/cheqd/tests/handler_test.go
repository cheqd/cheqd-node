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
	setup := Setup()
	keys, err := setup.CreateTestDIDs()
	require.NoError(t, err)

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
			name:  "Works",
			keys: map[string]KeyPair{
				"did:cheqd:test:0123456qwertyui2#key-1": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:0123456qwertyui2#key-1"},
			msg: &types.MsgCreateDidPayload{
				Id:             "did:cheqd:test:0123456qwertyui2",
				Authentication: []string{"did:cheqd:test:0123456qwertyui2#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:0123456qwertyui2#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:0123456qwertyui2",
					},
				},
			},
		},
		{
			valid: true,
			name:  "Works with Key Agreement",
			keys: map[string]KeyPair{
				"did:cheqd:test:0000KeyAgreement#key-1": GenerateKeyPair(),
				AliceKey1:                           keys[AliceKey1],
			},
			signers: []string{AliceKey1},
			msg: &types.MsgCreateDidPayload{
				Id:           "did:cheqd:test:0000KeyAgreement",
				KeyAgreement: []string{"did:cheqd:test:0000KeyAgreement#key-1"},
				Controller:   []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:0000KeyAgreement#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:0000KeyAgreement",
					},
				},
			},
		},
		{
			valid: true,
			name:  "Works with Assertion Method",
			keys: map[string]KeyPair{
				"did:cheqd:test:0AssertionMethod#key-1": GenerateKeyPair(),
				AliceKey1:                              keys[AliceKey1],
			},
			signers: []string{AliceKey1},
			msg: &types.MsgCreateDidPayload{
				Id:              "did:cheqd:test:0AssertionMethod",
				AssertionMethod: []string{"did:cheqd:test:0AssertionMethod#key-1"},
				Controller:      []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:0AssertionMethod#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:0AssertionMethod",
					},
				},
			},
		},
		{
			valid: true,
			name:  "Works with Capability Delegation",
			keys: map[string]KeyPair{
				"did:cheqd:test:000000000000CapabilityDelegation#key-1": GenerateKeyPair(),
				AliceKey1: keys[AliceKey1],
			},
			signers: []string{AliceKey1},
			msg: &types.MsgCreateDidPayload{
				Id:                   "did:cheqd:test:000000000000CapabilityDelegation",
				CapabilityDelegation: []string{"did:cheqd:test:000000000000CapabilityDelegation#key-1"},
				Controller:           []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:000000000000CapabilityDelegation#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:000000000000CapabilityDelegation",
					},
				},
			},
		},
		{
			valid: true,
			name:  "Works with Capability Invocation",
			keys: map[string]KeyPair{
				"did:cheqd:test:000000000000CapabilityInvocation#key-1": GenerateKeyPair(),
				AliceKey1: keys[AliceKey1],
			},
			signers: []string{AliceKey1},
			msg: &types.MsgCreateDidPayload{
				Id:                   "did:cheqd:test:000000000000CapabilityInvocation",
				CapabilityInvocation: []string{"did:cheqd:test:000000000000CapabilityInvocation#key-1"},
				Controller:           []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:000000000000CapabilityInvocation#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:000000000000CapabilityInvocation",
					},
				},
			},
		},
		{
			valid: true,
			name:  "With controller works",
			msg: &types.MsgCreateDidPayload{
				Id:         "did:cheqd:test:00000controller1",
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
				"did:cheqd:test:00123456qwertyui#key-1": GenerateKeyPair(),
				"did:cheqd:test:00123456qwertyui#key-2": GenerateKeyPair(),
				"did:cheqd:test:00123456qwertyui#key-3": GenerateKeyPair(),
				"did:cheqd:test:00123456qwertyui#key-4": GenerateKeyPair(),
				"did:cheqd:test:00123456qwertyui#key-5": GenerateKeyPair(),
				AliceKey1:                             keys[AliceKey1],
				BobKey1:                               keys[BobKey1],
				BobKey2:                               keys[BobKey2],
				BobKey3:                               keys[BobKey3],
				CharlieKey1:                           keys[CharlieKey1],
				CharlieKey2:                           keys[CharlieKey2],
				CharlieKey3:                           keys[CharlieKey3],
			},
			signers: []string{
				"did:cheqd:test:00123456qwertyui#key-1",
				"did:cheqd:test:00123456qwertyui#key-5",
				AliceKey1,
				BobKey1,
				BobKey2,
				BobKey3,
				CharlieKey1,
				CharlieKey2,
				CharlieKey3,
			},
			msg: &types.MsgCreateDidPayload{
				Id: "did:cheqd:test:00123456qwertyui",
				Authentication: []string{
					"did:cheqd:test:00123456qwertyui#key-1",
					"did:cheqd:test:00123456qwertyui#key-5",
				},
				Context:              []string{"abc", "de"},
				CapabilityInvocation: []string{"did:cheqd:test:00123456qwertyui#key-2"},
				CapabilityDelegation: []string{"did:cheqd:test:00123456qwertyui#key-3"},
				KeyAgreement:         []string{"did:cheqd:test:00123456qwertyui#key-4"},
				AlsoKnownAs:          []string{"did:cheqd:test:000123456eqweqwe"},
				Service: []*types.Service{
					{
						Id:              "did:cheqd:test:00123456qwertyui#service-1",
						Type:            "DIDCommMessaging",
						ServiceEndpoint: "ServiceEndpoint",
					},
				},
				Controller: []string{"did:cheqd:test:00123456qwertyui", AliceDID, BobDID, CharlieDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-2",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-3",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-4",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-5",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
				},
			},
		},
		{
			valid: false,
			name:  "Second controller did not sign request",
			msg: &types.MsgCreateDidPayload{
				Id:         "did:cheqd:test:00000controller1",
				Controller: []string{AliceDID, BobDID},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[AliceKey1],
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", BobDID),
		},
		{
			valid: false,
			name:  "Bad request",
			msg: &types.MsgCreateDidPayload{
				Id: "did:cheqd:test:00000controller1",
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
			msg: &types.MsgCreateDidPayload{
				Id:         "did:cheqd:test:00000controller1",
				Controller: []string{AliceDID, BobDID},
			},
			errMsg: "At least one signature should be present: invalid signature detected",
		},
		{
			valid: false,
			name:  "Empty request",
			msg: &types.MsgCreateDidPayload{
				Id:         "did:cheqd:test:00000controller1",
				Controller: []string{AliceDID, BobDID},
			},
			errMsg: "At least one signature should be present: invalid signature detected",
		},
		{
			valid: false,
			name:  "Controller not found",
			msg: &types.MsgCreateDidPayload{
				Id:         "did:cheqd:test:00000controller1",
				Controller: []string{AliceDID, "did:cheqd:test:00000000notfound"},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[AliceKey1],
			},
			errMsg: "did:cheqd:test:00000000notfound: DID Doc not found",
		},
		{
			valid: false,
			name:  "Wrong signature",
			msg: &types.MsgCreateDidPayload{
				Id:         "did:cheqd:test:00000controller1",
				Controller: []string{AliceDID},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[BobKey1],
			},
			errMsg: fmt.Sprintf("%v: invalid signature detected", AliceDID),
		},
		{
			valid: false,
			name:  "Controller verification method not found",
			msg: &types.MsgCreateDidPayload{
				Id:         "did:cheqd:test:00000controller1",
				Controller: []string{BobDID},
			},
			signers: []string{BobKey4},
			keys: map[string]KeyPair{
				BobKey4: keys[BobKey4],
			},
			errMsg: "did:cheqd:test:0000000000000bob#key-4: verification method not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Second controller verification method not found",
			msg: &types.MsgCreateDidPayload{
				Id:         "did:cheqd:test:00000controller1",
				Controller: []string{AliceDID, BobDID, CharlieDID},
			},
			signers: []string{AliceKey1, BobKey4, CharlieKey3},
			keys: map[string]KeyPair{
				AliceKey1:   keys[AliceKey1],
				BobKey4:     keys[BobKey4],
				CharlieKey3: keys[CharlieKey3],
			},
			errMsg: "did:cheqd:test:0000000000000bob#key-4: verification method not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "DID signed by wrong controller",
			msg: &types.MsgCreateDidPayload{
				Id:             "did:cheqd:test:00123456qwertyui",
				Authentication: []string{"did:cheqd:test:00123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
				},
			},
			signers: []string{AliceKey1},
			keys: map[string]KeyPair{
				AliceKey1: keys[AliceKey1],
			},
			errMsg: "signature did:cheqd:test:00123456qwertyui not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "DID self-signed by not existing verification method",
			msg: &types.MsgCreateDidPayload{
				Id:             "did:cheqd:test:00123456qwertyui",
				Authentication: []string{"did:cheqd:test:00123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
				},
			},
			signers: []string{"did:cheqd:test:00123456qwertyui#key-2"},
			keys: map[string]KeyPair{
				"did:cheqd:test:00123456qwertyui#key-2": GenerateKeyPair(),
			},
			errMsg: "did:cheqd:test:00123456qwertyui#key-2: verification method not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Self-signature not found",
			msg: &types.MsgCreateDidPayload{
				Id:             "did:cheqd:test:00123456qwertyui",
				Controller:     []string{AliceDID, "did:cheqd:test:00123456qwertyui"},
				Authentication: []string{"did:cheqd:test:00123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
				},
			},
			signers: []string{AliceKey1, "did:cheqd:test:00123456qwertyui#key-2"},
			keys: map[string]KeyPair{
				AliceKey1:                             keys[AliceKey1],
				"did:cheqd:test:00123456qwertyui#key-2": GenerateKeyPair(),
			},
			errMsg: "did:cheqd:test:00123456qwertyui#key-2: verification method not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "DID Doc already exists",
			keys: map[string]KeyPair{
				"did:cheqd:test:00123456qwertyui#key-1": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:00123456qwertyui#key-1"},
			msg: &types.MsgCreateDidPayload{
				Id:             "did:cheqd:test:00123456qwertyui",
				Authentication: []string{"did:cheqd:test:00123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
				},
			},
			errMsg: "DID is already used by DIDDoc did:cheqd:test:00123456qwertyui: DID Doc exists",
		},
		{
			valid: false,
			name:  "Verification Method ID doesnt match",
			msg: &types.MsgCreateDidPayload{
				Id:             "did:cheqd:test:00000controller1",
				Controller:     []string{AliceDID, CharlieDID},
				Authentication: []string{"#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
				},
			},
			signers: []string{AliceKey1, CharlieKey3},
			keys: map[string]KeyPair{
				AliceKey1:   keys[AliceKey1],
				CharlieKey3: keys[CharlieKey3],
			},
			errMsg: "did:cheqd:test:00123456qwertyui#key-1 not belong did:cheqd:test:00000controller1 DID Doc: invalid verification method",
		},
		{
			valid: false,
			name:  "Full Verification Method ID doesnt match",
			msg: &types.MsgCreateDidPayload{
				Id:             "did:cheqd:test:00000controller1",
				Controller:     []string{AliceDID, CharlieDID},
				Authentication: []string{"did:cheqd:test:00123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:00123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:00123456qwertyui",
					},
				},
			},
			signers: []string{AliceKey1, CharlieKey3},
			keys: map[string]KeyPair{
				AliceKey1:   keys[AliceKey1],
				CharlieKey3: keys[CharlieKey3],
			},
			errMsg: "did:cheqd:test:00123456qwertyui#key-1 not belong did:cheqd:test:00000controller1 DID Doc: invalid verification method",
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
				require.Equal(t, tc.msg.Id, did.Id)
				require.Equal(t, tc.msg.Controller, did.Controller)
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

func TestUpdateDid(t *testing.T) {
	setup := Setup()
	keys, err := setup.CreateTestDIDs()
	require.NoError(t, err)

	cases := []struct {
		valid   bool
		name    string
		signers []string
		msg     *types.MsgUpdateDidPayload
		errMsg  string
	}{
		{
			valid:   true,
			name:    "Key rotation works",
			signers: []string{AliceKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid:   false,
			name:    "Try to add controller without self-signature",
			signers: []string{BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID},
				Authentication: []string{AliceKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", AliceDID),
		},
		{
			valid:   false,
			name:    "Add controller and replace authentication without old signature do not work",
			signers: []string{BobKey1, AliceKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID},
				Authentication: []string{AliceKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("%v: verification method not found: invalid signature detected", AliceKey1),
		},
		{
			valid:   true,
			name:    "Add controller work",
			signers: []string{BobKey1, AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid:   true,
			name:    "Add controller without signature work (signatures of old controllers are present)",
			signers: []string{BobKey1, AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", CharlieDID),
		},
		{
			valid:   false,
			name:    "Replace controller work without new signature do not work",
			signers: []string{BobKey1, AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", CharlieDID),
		},
		{
			valid:   false,
			name:    "Replace controller without old signature do not work",
			signers: []string{AliceKey2, CharlieKey3},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", BobDID),
		},
		{
			valid:   true,
			name:    "Replace controller work",
			signers: []string{AliceKey2, CharlieKey3, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid:   true,
			name:    "Add second controller works",
			signers: []string{AliceKey2, CharlieKey3, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid:   true,
			name:    "Add verification method without signature controller work",
			signers: []string{CharlieKey3, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				KeyAgreement:   []string{AliceKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid:   false,
			name:    "Remove verification method without signature controller do not work",
			signers: []string{CharlieKey3, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", AliceDID),
		},
		{
			valid:   false,
			name:    "Remove verification method wrong authentication detected",
			signers: []string{AliceKey1, CharlieKey3, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("%v: verification method not found: invalid signature detected", AliceKey1),
		},
		{
			valid:   true,
			name:    "Add second authentication works",
			signers: []string{AliceKey2, CharlieKey3, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey1, AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
		{
			valid:   false,
			name:    "Remove self authentication without signature do not work",
			signers: []string{CharlieKey3, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", AliceDID),
		},
		{
			valid:   false,
			name:    "Change self controller verification without signature do not work",
			signers: []string{CharlieKey3, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey1, AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: CharlieDID,
					},
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", AliceDID),
		},
		{
			valid:   true,
			signers: []string{AliceKey2, CharlieKey3, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
		{
			valid:   false,
			name:    "Change controller to self without old controllers signatures does not work",
			signers: []string{AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", BobDID),
		},
		{
			valid:   true,
			name:    "Change controller to self works",
			signers: []string{AliceKey2, CharlieKey3, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
		{
			valid:   false,
			name:    "Change verification method controller without old signature",
			signers: []string{AliceKey2, CharlieKey3},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: CharlieDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", BobDID),
		},
		{
			valid:   false,
			name:    "Change verification method controller without new signature",
			signers: []string{AliceKey2, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: CharlieDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", CharlieDID),
		},
		{
			valid:   true,
			name:    "Change verification method controller",
			signers: []string{AliceKey2, BobKey1, CharlieKey3},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: CharlieDID,
					},
				},
			},
		},
		{
			valid:   false,
			name:    "Change to self verification method without controller signature",
			signers: []string{AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("signature %v not found: invalid signature detected", CharlieDID),
		},
		{
			valid:   true,
			name:    "Change to self verification method without controller signature",
			signers: []string{AliceKey2, CharlieKey3},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg

			for _, vm := range msg.VerificationMethod {
				encoded, err := multibase.Encode(multibase.Base58BTC, keys[vm.Id].PublicKey)
				require.NoError(t, err)
				vm.PublicKeyMultibase = encoded
			}

			signerKeys := map[string]ed25519.PrivateKey{}
			for _, signer := range tc.signers {
				signerKeys[signer] = keys[signer].PrivateKey
			}

			did, err := setup.SendUpdateDid(msg, signerKeys)

			if tc.valid {
				require.Nil(t, err)
				require.Equal(t, tc.msg.Id, did.Id)
				require.Equal(t, tc.msg.Controller, did.Controller)
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
	require.Equal(t, fmt.Sprintf("DID is already used by DIDDoc %v: DID Doc exists", AliceDID), err.Error())
}
