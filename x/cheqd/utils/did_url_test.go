package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsDidURL(t *testing.T) {
	cases := []struct {
		valid bool
		didUrl   string
	}{
		// Path: all the possible symbols
		{true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"},
		{true, "did:cheqd:testnet:123456789abcdefg/path/to/some/other/place/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff/"},
		{true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query"},
		{true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query#fragment"},
		{true, "did:cheqd:testnet:123456789abcdefg/12/ab/AB/-/./_/~/!/$/&/'/(/)/*/+/,/;/=/:/@/%20/%ff"},
		{true, "did:cheqd:testnet:123456789abcdefg/"},
		// Query: all the possible variants
		{true, "did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"},
		{true, "did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query=2?query=3/query=%A4"},
		// Fragment:
		{true, "did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"},
		{true, "did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"},
		// Wrong cases
		{false, "did:cheqd:testnet:123456789abcdefg/path%20%zz"},
		{false, "did:cheqd:testnet:123456789abcdefg/path?query%20%zz"},
		{false, "did:cheqd:testnet:123456789abcdefg/path?query#fragment%20%zz"},
		// Wrong splitting (^did:([^:]+?)(:([^:]+?))?:([^:]+)$)
		{false, "did123:cheqd:::123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd::123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:::123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdefg:did:cheqd:testnet:123456789abcdefg/path?query#fragment"},
		// Wrong namespace (^[a-zA-Z0-9]*)
		{false, "did:cheqd:testnet/:123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet_:123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet%:123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet*:123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet&:123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet@/:123456789abcdefg/path?query#fragment"},
		// Base58 checks (^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$)
		{false, "did:cheqd:testnet:123456789abcdefO/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdefI/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdefl/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdef0/path?query#fragment"},
		// Length checks (should be exactly 16 or 32)
		{false, "did:cheqd:testnet:123/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdefgABCDEF/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdefg123456789abcdefgABCDEF/path?query#fragment"},
	}

	for _, tc := range cases {
		isDid := IsValidDIDUrl(tc.didUrl, "", []string{})

		if tc.valid {
			require.True(t, isDid)
		} else {
			require.False(t, isDid)
		}
	}
}