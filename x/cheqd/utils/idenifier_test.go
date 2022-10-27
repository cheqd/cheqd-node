package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

var _ = Describe("Identifier Validation tests", func() {
	type TestCaseUUIDIdStruct struct {
		inputId    string
		expectedId string
	}

	DescribeTable("UUID test cases", func(testCase TestCaseUUIDIdStruct) {
		result := NormalizeIdentifier(testCase.inputId)
		Expect(result).To(Equal(testCase.expectedId))
	},

		Entry(
			"base58 identifier - not changed",
			TestCaseUUIDIdStruct{
				inputId:    "did:cheqd:testnet:aaaaaaaaaaaaaaaa",
				expectedId: "did:cheqd:testnet:aaaaaaaaaaaaaaaa",
			}),

		Entry(
			"Mixed case UUID",
			TestCaseUUIDIdStruct{
				inputId:    "did:cheqd:testnet:BAbbba14-f294-458a-9b9c-474d188680fd",
				expectedId: "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Low case UUID",
			TestCaseUUIDIdStruct{
				inputId:    "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
				expectedId: "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Upper case UUID",
			TestCaseUUIDIdStruct{
				inputId:    "did:cheqd:testnet:A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
				expectedId: "did:cheqd:testnet:a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			}),
	)

	DescribeTable("Id validation tests", func(isValid bool, id string) {
		isDid := IsValidID(id)

		if isValid {
			Expect(isDid).To(BeTrue())
		} else {
			Expect(isDid).To(BeFalse())
		}
	},

		Entry("Base58 string, 16 symbols", true, "123456789abcdefg"),
		Entry("Base58 string, 32 symbols", true, "123456789abcdefg123456789abcdefg"),
		Entry("UUID string", true, "3b9b8eec-5b5d-4382-86d8-9185126ff130"),
		Entry("Too short", false, "sdf"),
		Entry("Unexpected :", false, "sdf:sdf"),
		Entry("Too short", false, "12345"),
	)
})
