package types_test

import (
. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

. "github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("Message for DID creation", func() {

	var struct_           *MsgCreateDid
	var isValid           bool
	var errorMsg          string

	BeforeEach(func() {
		struct_ = &MsgCreateDid{}
		isValid = false
		errorMsg = ""
	})

	AfterEach(func() {
		err := struct_.ValidateBasic()

			if isValid {
				Expect(err).To(BeNil())
			} else {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(errorMsg))
			}
	})

	When("all fields are set properly", func() {
		It("Will pass", func() {
			struct_ = &MsgCreateDid{
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
			}
			isValid = true
		})
	})

	When("IDs are duplicated", func() {
		It("should fail the validation", func() {
			struct_ = &MsgCreateDid{
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
			}
			isValid = false
			errorMsg = "payload: (authentication: there should be no duplicates.).: basic validation failed"
		})

	})
})
