package tests

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/multiformats/go-multibase"
	"github.com/stretchr/testify/require"
)

func TestUpdateDid(t *testing.T) {
	var err error
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
		CharlieKey4:  GenerateKeyPair(),
		ImposterKey1: GenerateKeyPair(),
	}

	cases := []struct {
		valid      bool
		name       string
		signerKeys []SignerKey
		signers    []string
		msg        *types.MsgUpdateDidPayload
		errMsg     string
	}{
		{
			valid: true,
			name:  "Valid: Key rotation works",
			signerKeys: []SignerKey{
				{
					signer: AliceKey1,
					key:    keys[AliceKey1].PrivateKey,
				},
				{
					signer: AliceKey1,
					key:    keys[AliceKey2].PrivateKey,
				},
			},
			msg: &types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 AliceKey1,
						Type:               Ed25519VerificationKey2020,
						Controller:         AliceDID,
						PublicKeyMultibase: "z" + base58.Encode(keys[AliceKey2].PublicKey),
					},
				},
			},
		},
		// VM and Controller replacing tests
		{
			valid:   false,
			name:    "Not Valid: replacing controller and Verification method ID does not work without new sign",
			signers: []string{AliceKey2, BobKey1, AliceKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{CharlieDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", utils.NormalizeIdentifier(CharlieDID)),
		},
		{
			valid:   true,
			name:    "Valid: replacing controller and Verification method ID works with all signatures",
			signers: []string{AliceKey1, CharlieKey1, AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{CharlieDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", utils.NormalizeIdentifier(CharlieDID)),
		},
		// Verification method's tests
		// cases:
		// - replacing VM controller works
		// - replacing VM controller does not work without new signature
		// - replacing VM controller does not work without old signature     ??????
		// - replacing VM doesn't work without new signature
		// - replacing VM doesn't work without old signature
		// - replacing VM works with all signatures
		// --- adding new VM works
		// --- adding new VM without new signature
		// --- adding new VM without old signature
		{
			valid:   true,
			name:    "Valid: Replacing VM controller works with one signature",
			signers: []string{AliceKey1, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: BobDID,
					},
				},
			},
		},
		{
			valid:   false,
			name:    "Not Valid: Replacing VM controller does not work without new signature",
			signers: []string{AliceKey1},
			msg: &types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: BobDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID),
		},
		{
			valid:   false,
			name:    "Not Valid: Replacing VM does not work without new signature",
			signers: []string{AliceKey1},
			msg: &types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one valid signature by %s (new version): invalid signature detected", AliceDID),
		},
		{
			valid:   false,
			name:    "Not Valid: Replacing VM does not work without old signature",
			signers: []string{AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", AliceDID),
		},
		{
			valid:   true,
			name:    "Not Valid: Replacing VM works with all signatures",
			signers: []string{AliceKey1, AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
		},
		// Adding VM
		{
			valid:   true,
			name:    "Valid: Adding another verification method",
			signers: []string{AliceKey1, BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
					{
						Id:         AliceKey2,
						Type:       Ed25519VerificationKey2020,
						Controller: BobDID,
					},
				},
			},
		},
		{
			valid:   false,
			name:    "Not Valid: Adding another verification method without new sign",
			signers: []string{AliceKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
					{
						Id:         AliceKey2,
						Type:       Ed25519VerificationKey2020,
						Controller: BobDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID),
		},
		{
			valid:   false,
			name:    "Not Valid: Adding another verification method without old sign",
			signers: []string{AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
					{
						Id:         AliceKey2,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", AliceDID),
		},

		// Controller's tests
		// cases:
		// - replacing Controller works with all signatures
		// - replacing Controller doesn't work without old signature
		// - replacing Controller doesn't work without new signature
		// --- adding Controller works with all signatures
		// --- adding Controller doesn't work without old signature
		// --- adding Controller doesn't work without new signature
		{
			valid:   true,
			name:    "Valid: Replace controller works with all signatures",
			signers: []string{BobKey1, AliceKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{BobDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid:   false,
			name:    "Not Valid: Replace controller doesn't work without old signatures",
			signers: []string{BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{BobDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID),
		},
		{
			valid:   false,
			name:    "Not Valid: Replace controller doesn't work without new signatures",
			signers: []string{AliceKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{BobDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID),
		},
		// add Controller
		{
			valid:   true,
			name:    "Valid: Adding second controller works",
			signers: []string{AliceKey1, CharlieKey3},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID, CharlieDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid:   false,
			name:    "Not Valid: Adding controller without old signature",
			signers: []string{BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID, BobDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID),
		},
		{
			valid:   false,
			name:    "Not Valid: Add controller without new signature doesn't work",
			signers: []string{AliceKey1},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID, BobDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID),
		},

		{
			valid:   true,
			name:    "Valid: Adding verification method with the same controller works",
			signers: []string{AliceKey1, AliceKey2},
			msg: &types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey2,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
					{
						Id:         AliceKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: AliceDID,
					},
				},
			},
		},
		{
			valid:   true,
			name:    "Valid: Keeping VM with controller different then subject untouched during update should not require Bob signature",
			signers: []string{CharlieKey1},
			msg: &types.MsgUpdateDidPayload{
				Id: CharlieDID,
				Authentication: []string{
					CharlieKey1,
					CharlieKey2,
					CharlieKey3,
				},

				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         CharlieKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: BobDID,
					},
					{
						Id:         CharlieKey2,
						Type:       Ed25519VerificationKey2020,
						Controller: BobDID,
					},
					{
						Id:         CharlieKey3,
						Type:       Ed25519VerificationKey2020,
						Controller: BobDID,
					},
					{
						Id:         CharlieKey4,
						Type:       Ed25519VerificationKey2020,
						Controller: CharlieDID,
					},
				},
			},
		},
		{
			valid:   true,
			name:    "Valid: Removing verification method is possible with any kind of valid Bob's key",
			signers: []string{BobKey1},
			msg: &types.MsgUpdateDidPayload{
				Id: BobDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         BobKey1,
						Type:       Ed25519VerificationKey2020,
						Controller: BobDID,
					},
				},
			},
			errMsg: fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", BobDID),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setup := InitEnv(t, keys)
			msg := tc.msg

			for _, vm := range msg.VerificationMethod {
				if vm.PublicKeyMultibase == "" {
					vm.PublicKeyMultibase, err = multibase.Encode(multibase.Base58BTC, keys[vm.Id].PublicKey)
				}
				require.NoError(t, err)
			}

			signerKeys := []SignerKey{}
			if tc.signerKeys != nil {
				signerKeys = tc.signerKeys
			} else {
				for _, signer := range tc.signers {
					signerKeys = append(signerKeys, SignerKey{
						signer: signer,
						key:    keys[signer].PrivateKey,
					})
				}
			}

			did, err := setup.SendUpdateDid(msg, signerKeys)
			rawExpectedDid := tc.msg.ToDid()
			expectedMsg := types.NormalizeDID(&rawExpectedDid)

			if tc.valid {
				require.Nil(t, err)
				require.Equal(t, expectedMsg.Id, did.Id)
				require.Equal(t, expectedMsg.Controller, did.Controller)
				require.Equal(t, expectedMsg.VerificationMethod, did.VerificationMethod)
				require.Equal(t, expectedMsg.Authentication, did.Authentication)
				require.Equal(t, expectedMsg.AssertionMethod, did.AssertionMethod)
				require.Equal(t, expectedMsg.CapabilityInvocation, did.CapabilityInvocation)
				require.Equal(t, expectedMsg.CapabilityDelegation, did.CapabilityDelegation)
				require.Equal(t, expectedMsg.KeyAgreement, did.KeyAgreement)
				require.Equal(t, expectedMsg.AlsoKnownAs, did.AlsoKnownAs)
				require.Equal(t, expectedMsg.Service, did.Service)
				require.Equal(t, expectedMsg.Context, did.Context)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
