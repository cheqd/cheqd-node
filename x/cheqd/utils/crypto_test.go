package utils

import (
	"github.com/multiformats/go-multibase"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateEd25519PubKey(t *testing.T) {
	cases := []struct {
		name string
		key  string
		valid bool
		errorMsg string
	}{
		{"Valid: General Ed25519 public key", "zF1hVGXXK9rmx5HhMTpGnGQJiab9qrFJbQXBRhSmYjQWX", true, ""},
		{"Valid: General Ed25519 public key", "zF1hVGXXK9rmx5HhMTpGnGQJiab9qr1111111111111", false, "ed25519: bad public key length: 31"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, keyBytes, _ := multibase.Decode(tc.key)
			err := ValidateEd25519PubKey(keyBytes)

			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			}
		})
	}
}