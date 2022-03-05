package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerificationMethod(t *testing.T) {
	cases := []struct {
		name   string
		struct_ VerificationMethod
		isValid bool
		errorMsg string
	}{
		{
			name: "test case 1",
			struct_: VerificationMethod{
				Id:                 "did1:cheqd:testnet:123456789abcdefg#sdfsdf",
				Type:               "jwk",
				Controller:         "",
				PublicKeyJwk:       nil,
				PublicKeyMultibase: "multibase",
			},
			isValid: true,
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
