package utils

import (
	"github.com/multiformats/go-multibase"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateMultibase(t *testing.T) {
	cases := []struct {
		name string
		data  string
		encoding multibase.Encoding
		valid bool
	}{
		{"Valid: General pmbkey", "zABCDEFG123456789", multibase.Base58BTC, true},
		{"Not Valid: cannot be empty", "", multibase.Base58BTC, false},
		{"Not Valid: without z but base58", "ABCDEFG123456789", multibase.Base58BTC, false},
		{"Not Valid: without z and not base58", "OIl0ABCDEFG123456789", multibase.Base58BTC, false},
		{"Not Valid: with z but not base58", "zOIl0ABCDEFG123456789", multibase.Base58BTC, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMultibase(tc.data)

			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestValidateBase58(t *testing.T) {
	cases := []struct {
		name string
		data  string
		valid bool
	}{
		{"Valid: General pmbkey", "ABCDEFG123456789", true},
		{"Not Valid: cannot be empty", "", false},
		{"Not Valid: not base58", "OIl0ABCDEFG123456789", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateBase58(tc.data)

			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
