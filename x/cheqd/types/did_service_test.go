package types_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("DID Validation tests", func() {

	var struct_           Service
	var baseDid           string
	var allowedNamespaces []string
	var isValid           bool
	var errorMsg          string

	BeforeEach(func() {
		struct_ = Service{}
		baseDid = ""
		allowedNamespaces = []string{}
		isValid = false
		errorMsg = ""
	})

	AfterEach(func() {
		err := struct_.Validate(baseDid, allowedNamespaces)

			if isValid {
				Expect(err).To(BeNil())
			} else {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(errorMsg))
			}
	})

	It("Positive case", func() {
		struct_ = Service{
			Id:              "did:cheqd:aaaaaaaaaaaaaaaa#service1",
			Type:            "DIDCommMessaging",
			ServiceEndpoint: "endpoint",
		}
		baseDid = "did:cheqd:aaaaaaaaaaaaaaaa"
		allowedNamespaces = []string{""}
		isValid = true
		errorMsg = ""
	})

	When("Namespace is not allowed", func() {

		It("should fail", func() {
			struct_ = Service{
				Id:              "did:cheqd:aaaaaaaaaaaaaaaa#service1",
				Type:            "DIDCommMessaging",
				ServiceEndpoint: "endpoint",
			}
			allowedNamespaces = []string{"mainnet"}
			isValid = false
			errorMsg = "id: did namespace must be one of: mainnet."
		})
	})

	When("base DID is not the same as in id", func() {
		It("should fail", func() {
			struct_ = Service{
				Id:              "did:cheqd:aaaaaaaaaaaaaaaa#service1",
				Type:            "DIDCommMessaging",
				ServiceEndpoint: "endpoint",
			}
			baseDid = "did:cheqd:baaaaaaaaaaaaaab"
			isValid = false
			errorMsg = "id: must have prefix: did:cheqd:baaaaaaaaaaaaaab."
		})
	})
})
