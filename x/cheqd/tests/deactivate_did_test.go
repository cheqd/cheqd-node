package tests

import (
	"testing"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/require"
)

func TestDeactivateDid(t *testing.T) {
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
		msg        *types.MsgDeactivateDidPayload
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
			msg: &types.MsgDeactivateDidPayload{
				Id: AliceDID,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setup := InitEnv(t, keys)
			msg := tc.msg



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

			did, err := setup.SendDeactivateDid(msg, signerKeys)

			if tc.valid {
				require.Nil(t, err)
				require.Equal(t, tc.msg.Id, did.Id)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
