package types_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("SignInfo validation tests", func() {

	var struct_           SignInfo
	var allowedNamespaces []string
	var isValid           bool
	var errorMsg          string

	BeforeEach(func() {
		struct_ = SignInfo{}
		allowedNamespaces = []string{}
		isValid = false
		errorMsg = ""
	})

	AfterEach(func() {
		err := struct_.Validate(allowedNamespaces)

			if isValid {
				Expect(err).To(BeNil())
			} else {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(errorMsg))
			}
	})

	It("Positive case", func() {
		struct_ = SignInfo{
			VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
			Signature:            "aaa=",
		}
		isValid = true
		errorMsg = ""
	})

	When("namespace is not allowed", func() {
		It("should fail", func() {
			struct_ = SignInfo{
				VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#service1",
				Signature:            "DIDCommMessaging",
			}
			allowedNamespaces = []string{"mainnet"}
			isValid = false
			errorMsg = "verification_method_id: did namespace must be one of: mainnet."
		})
	})

	When("signature is not valid base64 string", func() {
		It("should fail", func() {
			struct_ = SignInfo{
				VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#service1",
				Signature:            "!@#",
			}
			isValid = false
			errorMsg = "signature: must be encoded in Base64."
		})
	})

})

var _ = Describe("Full SignInfo duplicates tests", func() {
	var structs_ []*SignInfo
	var isValid  bool

	BeforeEach(func() {
		structs_ = []*SignInfo{}
		isValid = false
	})

	AfterEach(func() {
		res_ := IsUniqueSignInfoList(structs_)
		Expect(res_).To(Equal(isValid))
	})

	When("signatures are different", func() {
		It("should pass", func() {
			structs_ = []*SignInfo{
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa=",
				},
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "bbb=",
				},
			}
			isValid  = true
		})
	})

	When("all fields are different", func() {
		It("should pass", func() {
			structs_ = []*SignInfo{
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa=",
				},
				{
					VerificationMethodId: "did:cheqd:bbbbbbbbbbbbbbbb#method1",
					Signature:            "bbb=",
				},
			}
			isValid = true
		})
	})

	When("all fields are the same", func() {
		It("should fail", func() {
			structs_ = []*SignInfo{
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa=",
				},
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa=",
				},
			}
			isValid = false
		})
	})

	When("all fields are the same and more elments", func() {
		It("should fail", func() {
			structs_ = []*SignInfo{
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa=",
				},
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa=",
				},
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa=",
				},
			}
			isValid = false
		})
	})
})
