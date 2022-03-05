package utils

import (
"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestController(t *testing.T) {
	cases := []struct {
		name string
		isValid bool
		did_list []string
		errorMsg string
	}{
		{"Valid: One element", true, []string{"did:cheqd:testnet:123456789abcdefg"}, ""},
		{"Valid: More then one", true, []string{"did:cheqd:testnet:123456789abcdefg", "did:cheqd:testnet:gfedcba987654321"}, ""},
		{"Not valid: Wrong DID, error is passing", false, []string{"did:cheqd:testnet:badDid"}, "unique id length should be 16 or 32 symbols"},
		{"Not valid: More then 1 error", false, []string{"did:cheqd:testnet:first", "did1:cheqd:testnet:first"}, "unique id length should be 16 or 32 symbols, did must match the following regex exactly"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateDIDList(tc.did_list, "", []string{})

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), tc.errorMsg))
			}
		})
	}
}