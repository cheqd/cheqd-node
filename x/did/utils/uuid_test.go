package utils_test

import (
	. "github.com/cheqd/cheqd-node/x/did/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UUID validation", func() {
	DescribeTable("ValidateUUID",

		func(uri string, isValid bool) {
			err_ := ValidateUUID(uri)
			if isValid {
				Expect(err_).To(BeNil())
			} else {
				Expect(err_).ToNot(BeNil())
			}
		},

		Entry("Valid: General UUID", "42d9c704-ecb0-11ec-8ea0-0242ac120002", true),
		Entry("Not Valid: wrong format", "not uuid", false),
		Entry("Valid: Another general UUID", "e1cdbc10-858c-4d7d-8a4b-5d45e90a81b3", true),
		Entry("Not Valid: Unsupported symbol {", "{42d9c704-ecb0-11ec-8ea0-0242ac120002}", false),
		Entry("Not Valid: Unexpeced prefix", "urn:uuid:42d9c704-ecb0-11ec-8ea0-0242ac120002", false),
		Entry("Not Valid: Wrong format", "42d9c704ecb011ec8ea00242ac120002", false),
	)
})
