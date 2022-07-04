package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsId(t *testing.T) {
	cases := []struct {
		valid bool
		id    string
	}{
		{true, "123456789abcdefg"},
		{true, "123456789abcdefg123456789abcdefg"},
		{true, "3b9b8eec-5b5d-4382-86d8-9185126ff130"},
		{false, "sdf"},
		{false, "sdf:sdf"},
		{false, "12345"},
	}

	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			isDid := IsValidID(tc.id)

			if tc.valid {
				require.True(t, isDid)
			} else {
				require.False(t, isDid)
			}
		})
	}
}

func TestIsDid(t *testing.T) {
	cases := []struct {
		name      string
		valid     bool
		did       string
		method    string
		allowedNS []string
	}{
		{"Valid: Inputs: Method and namespace are set", true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet"}},
		{"Valid: Inputs: Method and namespaces are set", true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet", "mainnet"}},
		{"Valid: Inputs: Method not set", true, "did:cheqd:testnet:123456789abcdefg", "", []string{"testnet"}},
		{"Valid: Inputs: Method and namespaces are empty", true, "did:cheqd:testnet:123456789abcdefg", "", []string{}},
		{"Valid: Namespace is absent in DID", true, "did:cheqd:123456789abcdefg", "", []string{}},
		// Generic method validation
		{"Valid: Inputs: Method is not set and passed for NOTcheqd", true, "did:NOTcheqd:123456789abcdefg", "", []string{}},
		{"Valid: Inputs: Method and Namespaces are not set and passed for NOTcheqd", true, "did:NOTcheqd:123456789abcdefg123456789abcdefg", "", []string{}},

		{"Valid: Inputs: Order of namespaces changed", true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"mainnet", "testnet"}},
		// Wrong splitting (^did:([^:]+?)(:([^:]+?))?:([^:]+)$)
		{"Not valid: DID is not started from 'did'", false, "did123:cheqd:::123456789abcdefg", "cheqd", []string{"testnet"}},
		{"Not valid: empty namespace", false, "did:cheqd::123456789abcdefg", "cheqd", []string{"testnet"}},
		{"Not valid: a lot of ':'", false, "did:cheqd:::123456789abcdefg", "cheqd", []string{"testnet"}},
		{"Not valid: several DIDs in one string", false, "did:cheqd:testnet:123456789abcdefg:did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet"}},
		// Wrong method
		{"Not valid: method in DID is not the same as from Inputs", false, "did:NOTcheqd:testnet:123456789abcdefg", "cheqd", []string{"mainnet", "testnet"}},
		{"Not valid: method in Inputs is not the same as from DID", false, "did:cheqd:testnet:123456789abcdefg", "NOTcheqd", []string{"mainnet", "testnet"}},
		// Wrong namespace (^[a-zA-Z0-9]*)
		{"Not valid: / is not allowed for namespace", false, "did:cheqd:testnet/:123456789abcdefg", "cheqd", []string{}},
		{"Not valid: _ is not allowed for namespace", false, "did:cheqd:testnet_:123456789abcdefg", "cheqd", []string{}},
		{"Not valid: % is not allowed for namespace", false, "did:cheqd:testnet%:123456789abcdefg", "cheqd", []string{}},
		{"Not valid: * is not allowed for namespace", false, "did:cheqd:testnet*:123456789abcdefg", "cheqd", []string{}},
		{"Not valid: & is not allowed for namespace", false, "did:cheqd:testnet&:123456789abcdefg", "cheqd", []string{}},
		{"Not valid: @ is not allowed for namespace", false, "did:cheqd:testnet@:123456789abcdefg", "cheqd", []string{}},
		{"Not valid: namespace from Inputs is not the same as from DID", false, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"not_testnet"}},
		// Base58 checks (^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$)
		{"Not valid: O - is not allowed for UniqueID", false, "did:cheqd:testnet:123456789abcdefO", "cheqd", []string{}},
		{"Not valid: I - is not allowed for UniqueID", false, "did:cheqd:testnet:123456789abcdefI", "cheqd", []string{}},
		{"Not valid: l - is not allowed for UniqueID", false, "did:cheqd:testnet:123456789abcdefl", "cheqd", []string{}},
		{"Not valid: 0 - is not allowed for UniqueID", false, "did:cheqd:testnet:123456789abcdef0", "cheqd", []string{}},
		// Length checks (should be exactly 16 or 32)
		{"Not valid: UniqueID less then 16 symbols", false, "did:cheqd:testnet:123", "cheqd", []string{}},
		{"Not valid: UniqueID more then 16 symbols but less then 32", false, "did:cheqd:testnet:123456789abcdefgABCDEF", "cheqd", []string{}},
		{"Not valid: UniqueID more then 32 symbols", false, "did:cheqd:testnet:123456789abcdefg123456789abcdefgABCDEF", "cheqd", []string{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			isDid := IsValidDID(tc.did, tc.method, tc.allowedNS)

			if tc.valid {
				require.True(t, isDid)
			} else {
				require.False(t, isDid)
			}
		})
	}
}

func TestSplitJoin(t *testing.T) {
	cases := []string{
		"did:cheqd:testnet:123456789abcdefg",
		"did:cheqd:testnet:123456789abcdefg",
		"did:cheqd:testnet:123456789abcdefg",
		"did:cheqd:testnet:123456789abcdefg",
		"did:cheqd:123456789abcdefg",
		"did:NOTcheqd:123456789abcdefg",
		"did:NOTcheqd:123456789abcdefg123456789abcdefg",
		"did:cheqd:testnet:123456789abcdefg",
	}

	for _, tc := range cases {
		// Test split/join
		t.Run("split/join "+tc, func(t *testing.T) {
			method, namespace, id := MustSplitDID(tc)
			require.Equal(t, tc, JoinDID(method, namespace, id))
		})
	}
}

func TestMustSplitDID(t *testing.T) {
	require.Panicsf(t, func() {
		MustSplitDID("not did")
	}, "must panic")

	method, namespace, id := MustSplitDID("did:cheqd:mainnet:qqqqqqqqqqqqqqqq")
	require.Equal(t, "cheqd", method)
	require.Equal(t, "mainnet", namespace)
	require.Equal(t, "qqqqqqqqqqqqqqqq", id)
}
