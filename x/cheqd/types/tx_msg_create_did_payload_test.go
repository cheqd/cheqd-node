package types_test

import (
	. "github.com/cheqd/cheqd-node/x/cheqd/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create DID Payload Validation tests", func() {

	type TestCaseUUIDDidStruct struct {
		inputId    string
		expectedId string
	}

	DescribeTable("UUID validation tests", func(testCase TestCaseUUIDDidStruct) {

		inputMsg := MsgCreateDidPayload{
			Id:             testCase.inputId,
			Authentication: []string{testCase.inputId + "#key1"},
			VerificationMethod: []*VerificationMethod{
				{
					Id:         testCase.inputId + "#key1",
					Type:       Ed25519VerificationKey2020,
					Controller: testCase.inputId,
				},
			},
		}
		expectedMsg := MsgCreateDidPayload{
			Id:             testCase.expectedId,
			Authentication: []string{testCase.expectedId + "#key1"},
			VerificationMethod: []*VerificationMethod{
				{
					Id:         testCase.expectedId + "#key1",
					Type:       Ed25519VerificationKey2020,
					Controller: testCase.expectedId,
				},
			},
		}
		inputMsg.Normalize()
		Expect(inputMsg).To(Equal(expectedMsg))
	},

		Entry(
			"base58 identifier - not changed",
			TestCaseUUIDDidStruct{
				inputId:    "did:cheqd:testnet:aaaaaaaaaaaaaaaa",
				expectedId: "did:cheqd:testnet:aaaaaaaaaaaaaaaa",
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
