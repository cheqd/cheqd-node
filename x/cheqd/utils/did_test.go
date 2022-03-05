package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsDid(t *testing.T) {
	cases := []struct {
		valid bool
		did   string
		method string
		allowedNS []string
	}{
		{true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet"}},
		{true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet", "mainnet"}},
		{true, "did:cheqd:testnet:123456789abcdefg", "", []string{"testnet"}},
		{true, "did:cheqd:testnet:123456789abcdefg", "", []string{}},
		{true, "did:cheqd:123456789abcdefg", "", []string{}},
		// Generic method validation
		{true, "did:NOTcheqd:123456789abcdefg", "", []string{}},
		{true, "did:NOTcheqd:123456789abcdefg123456789abcdefg", "", []string{}},

		{true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"mainnet", "testnet"}},
		// Wrong splitting (^did:([^:]+?)(:([^:]+?))?:([^:]+)$)
		{false, "did123:cheqd:::123456789abcdefg", "cheqd", []string{"testnet"}},
		{false, "did:cheqd::123456789abcdefg", "cheqd", []string{"testnet"}},
		{false, "did:cheqd:::123456789abcdefg", "cheqd", []string{"testnet"}},
		{false, "did:cheqd:testnet:123456789abcdefg:did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet"}},
		// Wrong method
		{false, "did:NOTcheqd:testnet:123456789abcdefg", "cheqd", []string{"mainnet", "testnet"}},
		{false, "did:cheqd:testnet:123456789abcdefg", "NOTcheqd", []string{"mainnet", "testnet"}},
		// Wrong namespace (^[a-zA-Z0-9]*)
		{false, "did:cheqd:testnet/:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet_:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet%:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet*:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet&:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet@/:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"not_testnet"}},
		// Base58 checks (^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$)
		{false, "did:cheqd:testnet:123456789abcdefO", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdefI", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdefl", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdef0", "cheqd", []string{}},
		// Length checks (should be exactly 16 or 32)
		{false, "did:cheqd:testnet:123", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdefgABCDEF", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdefg123456789abcdefgABCDEF", "cheqd", []string{}},

	}

	for _, tc := range cases {
		isDid := IsValidDID(tc.did, tc.method, tc.allowedNS)

		if tc.valid {
			require.True(t, isDid)
		} else {
			require.False(t, isDid)
		}
	}
}


//func TestDIDURLErrors(t *testing.T) {
//	cases := []struct {
//		didUrl   string
//		error_code int
//		error_string string
//	}{
//		{"did:cheqd:testnet:123456789abcdefg/", ErrStaticDIDURLPathAbemptyNotValid},
//	}
//
//	for _, tc := range cases {
//		err := ValidateDIDUrl(tc.didUrl, "", []string{})
//		require.NotNil(t, err)
//		require.Exactly(t, err.code, tc.error_code)
//
//}



