package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgUpdateDidValidation(t *testing.T) {
	cases := []struct {
		name     string
		struct_  *MsgUpdateDid
		isValid  bool
		errorMsg string
	}{
		{
			name: "positive",
			struct_: &MsgUpdateDid{
				Payload: &MsgUpdateDidPayload{
					Id:         "did:cheqd:testnet:123456789abcdefg",
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 "did:cheqd:testnet:123456789abcdefg#key1",
							Type:               "Ed25519VerificationKey2020",
							Controller:         "did:cheqd:testnet:123456789abcdefg",
							PublicKeyMultibase: "multibase",
						},
					},
					Authentication:       []string{"did:cheqd:testnet:123456789abcdefg#key1", "did:cheqd:testnet:123456789abcdefg#aaa"},
					VersionId: "version1",
				},
				Signatures: nil,
			},
			isValid:  true,
		},
		{
			name: "negative: relationship duplicates",
			struct_: &MsgUpdateDid{
				Payload: &MsgUpdateDidPayload{
					Id:         "did:cheqd:testnet:123456789abcdefg",
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 "did:cheqd:testnet:123456789abcdefg#key1",
							Type:               "Ed25519VerificationKey2020",
							Controller:         "did:cheqd:testnet:123456789abcdefg",
							PublicKeyMultibase: "multibase",
						},
					},
					Authentication:       []string{"did:cheqd:testnet:123456789abcdefg#key1", "did:cheqd:testnet:123456789abcdefg#key1"},
					VersionId: "version1",
				},
				Signatures: nil,
			},
			isValid:  false,
			errorMsg: "payload: (authentication: there should be no duplicates.).: basic validation failed",
		},
		{
			name: "negative: version id is required",
			struct_: &MsgUpdateDid{
				Payload: &MsgUpdateDidPayload{
					Id:         "did:cheqd:testnet:123456789abcdefg",
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 "did:cheqd:testnet:123456789abcdefg#key1",
							Type:               "Ed25519VerificationKey2020",
							Controller:         "did:cheqd:testnet:123456789abcdefg",
							PublicKeyMultibase: "multibase",
						},
					},
					Authentication:       []string{"did:cheqd:testnet:123456789abcdefg#key1", "did:cheqd:testnet:123456789abcdefg#aaa"},
				},
				Signatures: nil,
			},
			isValid:  false,
			errorMsg: "payload: (version_id: cannot be blank.).: basic validation failed",
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
