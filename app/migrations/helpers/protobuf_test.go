package helpers

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	expecteRFC3339, _        = time.Parse(time.RFC3339, "2022-09-06T16:19:39Z")
	expecteRFC3339Nano, _    = time.Parse(time.RFC3339Nano, "2022-09-06T16:19:39.464251406Z")
	expectedOldTimeFormat, _ = time.Parse(OldTimeFormat, "2022-02-22 13:32:19.464251406 +0000 UTC")
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

var _ = DescribeTable(
	"Test MustParseFromStringTimeToGoTime",
	func(inputS string, outputTT time.Time) {
		timeTime := MustParseFromStringTimeToGoTime(inputS)
		Expect(timeTime).To(Equal(outputTT))
	},

	Entry("Valid: General conversion RFC3339", "2022-09-06T16:19:39Z", expecteRFC3339),
	Entry("Valid: General conversion RFC3339Nano", "2022-09-06T16:19:39.464251406Z", expecteRFC3339Nano),
	Entry("Valid: General conversion OldTimeFormat", "2022-02-22 13:32:19.464251406 +0000 UTC", expectedOldTimeFormat),
)
