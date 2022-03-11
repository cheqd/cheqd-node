package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgCreateDidValidation(t *testing.T) {
	cases := []struct {
		name     string
		struct_  *MsgCreateDid
		isValid  bool
		errorMsg string
	}{
		{
			name: "positive",
			struct_: &MsgCreateDid{
				Payload: &MsgCreateDidPayload{
					Id:         "did:cheqd:testnet:123456789abcdefg",
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 "did:cheqd:testnet:123456789abcdefg#key1",
							Type:               "Ed25519VerificationKey2020",
							Controller:         "did:cheqd:testnet:123456789abcdefg",
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
					Authentication:       []string{"did:cheqd:testnet:123456789abcdefg#key1", "did:cheqd:testnet:123456789abcdefg#aaa"},
				},
				Signatures: nil,
			},
			isValid:  true,
		},
		{
			name: "negative: relationship duplicates",
			struct_: &MsgCreateDid{
				Payload: &MsgCreateDidPayload{
					Id:         "did:cheqd:testnet:123456789abcdefg",
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 "did:cheqd:testnet:123456789abcdefg#key1",
							Type:               "Ed25519VerificationKey2020",
							Controller:         "did:cheqd:testnet:123456789abcdefg",
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
					Authentication:       []string{"did:cheqd:testnet:123456789abcdefg#key1", "did:cheqd:testnet:123456789abcdefg#key1"},
				},
				Signatures: nil,
			},
			isValid:  false,
			errorMsg: "payload: (authentication: there should be no duplicates.).: basic validation failed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.struct_.ValidateBasic()

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, err.Error(), tc.errorMsg)
			}
		})
	}
}
