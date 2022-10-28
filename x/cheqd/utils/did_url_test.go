package utils_test

import (
	. "github.com/cheqd/cheqd-node/x/cheqd/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DID-URL tests", func() {
	type TestCaseUUIDIdStruct struct {
		inputDIDUrl    string
		expectedDIDUrl string
	}

	DescribeTable("Check the DID-URL join functionality (without functionality)",

		func(did_url string) {
			did, path, query, fragment := MustSplitDIDUrl(did_url)
			joined_did_url := JoinDIDUrl(did, path, query, fragment)
			Expect(joined_did_url).To(Equal(did_url))
		},
		Entry("All symbols", "did:cheqd:testnet:zABCDEFG123456789abcd/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"),
		Entry("All symbols for path", "did:cheqd:testnet:zABCDEFG123456789abcd/path/to/some/other/place/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff/"),
		Entry("All symbols for path and query", "did:cheqd:testnet:zABCDEFG123456789abcd/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query"),
		Entry("All symbols for path and query and fragment", "did:cheqd:testnet:zABCDEFG123456789abcd/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query#fragment"),
		Entry("All variants for path", "did:cheqd:testnet:zABCDEFG123456789abcd/12/ab/AB/-/./_/~/!/$/&/'/(/)/*/+/,/;/=/:/@/%20/%ff"),
		Entry("Empty path", "did:cheqd:testnet:zABCDEFG123456789abcd/"),
		Entry("All symbols for query", "did:cheqd:testnet:zABCDEFG123456789abcd/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"),
		Entry("All symbols for query and a lot of queries", "did:cheqd:testnet:zABCDEFG123456789abcd/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query=2?query=3/query=%A4"),
		Entry("All symbols for fragment", "did:cheqd:testnet:zABCDEFG123456789abcd/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"),
		Entry("All symbols for fragment and query", "did:cheqd:testnet:zABCDEFG123456789abcd/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"),
		Entry("Just fragment", "did:cheqd:testnet:zABCDEFG123456789abcd#fragment"),
		Entry("Just query", "did:cheqd:testnet:zABCDEFG123456789abcd?query"),
	)

	DescribeTable("Check the DID-URL Validation",

		func(expected bool, did_url string) {
			isValid := IsValidDIDUrl(did_url, "", []string{})

			Expect(isValid).To(Equal(expected))
		},
		// Path: all the possible symbols
		Entry("Valid: the whole alphabet for path", true, "did:cheqd:testnet:zABCDEFG123456789abcd/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"),
		Entry("Valid: several paths", true, "did:cheqd:testnet:zABCDEFG123456789abcd/path/to/some/other/place/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff/"),
		Entry("Valid: the whole alphabet with query", true, "did:cheqd:testnet:zABCDEFG123456789abcd/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query"),
		Entry("Valid: the whole alphabet with query and fragment", true, "did:cheqd:testnet:zABCDEFG123456789abcd/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query#fragment"),
		Entry("Valid: each possible symbols as a path", true, "did:cheqd:testnet:zABCDEFG123456789abcd/12/ab/AB/-/./_/~/!/$/&/'/(/)/*/+/,/;/=/:/@/%20/%ff"),
		Entry("Valid: empty path", true, "did:cheqd:testnet:zABCDEFG123456789abcd/"),
		// Query: all the possible variants
		Entry("Valid: the whole alphabet for query", true, "did:cheqd:testnet:zABCDEFG123456789abcd/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"),
		Entry("Valid: the whole alphabet for query and another query", true, "did:cheqd:testnet:zABCDEFG123456789abcd/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query=2?query=3/query=%A4"),
		// Fragment:
		Entry("Valid: the whole alphabet for fragment", true, "did:cheqd:testnet:zABCDEFG123456789abcd/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"),
		Entry("Valid: the whole alphabet with query and apth", true, "did:cheqd:testnet:zABCDEFG123456789abcd/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"),
		Entry("Valid: only fragment", true, "did:cheqd:testnet:zABCDEFG123456789abcd#fragment"),
		Entry("Valid: only query", true, "did:cheqd:testnet:zABCDEFG123456789abcd?query"),
		// Wrong cases
		Entry("Not valid: wrong HEXDIG for path (pct-encoded phrase)", false, "did:cheqd:testnet:zABCDEFG123456789abcd/path%20%zz"),
		Entry("Not valid: wrong HEXDIG for query (pct-encoded phrase)", false, "did:cheqd:testnet:zABCDEFG123456789abcd/path?query%20%zz"),
		Entry("Not valid: wrong HEXDIG for fragment (pct-encoded phrase)", false, "did:cheqd:testnet:zABCDEFG123456789abcd/path?query#fragment%20%zz"),
		// Wrong splitting (^did:([^:]+?)(:([^:]+?))?:([^:]+)$)
		Entry("Not valid: starts with not 'did'", false, "did123:cheqd:::zABCDEFG123456789abcd/path?query#fragment"),
		Entry("Not valid: empty namespace", false, "did:cheqd::zABCDEFG123456789abcd/path?query#fragment"),
		Entry("Not valid: a lot of ':'", false, "did:cheqd:::zABCDEFG123456789abcd/path?query#fragment"),
		Entry("Not valid: two DIDs in one", false, "did:cheqd:testnet:zABCDEFG123456789abcd:did:cheqd:testnet:zABCDEFG123456789abcd/path?query#fragment"),
		// Wrong namespace (^[a-zA-Z0-9]*)
		Entry("Not valid:  / - is not allowed for namespace", false, "did:cheqd:testnet/:zABCDEFG123456789abcd/path?query#fragment"),
		Entry("Not valid: _  - is not allowed for namespace", false, "did:cheqd:testnet_:zABCDEFG123456789abcd/path?query#fragment"),
		Entry("Not valid: % - is not allowed for namespace", false, "did:cheqd:testnet%:zABCDEFG123456789abcd/path?query#fragment"),
		Entry("Not valid: * - is not allowed for namespace", false, "did:cheqd:testnet*:zABCDEFG123456789abcd/path?query#fragment"),
		Entry("Not valid: & - is not allowed for namespace", false, "did:cheqd:testnet&:zABCDEFG123456789abcd/path?query#fragment"),
		Entry("Not valid: @ - is not allowed for namespace", false, "did:cheqd:testnet@/:zABCDEFG123456789abcd/path?query#fragment"),
		// Base58 checks (^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$)
		Entry("Not valid: O - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdefO/path?query#fragment"),
		Entry("Not valid: I - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdefI/path?query#fragment"),
		Entry("Not valid: l - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdefl/path?query#fragment"),
		Entry("Not valid: 0 - is not allowed for base58", false, "did:cheqd:testnet:123456789abcdef0/path?query#fragment"),
		// Length checks (should be exactly 16 or 32)
		Entry("Not valid: UniqueID less then 16 bytes", false, "did:cheqd:testnet:123/path?query#fragment"),
		Entry("Not valid: UniqueID more then 16 bytes", false, "did:cheqd:testnet:zABCDEFG123456789abcdzABCDEFG123456789abcdABCDEF/path?query#fragment"),
		Entry("Not valid: Split should return error", false, "qwerty"),
	)

	DescribeTable("UUID test cases", func(testCase TestCaseUUIDIdStruct) {
		result := NormalizeDIDUrl(testCase.inputDIDUrl)
		Expect(result).To(Equal(testCase.expectedDIDUrl))
	},

		Entry(
			"base58 identifier - not changed",
			TestCaseUUIDIdStruct{
				inputDIDUrl:    "did:cheqd:testnet:zABCDEFG123456789abcd#key1",
				expectedDIDUrl: "did:cheqd:testnet:zABCDEFG123456789abcd#key1",
			}),

		Entry(
			"Mixed case UUID",
			TestCaseUUIDIdStruct{
				inputDIDUrl:    "did:cheqd:testnet:BAbbba14-f294-458a-9b9c-474d188680fd#key1",
				expectedDIDUrl: "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd#key1",
			}),

		Entry(
			"Low case UUID",
			TestCaseUUIDIdStruct{
				inputDIDUrl:    "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd#key1",
				expectedDIDUrl: "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd#key1",
			}),

		Entry(
			"Upper case UUID",
			TestCaseUUIDIdStruct{
				inputDIDUrl:    "did:cheqd:testnet:A86F9CAE-0902-4a7c-a144-96b60ced2FC9#key1",
				expectedDIDUrl: "did:cheqd:testnet:a86f9cae-0902-4a7c-a144-96b60ced2fc9#key1",
			}),
	)
})
