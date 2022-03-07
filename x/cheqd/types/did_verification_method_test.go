package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerificationMethod(t *testing.T) {
	cases := []struct {
		name     string
		struct_  VerificationMethod
		isValid  bool
		errorMsg string
	}{
		{
			name: "valid method with multibase key",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa",
				Type:               "Ed25519VerificationKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       nil,
				PublicKeyMultibase: "multibase",
			},
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "valid method with jwk key",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa",
				Type:               "JsonWebKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       []*KeyValuePair{
					{
						Key:   "key",
						Value: "value",
					},
				},
				PublicKeyMultibase: "",
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
