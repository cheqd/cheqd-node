package utils_test

import (
	. "github.com/cheqd/cheqd-node/x/cheqd/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DID checks", func() {
	DescribeTable("Check is valid ID (for example did:cheqD:testnet:123456789abcdefg and ID is 123456789abcdefg)",

		func(expected bool, did string) {
			Expect(IsValidID(did)).To(Equal(expected))
		},
		Entry("Valid base58, 16 symbols", true, "123456789abcdefg"),
		Entry("Valid indy-style identifier", true, "123456789abcdefgre"),
		Entry("Valid base58, 32 symbols", true, "123456789abcdefg123456789abcdefg"),
		Entry("Valid UUID", true, "3b9b8eec-5b5d-4382-86d8-9185126ff130"),
		Entry("Not valid, not base58 symbols", false, "12345678abcdIlO0"),
		Entry("Not valid, length", false, "sdf"),
		Entry("Not valid, length and format", false, "sdf:sdf"),
		Entry("Not valid, length", false, "12345"),
	)

	DescribeTable("DID validation",

		func(expected bool, did string, method string, allowedNamespaces []string) {
			Expect(IsValidDID(did, method, allowedNamespaces)).To(Equal(expected))
		},

		Entry("Valid: Inputs: Method and namespace are set", true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet"}),
		Entry("Valid: Inputs: Method and namespaces are set", true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet", "mainnet"}),
		Entry("Valid: Inputs: Method not set", true, "did:cheqd:testnet:123456789abcdefg", "", []string{"testnet"}),
		Entry("Valid: Inputs: Method and namespaces are empty", true, "did:cheqd:testnet:123456789abcdefg", "", []string{}),
		Entry("Valid: Namespace is absent in DID", true, "did:cheqd:123456789abcdefg", "", []string{}),
		// Generic method validation
		Entry("Valid: Inputs: Method is not set and passed for NOTcheqd", true, "did:NOTcheqd:123456789abcdefg", "", []string{}),
		Entry("Valid: Inputs: Method and Namespaces are not set and passed for NOTcheqd", true, "did:NOTcheqd:123456789abcdefg123456789abcdefg", "", []string{}),

		Entry("Valid: Inputs: Order of namespaces changed", true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"mainnet", "testnet"}),
		Entry("Valid: UniqueID more then 16 symbols and less then 32", true, "did:cheqd:testnet:123456789abcdefgABCDEF", "cheqd", []string{"mainnet", "testnet"}),
		// Wrong splitting (^did:([^:]+?)(:([^:]+?))?:([^:]+)$)
		Entry("Not valid: DID is not started from 'did'", false, "did123:cheqd:::123456789abcdefg", "cheqd", []string{"testnet"}),
		Entry("Not valid: empty namespace", false, "did:cheqd::123456789abcdefg", "cheqd", []string{"testnet"}),
		Entry("Not valid: a lot of ':'", false, "did:cheqd:::123456789abcdefg", "cheqd", []string{"testnet"}),
		Entry("Not valid: several DIDs in one string", false, "did:cheqd:testnet:123456789abcdefg:did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet"}),
		// Wrong method
		Entry("Not valid: method in DID is not the same as from Inputs", false, "did:NOTcheqd:testnet:123456789abcdefg", "cheqd", []string{"mainnet", "testnet"}),
		Entry("Not valid: method in Inputs is not the same as from DID", false, "did:cheqd:testnet:123456789abcdefg", "NOTcheqd", []string{"mainnet", "testnet"}),
		// Wrong namespace (^[a-zA-Z0-9]*)
		Entry("Not valid: / is not allowed for namespace", false, "did:cheqd:testnet/:123456789abcdefg", "cheqd", []string{}),
		Entry("Not valid: _ is not allowed for namespace", false, "did:cheqd:testnet_:123456789abcdefg", "cheqd", []string{}),
		Entry("Not valid: % is not allowed for namespace", false, "did:cheqd:testnet%:123456789abcdefg", "cheqd", []string{}),
		Entry("Not valid: * is not allowed for namespace", false, "did:cheqd:testnet*:123456789abcdefg", "cheqd", []string{}),
		Entry("Not valid: & is not allowed for namespace", false, "did:cheqd:testnet&:123456789abcdefg", "cheqd", []string{}),
		Entry("Not valid: @ is not allowed for namespace", false, "did:cheqd:testnet@:123456789abcdefg", "cheqd", []string{}),
		Entry("Not valid: namespace from Inputs is not the same as from DID", false, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"not_testnet"}),
		// Base58 checks (^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$)
		Entry("Not valid: O - is not allowed for UniqueID", false, "did:cheqd:testnet:123456789abcdefO", "cheqd", []string{}),
		Entry("Not valid: I - is not allowed for UniqueID", false, "did:cheqd:testnet:123456789abcdefI", "cheqd", []string{}),
		Entry("Not valid: l - is not allowed for UniqueID", false, "did:cheqd:testnet:123456789abcdefl", "cheqd", []string{}),
		Entry("Not valid: 0 - is not allowed for UniqueID", false, "did:cheqd:testnet:123456789abcdef0", "cheqd", []string{}),
		// Length checks (should be exactly 16 or 32)
		Entry("Not valid: UniqueID less then 16 symbols", false, "did:cheqd:testnet:123", "cheqd", []string{}),
		Entry("Not valid: UniqueID more then 16 symbols but less then 32", false, "did:cheqd:testnet:123456789abcdefgABCDEF", "cheqd", []string{}),
		Entry("Not valid: UniqueID more then 32 symbols", false, "did:cheqd:testnet:123456789abcdefg123456789abcdefgABCDEF", "cheqd", []string{}),
	)

	// The next test will check more functionality aspects. This one is for testing corner cases
	Describe("Check splitting functionality", func() {
		Context("Valid DID", func() {
			It("should return expected method, namespace and DID", func() {
				method, namespace, id := MustSplitDID("did:cheqd:mainnet:qqqqqqqqqqqqqqqq")
				Expect(method).To(Equal("cheqd"))
				Expect(namespace).To(Equal("mainnet"))
				Expect(id).To(Equal("qqqqqqqqqqqqqqqq"))
			})
		})

		Context("Not valid DID string at all", func() {
			It("should panic", func() {
				panicDID := "Not 	"
				Expect(func() {
					MustSplitDID(panicDID)
				}).To(Panic())
			})
		})
	})

	DescribeTable("Check DID splitting and joining",

		func(did string) {
			method, namespace, id := MustSplitDID(did)
			Expect(did).To(Equal(JoinDID(method, namespace, id)))
		},
		Entry("Full DID", "did:cheqd:testnet:123456789abcdefg"),
		Entry("Without namespace", "did:cheqd:123456789abcdefg"),
		Entry("Not cheqd method", "did:NOTcheqd:123456789abcdefg"),
		Entry("32-symbols ID", "did:NOTcheqd:123456789abcdefg123456789abcdefg"),
	)
})
