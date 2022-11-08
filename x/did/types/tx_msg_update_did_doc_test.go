package types_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/did/types"
)

var _ = Describe("Message for DID updating", func() {
	type TestCaseMsgUpdateDID struct {
		msg      *MsgUpdateDidDoc
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
				msg: &MsgUpdateDidDoc{
					Payload: &MsgUpdateDidDocPayload{
						Id: "did:cheqd:testnet:zABCDEFG123456789abcd",
						VerificationMethod: []*VerificationMethod{
							{
								Id:                   "did:cheqd:testnet:zABCDEFG123456789abcd#key1",
								Type:                 "Ed25519VerificationKey2020",
								Controller:           "did:cheqd:testnet:zABCDEFG123456789abcd",
								VerificationMaterial: ValidEd25519VerificationMaterial,
							},
						},
						Authentication:    []string{"did:cheqd:testnet:zABCDEFG123456789abcd#key1", "did:cheqd:testnet:zABCDEFG123456789abcd#aaa"},
						VersionId:         "version1",
						PreviousVersionId: uuid.NewString(),
					},
					Signatures: nil,
				},
				isValid: true,
			}),

		Entry(
			"IDs are duplicated",
			TestCaseMsgUpdateDID{
				msg: &MsgUpdateDidDoc{
					Payload: &MsgUpdateDidDocPayload{
						Id: "did:cheqd:testnet:zABCDEFG123456789abcd",
						VerificationMethod: []*VerificationMethod{
							{
								Id:                   "did:cheqd:testnet:zABCDEFG123456789abcd#key1",
								Type:                 "Ed25519VerificationKey2020",
								Controller:           "did:cheqd:testnet:zABCDEFG123456789abcd",
								VerificationMaterial: ValidEd25519VerificationMaterial,
							},
						},
						Authentication: []string{"did:cheqd:testnet:zABCDEFG123456789abcd#key1", "did:cheqd:testnet:zABCDEFG123456789abcd#key1"},
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
				msg: &MsgUpdateDidDoc{
					Payload: &MsgUpdateDidDocPayload{
						Id: "did:cheqd:testnet:zABCDEFG123456789abcd",
						VerificationMethod: []*VerificationMethod{
							{
								Id:                   "did:cheqd:testnet:zABCDEFG123456789abcd#key1",
								Type:                 "Ed25519VerificationKey2020",
								Controller:           "did:cheqd:testnet:zABCDEFG123456789abcd",
								VerificationMaterial: ValidEd25519VerificationMaterial,
							},
						},
						Authentication:    []string{"did:cheqd:testnet:zABCDEFG123456789abcd#key1", "did:cheqd:testnet:zABCDEFG123456789abcd#aaa"},
						PreviousVersionId: uuid.NewString(),
					},
					Signatures: nil,
				},
				isValid:  false,
				errorMsg: "payload: (version_id: cannot be blank.).: basic validation failed",
			}),
	)
})
