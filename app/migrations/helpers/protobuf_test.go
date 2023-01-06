package helpers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable(
	"Test GenerateEd25519VerificationKey2020VerificationMaterial",
	func(v1PubKey string, v2PubKey string) {
		key, err := GenerateEd25519VerificationKey2020VerificationMaterial(v1PubKey)
		Expect(err).To(BeNil())
		Expect(key).To(Equal(v2PubKey))
	},

	Entry("Valid: General conversion", "zDw21irq4wBfyTvxAG9L8PQj6b79iyTyzyV6XVj9SfyRR", "z6MksPH4K75WGjASaRnrwiHyEWH6QgRaPMEMfW1TL17TbCCo"),
	// Mainnet case
	Entry("Valid: Real case", "zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkYsWCo7fztHtepn", "z6Mkta7joRuvDh7UnoESdgpr9dDUMh5LvdoECDi3WGrJoscA"),
)
