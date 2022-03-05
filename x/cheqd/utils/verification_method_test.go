package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPublicKeyMultibase(t *testing.T) {
	cases := []struct {
		name string
		valid bool
		key   string
	}{
		{"Valid: General pmbkey", true, "zABCDEFG123456789"},
		{"Not Valid: cannot be empty", false, ""},
		{"Not Valid: without z but base58", false, "ABCDEFG123456789"},
		{"Not Valid: without z and not base58", false, "OIl0ABCDEFG123456789"},
		{"Not Valid: with z but not base58", false, "zOIl0ABCDEFG123456789"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePublicKeyMultibase(tc.key)

			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

