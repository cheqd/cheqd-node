package types_test

import (
	. "github.com/cheqd/cheqd-node/x/did/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update DID Payload Validation tests", func() {
	type TestCaseUUIDDidStruct struct {
		inputID    string
		expectedID string
	}

	DescribeTable("UUID validation tests", func(testCase TestCaseUUIDDidStruct) {
		inputMsg := MsgUpdateDidDocPayload{
			Id:             testCase.inputID,
			Authentication: []string{testCase.inputID + "#key1"},
			VerificationMethod: []*VerificationMethod{
				{
					Id:                     testCase.inputID + "#key1",
					VerificationMethodType: Ed25519VerificationKey2020{}.Type(),
					Controller:             testCase.inputID,
				},
			},
			VersionId: "1234567890",
		}
		expectedMsg := MsgUpdateDidDocPayload{
			Id:             testCase.expectedID,
			Authentication: []string{testCase.expectedID + "#key1"},
			VerificationMethod: []*VerificationMethod{
				{
					Id:                     testCase.expectedID + "#key1",
					VerificationMethodType: Ed25519VerificationKey2020{}.Type(),
					Controller:             testCase.expectedID,
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
				inputID:    "did:cheqd:testnet:zABCDEFG123456789abcd",
				expectedID: "did:cheqd:testnet:zABCDEFG123456789abcd",
			}),

		Entry(
			"Mixed case UUID",
			TestCaseUUIDDidStruct{
				inputID:    "did:cheqd:testnet:BAbbba14-f294-458a-9b9c-474d188680fd",
				expectedID: "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Low case UUID",
			TestCaseUUIDDidStruct{
				inputID:    "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
				expectedID: "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Upper case UUID",
			TestCaseUUIDDidStruct{
				inputID:    "did:cheqd:testnet:A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
				expectedID: "did:cheqd:testnet:a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			}),
	)
})
