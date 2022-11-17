package utils_test

import (
	. "github.com/cheqd/cheqd-node/x/did/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("URI validation", func() {
	DescribeTable("ValidateURI",

		func(isValid bool, uri string) {
			err_ := ValidateURI(uri)
			if isValid {
				Expect(err_).To(BeNil())
			} else {
				Expect(err_).ToNot(BeNil())
			}
		},

		Entry("Valid: General http URI path", true, "http://a.com/a/b/c/d/?query=123#fragment=another_part"),
		Entry("Valid: General https URI path", true, "https://a.com/a/b/c/d/?query=123#fragment=another_part"),
		Entry("Valid: only alphabet symbols", true, "SomeAnotherPath"),
	)
})
