package types_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("Message for DID updating", func() {
	type TestCaseMsgUpdateDID struct {
		msg      *MsgUpdateDid
		isValid  bool
		errorMsg string
	}

	DescribeTable("Tests for message for DID updating", func(testCase TestCaseMsgUpdateDID) {
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
			TestCaseMsgUpdateDID{
				msg: &MsgUpdateDid{
					Payload: &MsgUpdateDidPayload{
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
						VersionId:      "version1",
					},
					Signatures: nil,
				},
				isValid: true,
			}),

		Entry(
			"IDs are duplicated",
			TestCaseMsgUpdateDID{
				msg: &MsgUpdateDid{
					Payload: &MsgUpdateDidPayload{
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
						VersionId:      "version1",
					},
					Signatures: nil,
				},
				isValid:  false,
				errorMsg: "payload: (authentication: there should be no duplicates.).: basic validation failed",
			}),
		Entry(
			"VersionId is empty",
			TestCaseMsgUpdateDID{
				msg: &MsgUpdateDid{
					Payload: &MsgUpdateDidPayload{
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
				isValid:  false,
				errorMsg: "payload: (version_id: cannot be blank.).: basic validation failed",
			}),
	)
})
