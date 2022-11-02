package types_test

import (
	. "github.com/cheqd/cheqd-node/x/cheqd/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update DID Payload Validation tests", func() {
	type TestCaseUUIDDidStruct struct {
		inputId    string
		expectedId string
	}

	DescribeTable("UUID validation tests", func(testCase TestCaseUUIDDidStruct) {
		inputMsg := MsgUpdateDidDocPayload{
			Id:             testCase.inputId,
			Authentication: []string{testCase.inputId + "#key1"},
			VerificationMethod: []*VerificationMethod{
				{
					Id:         testCase.inputId + "#key1",
					Type:       Ed25519VerificationKey2020{}.Type(),
					Controller: testCase.inputId,
				},
			},
			VersionId: "1234567890",
		}
		expectedMsg := MsgUpdateDidDocPayload{
			Id:             testCase.expectedId,
			Authentication: []string{testCase.expectedId + "#key1"},
			VerificationMethod: []*VerificationMethod{
				{
					Id:         testCase.expectedId + "#key1",
					Type:       Ed25519VerificationKey2020{}.Type(),
					Controller: testCase.expectedId,
				},
			},
			VersionId: "1234567890",
		}
		inputMsg.Normalize()
		Expect(inputMsg).To(Equal(expectedMsg))
	},

		Entry(
			"base58 identifier - not changed",
			TestCaseUUIDDidStruct{
				inputId:    "did:cheqd:testnet:zABCDEFG123456789abcd",
				expectedId: "did:cheqd:testnet:zABCDEFG123456789abcd",
			}),

		Entry(
			"Mixed case UUID",
			TestCaseUUIDDidStruct{
				inputId:    "did:cheqd:testnet:BAbbba14-f294-458a-9b9c-474d188680fd",
				expectedId: "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Low case UUID",
			TestCaseUUIDDidStruct{
				inputId:    "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
				expectedId: "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Upper case UUID",
			TestCaseUUIDDidStruct{
				inputId:    "did:cheqd:testnet:A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
				expectedId: "did:cheqd:testnet:a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			}),
	)
})
