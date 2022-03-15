package tests

import (
	"crypto/ed25519"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/multiformats/go-multibase"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateDid(t *testing.T) {
	keys := map[string]KeyPair{
		AliceKey1:    GenerateKeyPair(),
		AliceKey2:    GenerateKeyPair(),
		BobKey1:      GenerateKeyPair(),
		BobKey2:      GenerateKeyPair(),
		BobKey3:      GenerateKeyPair(),
		BobKey4:      GenerateKeyPair(),
		CharlieKey1:  GenerateKeyPair(),
		CharlieKey2:  GenerateKeyPair(),
		CharlieKey3:  GenerateKeyPair(),
		ImposterKey1: GenerateKeyPair(),
	}

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
			signers: []string{AliceKey1, AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id:             AliceDID,
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
			errMsg: "signature did:cheqd:test:alice not found: invalid signature detected",
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
			errMsg: "did:cheqd:test:alice#key-1: verification method not found: invalid signature detected",
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
			errMsg: "signature did:cheqd:test:charlie not found: invalid signature detected",
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
			errMsg: "signature did:cheqd:test:charlie not found: invalid signature detected",
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
			errMsg: "signature did:cheqd:test:bob not found: invalid signature detected",
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
			errMsg: "signature did:cheqd:test:alice not found: invalid signature detected",
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
			errMsg: "did:cheqd:test:alice#key-1: verification method not found: invalid signature detected",
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
			errMsg: "signature did:cheqd:test:alice not found: invalid signature detected",
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
			errMsg: "signature did:cheqd:test:alice not found: invalid signature detected",
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
			errMsg: "signature did:cheqd:test:bob not found: invalid signature detected",
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
			errMsg: "signature did:cheqd:test:bob not found: invalid signature detected",
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
			errMsg: "signature did:cheqd:test:charlie not found: invalid signature detected",
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
			errMsg: "signature did:cheqd:test:charlie not found: invalid signature detected",
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
			setup := InitEnv(t, keys)
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
