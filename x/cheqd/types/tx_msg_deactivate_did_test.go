package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgDeactivateDidValidation(t *testing.T) {
	cases := []struct {
		name     string
		struct_  *MsgDeactivateDid
		isValid  bool
		errorMsg string
	}{
		{
			name: "positive",
			struct_: &MsgDeactivateDid{
				Payload: &MsgDeactivateDidPayload{
					Id: "did:cheqd:testnet:123456789abcdefg",
				},
				Signatures: nil,
			},
			isValid: true,
		},
		{
			name: "negative: invalid did method",
			struct_: &MsgDeactivateDid{
				Payload: &MsgDeactivateDidPayload{
					Id: "did:cheqdttt:testnet:123456789abcdefg",
				},
				Signatures: nil,
			},
			isValid:  false,
			errorMsg: "payload: (id: did method must be: cheqd.).: basic validation failed",
		},
		{
			name: "negative: id is required",
			struct_: &MsgDeactivateDid{
				Payload:    &MsgDeactivateDidPayload{},
				Signatures: nil,
			},
			isValid:  false,
			errorMsg: "payload: (id: cannot be blank.).: basic validation failed",
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
