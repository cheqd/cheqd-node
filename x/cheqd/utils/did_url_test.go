package utils_test

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

var _ = Describe("DID-URL tests", func() {
	DescribeTable("Check the DID-URL join functionality (without functionality)",

	func(did_url string) {
		did, path, query, fragment := MustSplitDIDUrl(did_url)
		joined_did_url := JoinDIDUrl(did, path, query, fragment)
		Expect(joined_did_url).To(Equal(did_url))
	},
	Entry("All symbols", "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"),
	Entry("All symbols for path", "did:cheqd:testnet:123456789abcdefg/path/to/some/other/place/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff/"),
	Entry("All symbols for path and query", "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query"),
	Entry("All symbols for path and query and fragment", "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query#fragment"),
	Entry("All variants for path", "did:cheqd:testnet:123456789abcdefg/12/ab/AB/-/./_/~/!/$/&/'/(/)/*/+/,/;/=/:/@/%20/%ff"),
	Entry("Empty path", "did:cheqd:testnet:123456789abcdefg/"),
	Entry("All symbols for query", "did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"),
	Entry("All symbols for query and a lot of queries", "did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query=2?query=3/query=%A4"),
	Entry("All symbols for fragment", "did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"),
	Entry("All symbols for fragment and query", "did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"),
	Entry("Just fragment", "did:cheqd:testnet:123456789abcdefg#fragment"),
	Entry("Just query", "did:cheqd:testnet:123456789abcdefg?query"),
	)

	DescribeTable("Check the DID-URL Validation",

	func(expected bool, did_url string) {
		isValid := IsValidDIDUrl(did_url, "", []string{})

		Expect(isValid).To(Equal(expected))
	},
		// Path: all the possible symbols
		Entry("Valid: the whole alphabet for path", true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"),
		Entry("Valid: several paths", true, "did:cheqd:testnet:123456789abcdefg/path/to/some/other/place/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff/"),
		Entry("Valid: the whole alphabet with query", true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query"),
		Entry("Valid: the whole alphabet with query and fragment", true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query#fragment"),
		Entry("Valid: each possible symbols as a path", true, "did:cheqd:testnet:123456789abcdefg/12/ab/AB/-/./_/~/!/$/&/'/(/)/*/+/,/;/=/:/@/%20/%ff"),
		Entry("Valid: empty path", true, "did:cheqd:testnet:123456789abcdefg/"),
		// Query: all the possible variants
		Entry("Valid: the whole alphabet for query", true, "did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"),
		Entry("Valid: the whole alphabet for query and another query", true, "did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query=2?query=3/query=%A4"),
		// Fragment:
		Entry("Valid: the whole alphabet for fragment", true, "did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"),
		Entry("Valid: the whole alphabet with query and apth", true, "did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"),
		Entry("Valid: only fragment", true, "did:cheqd:testnet:123456789abcdefg#fragment"),
		Entry("Valid: only query", true, "did:cheqd:testnet:123456789abcdefg?query"),
		// Wrong cases
		Entry("Not valid: wrong HEXDIG for path (pct-encoded phrase)", false, "did:cheqd:testnet:123456789abcdefg/path%20%zz"),
		Entry("Not valid: wrong HEXDIG for query (pct-encoded phrase)", false, "did:cheqd:testnet:123456789abcdefg/path?query%20%zz"),
		Entry("Not valid: wrong HEXDIG for fragment (pct-encoded phrase)", false, "did:cheqd:testnet:123456789abcdefg/path?query#fragment%20%zz"),
		// Wrong splitting (^did:([^:]+?)(:([^:]+?))?:([^:]+)$)
		Entry("Not valid: starts with not 'did'", false, "did123:cheqd:::123456789abcdefg/path?query#fragment"),
		Entry("Not valid: empty namespace", false, "did:cheqd::123456789abcdefg/path?query#fragment"),
		Entry("Not valid: a lot of ':'", false, "did:cheqd:::123456789abcdefg/path?query#fragment"),
		Entry("Not valid: two DIDs in one", false, "did:cheqd:testnet:123456789abcdefg:did:cheqd:testnet:123456789abcdefg/path?query#fragment"),
		// Wrong namespace (^[a-zA-Z0-9]*)
		Entry("Not valid:  / - is not allowed for namespace", false, "did:cheqd:testnet/:123456789abcdefg/path?query#fragment"),
		Entry("Not valid: _  - is not allowed for namespace", false, "did:cheqd:testnet_:123456789abcdefg/path?query#fragment"),
		Entry("Not valid: % - is not allowed for namespace", false, "did:cheqd:testnet%:123456789abcdefg/path?query#fragment"),
		Entry("Not valid: * - is not allowed for namespace", false, "did:cheqd:testnet*:123456789abcdefg/path?query#fragment"),
		Entry("Not valid: & - is not allowed for namespace", false, "did:cheqd:testnet&:123456789abcdefg/path?query#fragment"),
		Entry("Not valid: @ - is not allowed for namespace", false, "did:cheqd:testnet@/:123456789abcdefg/path?query#fragment"),
		// Base58 checks (^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$)
		Entry("Not valid: O - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdefO/path?query#fragment"),
		Entry("Not valid: I - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdefI/path?query#fragment"),
		Entry("Not valid: l - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdefl/path?query#fragment"),
		Entry("Not valid: 0 - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdef0/path?query#fragment"),
		// Length checks (should be exactly 16 or 32)
		Entry("Not valid: UniqueID less then 16 symbols", false, "did:cheqd:testnet:123/path?query#fragment"),
		Entry("Not valid: UniqueID more then 16 symbols but less then 32", false, "did:cheqd:testnet:123456789abcdefgABCDEF/path?query#fragment"),
		Entry("Not valid: UniqueID more then 32 symbols", false, "did:cheqd:testnet:123456789abcdefg123456789abcdefgABCDEF/path?query#fragment"),
		Entry("Not valid: Split should return error", false, "qwerty"),
	)
})

// func TestIsDidURL(t *testing.T) {
// 	cases := []struct {
// 		name   string
// 		valid  bool
// 		didUrl string
// 	}{
// 		// Path: all the possible symbols
// 		{"Valid: the whole alphabet for path", true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"},
// 		{"Valid: several paths", true, "did:cheqd:testnet:123456789abcdefg/path/to/some/other/place/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff/"},
// 		{"Valid: the whole alphabet with query", true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query"},
// 		{"Valid: the whole alphabet with query and fragment", true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query#fragment"},
// 		{"Valid: each possible symbols as a path", true, "did:cheqd:testnet:123456789abcdefg/12/ab/AB/-/./_/~/!/$/&/'/(/)/*/+/,/;/=/:/@/%20/%ff"},
// 		{"Valid: empty path", true, "did:cheqd:testnet:123456789abcdefg/"},
// 		// Query: all the possible variants
// 		{"Valid: the whole alphabet for query", true, "did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"},
// 		{"Valid: the whole alphabet for query and another query", true, "did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query=2?query=3/query=%A4"},
// 		// Fragment:
// 		{"Valid: the whole alphabet for fragment", true, "did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"},
// 		{"Valid: the whole alphabet with query and apth", true, "did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"},
// 		{"Valid: only fragment", true, "did:cheqd:testnet:123456789abcdefg#fragment"},
// 		{"Valid: only query", true, "did:cheqd:testnet:123456789abcdefg?query"},
// 		// Wrong cases
// 		{"Not valid: wrong HEXDIG for path (pct-encoded phrase)", false, "did:cheqd:testnet:123456789abcdefg/path%20%zz"},
// 		{"Not valid: wrong HEXDIG for query (pct-encoded phrase)", false, "did:cheqd:testnet:123456789abcdefg/path?query%20%zz"},
// 		{"Not valid: wrong HEXDIG for fragment (pct-encoded phrase)", false, "did:cheqd:testnet:123456789abcdefg/path?query#fragment%20%zz"},
// 		// Wrong splitting (^did:([^:]+?)(:([^:]+?))?:([^:]+)$)
// 		{"Not valid: starts with not 'did'", false, "did123:cheqd:::123456789abcdefg/path?query#fragment"},
// 		{"Not valid: empty namespace", false, "did:cheqd::123456789abcdefg/path?query#fragment"},
// 		{"Not valid: a lot of ':'", false, "did:cheqd:::123456789abcdefg/path?query#fragment"},
// 		{"Not valid: two DIDs in one", false, "did:cheqd:testnet:123456789abcdefg:did:cheqd:testnet:123456789abcdefg/path?query#fragment"},
// 		// Wrong namespace (^[a-zA-Z0-9]*)
// 		{"Not valid:  / - is not allowed for namespace", false, "did:cheqd:testnet/:123456789abcdefg/path?query#fragment"},
// 		{"Not valid: _  - is not allowed for namespace", false, "did:cheqd:testnet_:123456789abcdefg/path?query#fragment"},
// 		{"Not valid: % - is not allowed for namespace", false, "did:cheqd:testnet%:123456789abcdefg/path?query#fragment"},
// 		{"Not valid: * - is not allowed for namespace", false, "did:cheqd:testnet*:123456789abcdefg/path?query#fragment"},
// 		{"Not valid: & - is not allowed for namespace", false, "did:cheqd:testnet&:123456789abcdefg/path?query#fragment"},
// 		{"Not valid: @ - is not allowed for namespace", false, "did:cheqd:testnet@/:123456789abcdefg/path?query#fragment"},
// 		// Base58 checks (^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$)
// 		{"Not valid: O - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdefO/path?query#fragment"},
// 		{"Not valid: I - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdefI/path?query#fragment"},
// 		{"Not valid: l - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdefl/path?query#fragment"},
// 		{"Not valid: 0 - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdef0/path?query#fragment"},
// 		// Length checks (should be exactly 16 or 32)
// 		{"Not valid: UniqueID less then 16 symbols", false, "did:cheqd:testnet:123/path?query#fragment"},
// 		{"Not valid: UniqueID more then 16 symbols but less then 32", false, "did:cheqd:testnet:123456789abcdefgABCDEF/path?query#fragment"},
// 		{"Not valid: UniqueID more then 32 symbols", false, "did:cheqd:testnet:123456789abcdefg123456789abcdefgABCDEF/path?query#fragment"},
// 		{"Not valid: Split should return error", false, "qwerty"},
// 	}

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			isDid := IsValidDIDUrl(tc.didUrl, "", []string{})

// 			if tc.valid {
// 				require.True(t, isDid)
// 			} else {
// 				require.False(t, isDid)
// 			}
// 		})
// 	}
// }

// func TestDidURLJoin(t *testing.T) Entry("", 
// 	cases := []stringEntry("", 
// 		"did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff",
// 		"did:cheqd:testnet:123456789abcdefg/path/to/some/other/place/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff/",
// 		"did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query",
// 		"did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query#fragment",
// 		"did:cheqd:testnet:123456789abcdefg/12/ab/AB/-/./_/~/!/$/&/'/(/)/*/+/,/;/=/:/@/%20/%ff",
// 		"did:cheqd:testnet:123456789abcdefg/",
// 		"did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?",
// 		"did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query=2?query=3/query=%A4",
// 		"did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff",
// 		"did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?",
// 		"did:cheqd:testnet:123456789abcdefg#fragment",
// 		"did:cheqd:testnet:123456789abcdefg?query",
// 	}

// 	for _, tc := range cases Entry("", 
// 		t.Run("split/join"+tc, func(t *testing.T) {
// 			did, path, query, fragment := MustSplitDIDUrl(tc)
// 			require.Equal(t, tc, JoinDIDUrl(did, path, query, fragment))
// 		})
// 	}
// }
