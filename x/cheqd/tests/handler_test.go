package tests

import (
	"crypto/ed25519"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDID(t *testing.T) {
	setup := Setup()
	keys := setup.CreatePreparedDID()

	cases := []struct {
		valid   bool
		name    string
		keys    map[string]KeyPair
		signers []string
		msg     *v1.MsgCreateDidPayload
		errMsg  string
	}{
		{
			valid: true,
			name:  "Works",
			keys: map[string]KeyPair{
				"did:cheqd:test:123456qwertyui2#key-1": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:123456qwertyui2#key-1"},
			msg: &v1.MsgCreateDidPayload{
				Id:             "did:cheqd:test:123456qwertyui2",
				Authentication: []string{"did:cheqd:test:123456qwertyui2#key-1"},
				VerificationMethod: []*v1.VerificationMethod{
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
			name:  "Works with Key Agreement",
			keys: map[string]KeyPair{
				"did:cheqd:test:KeyAgreement#key-1": GenerateKeyPair(),
				AliceKey1:                           keys[AliceKey1],
			},
			signers: []string{AliceKey1},
			msg: &v1.MsgCreateDidPayload{
				Id:           "did:cheqd:test:KeyAgreement",
				KeyAgreement: []string{"did:cheqd:test:KeyAgreement#key-1"},
				Controller:   []string{AliceDID},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         "did:cheqd:test:KeyAgreement#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:KeyAgreement",
					},
				},
			},
		},
		{
			valid: true,
			name:  "Works with Assertion Method",
			keys: map[string]KeyPair{
				"did:cheqd:test:AssertionMethod#key-1": GenerateKeyPair(),
				AliceKey1:                              keys[AliceKey1],
			},
			signers: []string{AliceKey1},
			msg: &v1.MsgCreateDidPayload{
				Id:              "did:cheqd:test:AssertionMethod",
				AssertionMethod: []string{"did:cheqd:test:AssertionMethod#key-1"},
				Controller:      []string{AliceDID},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         "did:cheqd:test:AssertionMethod#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:AssertionMethod",
					},
				},
			},
		},
		{
			valid: true,
			name:  "Works with Capability Delegation",
			keys: map[string]KeyPair{
				"did:cheqd:test:CapabilityDelegation#key-1": GenerateKeyPair(),
				AliceKey1: keys[AliceKey1],
			},
			signers: []string{AliceKey1},
			msg: &v1.MsgCreateDidPayload{
				Id:                   "did:cheqd:test:CapabilityDelegation",
				CapabilityDelegation: []string{"did:cheqd:test:CapabilityDelegation#key-1"},
				Controller:           []string{AliceDID},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         "did:cheqd:test:CapabilityDelegation#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:CapabilityDelegation",
					},
				},
			},
		},
		{
			valid: true,
			name:  "Works with Capability Invocation",
			keys: map[string]KeyPair{
				"did:cheqd:test:CapabilityInvocation#key-1": GenerateKeyPair(),
				AliceKey1: keys[AliceKey1],
			},
			signers: []string{AliceKey1},
			msg: &v1.MsgCreateDidPayload{
				Id:                   "did:cheqd:test:CapabilityInvocation",
				CapabilityInvocation: []string{"did:cheqd:test:CapabilityInvocation#key-1"},
				Controller:           []string{AliceDID},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         "did:cheqd:test:CapabilityInvocation#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:CapabilityInvocation",
					},
				},
			},
		},
		{
			valid: true,
			name:  "With controller works",
			msg: &v1.MsgCreateDidPayload{
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
			msg: &v1.MsgCreateDidPayload{
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
				Service: []*v1.Service{
					{
						Id:              "did:cheqd:test:123456qwertyui#service-1",
						Type:            "DIDCommMessaging",
						ServiceEndpoint: "ServiceEndpoint",
					},
				},
				Controller: []string{"did:cheqd:test:123456qwertyui", AliceDID, BobDID, CharlieDID},
				VerificationMethod: []*v1.VerificationMethod{
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
			msg: &v1.MsgCreateDidPayload{
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
			msg: &v1.MsgCreateDidPayload{
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
			msg: &v1.MsgCreateDidPayload{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{AliceDID, BobDID},
			},
			errMsg: "At least one signature should be present: invalid signature detected",
		},
		{
			valid: false,
			name:  "Empty request",
			msg: &v1.MsgCreateDidPayload{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{AliceDID, BobDID},
			},
			errMsg: "At least one signature should be present: invalid signature detected",
		},
		{
			valid: false,
			name:  "Controller not found",
			msg: &v1.MsgCreateDidPayload{
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
			msg: &v1.MsgCreateDidPayload{
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
			msg: &v1.MsgCreateDidPayload{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{BobDID},
			},
			signers: []string{BobKey4},
			keys: map[string]KeyPair{
				BobKey4: keys[BobKey4],
			},
			errMsg: "did:cheqd:test:bob#key-4: verification method not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Second controller verification method not found",
			msg: &v1.MsgCreateDidPayload{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{AliceDID, BobDID, CharlieDID},
			},
			signers: []string{AliceKey1, BobKey4, CharlieKey3},
			keys: map[string]KeyPair{
				AliceKey1:   keys[AliceKey1],
				BobKey4:     keys[BobKey4],
				CharlieKey3: keys[CharlieKey3],
			},
			errMsg: "did:cheqd:test:bob#key-4: verification method not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "DID signed by wrong controller",
			msg: &v1.MsgCreateDidPayload{
				Id:             "did:cheqd:test:123456qwertyui",
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*v1.VerificationMethod{
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
			msg: &v1.MsgCreateDidPayload{
				Id:             "did:cheqd:test:123456qwertyui",
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*v1.VerificationMethod{
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
			errMsg: "did:cheqd:test:123456qwertyui#key-2: verification method not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Self-signature not found",
			msg: &v1.MsgCreateDidPayload{
				Id:             "did:cheqd:test:123456qwertyui",
				Controller:     []string{AliceDID, "did:cheqd:test:123456qwertyui"},
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*v1.VerificationMethod{
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
			errMsg: "did:cheqd:test:123456qwertyui#key-2: verification method not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "DID Doc already exists",
			keys: map[string]KeyPair{
				"did:cheqd:test:123456qwertyui#key-1": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:123456qwertyui#key-1"},
			msg: &v1.MsgCreateDidPayload{
				Id:             "did:cheqd:test:123456qwertyui",
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
				},
			},
			errMsg: "DID is already used by DIDDoc did:cheqd:test:123456qwertyui: DID Doc exists",
		},
		{
			valid: false,
			name:  "Verification Method ID doesnt match",
			msg: &v1.MsgCreateDidPayload{
				Id:             "did:cheqd:test:controller1",
				Controller:     []string{AliceDID, CharlieDID},
				Authentication: []string{"#key-1"},
				VerificationMethod: []*v1.VerificationMethod{
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
			msg: &v1.MsgCreateDidPayload{
				Id:             "did:cheqd:test:controller1",
				Controller:     []string{AliceDID, CharlieDID},
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*v1.VerificationMethod{
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
	keys := setup.CreatePreparedDID()

	cases := []struct {
		valid   bool
		name    string
		keys    map[string]KeyPair
		signers []string
		msg     *v1.MsgUpdateDidPayload
		errMsg  string
	}{
		{
			valid: true,
			name:  "Works",
			keys: map[string]KeyPair{
				AliceKey2: keys[AliceKey2],
			},
			signers: []string{AliceKey2},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid: false,
			name:  "Try to add controller without self-signature",
			keys: map[string]KeyPair{
				BobKey1: keys[BobKey1],
			},
			signers: []string{BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID},
				Authentication: []string{AliceKey1},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: "signature did:cheqd:test:alice not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Add controller and replace authentication without old signature do not work",
			keys: map[string]KeyPair{
				BobKey1:   keys[BobKey1],
				AliceKey1: keys[AliceKey1],
			},
			signers: []string{BobKey1, AliceKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID},
				Authentication: []string{AliceKey1},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: "did:cheqd:test:alice#key-1: verification method not found: invalid signature detected",
		},
		{
			valid: true,
			name:  "Add controller work",
			keys: map[string]KeyPair{
				BobKey1:   keys[BobKey1],
				AliceKey2: keys[AliceKey2],
			},
			signers: []string{BobKey1, AliceKey2},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid: false,
			name:  "Add controller work without signature do not work",
			keys: map[string]KeyPair{
				BobKey1:   keys[BobKey1],
				AliceKey2: keys[AliceKey2],
			},
			signers: []string{BobKey1, AliceKey2},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: "signature did:cheqd:test:charlie not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Replace controller work without new signature do not work",
			keys: map[string]KeyPair{
				BobKey1:   keys[BobKey1],
				AliceKey2: keys[AliceKey2],
			},
			signers: []string{BobKey1, AliceKey2},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: "signature did:cheqd:test:charlie not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Replace controller without old signature do not work",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2, CharlieKey3},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: "signature did:cheqd:test:bob not found: invalid signature detected",
		},
		{
			valid: true,
			name:  "Replace controller work",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2, CharlieKey3, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid: true,
			name:  "Add second controller works",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2, CharlieKey3, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid: true,
			name:  "Add verification method without signature controller work",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey1:   keys[AliceKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{CharlieKey3, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				KeyAgreement:   []string{AliceKey1},
				VerificationMethod: []*v1.VerificationMethod{
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
			valid: false,
			name:  "Remove verification method without signature controller do not work",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey1:   keys[AliceKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{CharlieKey3, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: "signature did:cheqd:test:alice not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Remove verification method wrong authentication detected",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey1:   keys[AliceKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey1, CharlieKey3, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: "did:cheqd:test:alice#key-1: verification method not found: invalid signature detected",
		},
		{
			valid: true,
			name:  "Add second authentication works",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey1:   keys[AliceKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2, CharlieKey3, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey1, AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
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
			valid: false,
			name:  "Remove self authentication without signature do not work",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{CharlieKey3, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
			errMsg: "signature did:cheqd:test:alice not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Change self controller verification without signature do not work",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{CharlieKey3, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey1, AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
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
			errMsg: "signature did:cheqd:test:alice not found: invalid signature detected",
		},
		{
			valid: true,
			name:  "Remove self authentication works",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2, CharlieKey3, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{BobDID, CharlieDID},
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
		{
			valid: false,
			name:  "Change controller to self without old controllers signatures does not work",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
			errMsg: "signature did:cheqd:test:bob not found: invalid signature detected",
		},
		{
			valid: true,
			name:  "Change controller to self works",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2, CharlieKey3, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
		{
			valid: false,
			name:  "Change verification method controller without old signature",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2, CharlieKey3},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: CharlieDID,
					},
				},
			},
			errMsg: "signature did:cheqd:test:bob not found: invalid signature detected",
		},
		{
			valid: false,
			name:  "Change verification method controller without new signature",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2, BobKey1},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: CharlieDID,
					},
				},
			},
			errMsg: "signature did:cheqd:test:charlie not found: invalid signature detected",
		},
		{
			valid: true,
			name:  "Change verification method controller",
			keys: map[string]KeyPair{
				BobKey1:     keys[BobKey1],
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2, BobKey1, CharlieKey3},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: CharlieDID,
					},
				},
			},
		},
		{
			valid: false,
			name:  "Change to self verification method without controller signature",
			keys: map[string]KeyPair{
				AliceKey2: keys[AliceKey2],
			},
			signers: []string{AliceKey2},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
			errMsg: "signature did:cheqd:test:charlie not found: invalid signature detected",
		},
		{
			valid: true,
			name:  "Change to self verification method without controller signature",
			keys: map[string]KeyPair{
				AliceKey2:   keys[AliceKey2],
				CharlieKey3: keys[CharlieKey3],
			},
			signers: []string{AliceKey2, CharlieKey3},
			msg: &v1.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey2},
				VerificationMethod: []*v1.VerificationMethod{
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
				vm.PublicKeyMultibase = "z" + base58.Encode(tc.keys[vm.Id].PublicKey)
			}

			signerKeys := map[string]ed25519.PrivateKey{}
			for _, signer := range tc.signers {
				signerKeys[signer] = tc.keys[signer].PrivateKey
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

	_, _, _ = setup.InitDid("did:cheqd:test:alice")
	_, _, err := setup.InitDid("did:cheqd:test:alice")

	require.Error(t, err)
	require.Equal(t, "DID is already used by DIDDoc did:cheqd:test:alice: DID Doc exists", err.Error())
}
