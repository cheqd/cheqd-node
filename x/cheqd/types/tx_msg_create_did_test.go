package types_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("Message for DID creation", func() {
	type TestCaseMsgCreateDID struct {
		msg *MsgCreateDid
		isValid bool
		errorMsg string
	}

	DescribeTable("Tests for message for DID creation", func(testCase TestCaseMsgCreateDID) {
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
		TestCaseMsgCreateDID{
			msg: &MsgCreateDid{
				Payload: &MsgCreateDidPayload{
					Id: "did:cheqd:testnet:123456789abcdefg",
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 "did:cheqd:testnet:123456789abcdefg#key1",
							Type:               "Ed25519VerificationKey2020",
							Controller:         "did:cheqd:testnet:123456789abcdefg",
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
					Authentication: []string{"did:cheqd:testnet:123456789abcdefg#key1", "did:cheqd:testnet:123456789abcdefg#aaa"},
				},
				Signatures: nil,
			},
			isValid: true}),

	Entry(
		"IDs are duplicated", 
		TestCaseMsgCreateDID{
			msg: &MsgCreateDid{
				Payload: &MsgCreateDidPayload{
					Id: "did:cheqd:testnet:123456789abcdefg",
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 "did:cheqd:testnet:123456789abcdefg#key1",
							Type:               "Ed25519VerificationKey2020",
							Controller:         "did:cheqd:testnet:123456789abcdefg",
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
					Authentication: []string{"did:cheqd:testnet:123456789abcdefg#key1", "did:cheqd:testnet:123456789abcdefg#key1"},
				},
				Signatures: nil,
			},
			isValid: false,
			errorMsg: "payload: (authentication: there should be no duplicates.).: basic validation failed",
		}),
	)
})
