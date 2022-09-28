package tests

import (
	"crypto/ed25519"
	"fmt"
	"testing"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/require"
)

func TestDeactivateDid(t *testing.T) {
	keys := GenerateTestKeys()

	cases := []struct {
		valid       bool
		name        string
		signerKeys  []SignerKey
		msg         *types.MsgDeactivateDidPayload
		deactivared bool
		errMsg      string
	}{
		{
			valid: true,
			name:  "Valid: Deactivate DID",
			signerKeys: []SignerKey{
				{
					signer: AliceKey1,
					key:    keys[AliceKey1].PrivateKey,
				},
			},
			msg: &types.MsgDeactivateDidPayload{
				Id: AliceDID,
			},
		},
		{
			valid: false,
			name:  "Not Valid: Not found",
			signerKeys: []SignerKey{
				{
					signer: AliceKey1,
					key:    keys[AliceKey1].PrivateKey,
				},
			},
			msg: &types.MsgDeactivateDidPayload{
				Id: NotFounDID,
			},
			errMsg: NotFounDID + ": DID Doc not found",
		},
		{
			valid: false,
			name:  "Not Valid: Already deactivated",
			signerKeys: []SignerKey{
				{
					signer: DeactivatedDIDKey,
					key:    keys[DeactivatedDIDKey].PrivateKey,
				},
			},
			msg: &types.MsgDeactivateDidPayload{
				Id: DeactivatedDID,
			},
			deactivared: true,
			errMsg:      DeactivatedDID + ": DID Doc already deactivated",
		},
		{
			valid: false,
			name:  "Not Valid: Invalid signature",
			signerKeys: []SignerKey{
				{
					signer: BobKey1,
					key:    keys[BobKey1].PrivateKey,
				},
			},
			msg: &types.MsgDeactivateDidPayload{
				Id: AliceDID,
			},
			deactivared: false,
			errMsg:      fmt.Sprintf("signer: %s: signature is required but not found", AliceDID),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setup := InitEnv(t, keys)
			msg := tc.msg

			signerKeys := map[string]ed25519.PrivateKey{}
			for _, signature := range tc.signerKeys {
				signerKeys[signature.signer] = signature.key
			}

			did, err := setup.SendDeactivateDid(msg, signerKeys)

			if tc.valid {
				require.Nil(t, err)
				require.True(t, did.Deactivated)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
				if did != nil {
					require.Equal(t, tc.deactivared, did.Deactivated)
				}
			}
		})
	}
}
