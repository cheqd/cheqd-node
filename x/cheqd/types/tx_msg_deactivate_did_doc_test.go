package types_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("Message for DID updating", func() {
	type TestCaseMsgDeactivateDID struct {
		msg      *MsgDeactivateDidDoc
		isValid  bool
		errorMsg string
	}

	DescribeTable("Tests for message for DID deactivating", func(testCase TestCaseMsgDeactivateDID) {
		err := testCase.msg.ValidateBasic()

		if testCase.isValid {
			Expect(err).To(BeNil())
		} else {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(testCase.errorMsg))
		}
	},

		Entry(
			"All fields are set properly",
			TestCaseMsgDeactivateDID{
				msg: &MsgDeactivateDidDoc{
					Payload: &MsgDeactivateDidDocPayload{
						Id: "did:cheqd:testnet:zABCDEFG123456789abcd",
					},
					Signatures: nil,
				},
				isValid: true,
			}),

		Entry(
			"Negative: Invalid DID Method",
			TestCaseMsgDeactivateDID{
				msg: &MsgDeactivateDidDoc{
					Payload: &MsgDeactivateDidDocPayload{
						Id: "did:cheqdttt:testnet:zABCDEFG123456789abcd",
					},
					Signatures: nil,
				},
				isValid:  false,
				errorMsg: "payload: (id: did method must be: cheqd.).: basic validation failed",
			}),

		Entry(
			"Negative: Id is required",
			TestCaseMsgDeactivateDID{
				msg: &MsgDeactivateDidDoc{
					Payload:    &MsgDeactivateDidDocPayload{},
					Signatures: nil,
				},
				isValid:  false,
				errorMsg: "payload: (id: cannot be blank.).: basic validation failed",
			}),
	)
})
