package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgCreateDidPayloadValidation(t *testing.T) {
	cases := []struct {
		name     string
		struct_  *MsgCreateDidPayload
		isValid  bool
		errorMsg string
	}{
		{
			name: "test case 1",
			struct_: &MsgCreateDidPayload{
				Context:    nil,
				Id:         "bad did",
				Controller: nil,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                 "did1:cheqd:testnet:123456789abcdefg#sdfsdf",
						Type:               "jwk",
						Controller:         "",
						PublicKeyJwk:       nil,
						PublicKeyMultibase: "multibase",
					},
				},
				Authentication:       nil,
				AssertionMethod:      nil,
				CapabilityInvocation: nil,
				CapabilityDelegation: nil,
				KeyAgreement:         nil,
				AlsoKnownAs:          nil,
				Service:              nil,
			},
			isValid:  true,
			errorMsg: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.struct_.Validate()

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, err.Error(), tc.errorMsg)
			}
		})
	}
}
