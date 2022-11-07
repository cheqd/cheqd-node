package types_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/did/types"
)

var _ = Describe("SignInfo tests", func() {
	type TestCaseSignInfoStruct struct {
		si                SignInfo
		allowedNamespaces []string
		isValid           bool
		errorMsg          string
	}

	DescribeTable("SignInfo validation tests", func(testCase TestCaseSignInfoStruct) {
		err := testCase.si.Validate(testCase.allowedNamespaces)

		if testCase.isValid {
			Expect(err).To(BeNil())
		} else {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(testCase.errorMsg))
		}
	},

		Entry(
			"Positive case",
			TestCaseSignInfoStruct{
				si: SignInfo{
					VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#method1",
					Signature:            []byte("aaa="),
				},
				isValid:  true,
				errorMsg: "",
			}),

		Entry(
			"Namespace is not allowed",
			TestCaseSignInfoStruct{
				si: SignInfo{
					VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#service1",
					Signature:            []byte("DIDCommMessaging"),
				},
				allowedNamespaces: []string{"mainnet"},
				isValid:           false,
				errorMsg:          "verification_method_id: did namespace must be one of: mainnet.",
			}),
	)
})

var _ = Describe("Full SignInfo duplicates tests", func() {
	type TestCaseSignInfosStruct struct {
		signInfos []*SignInfo
		isValid   bool
	}

	DescribeTable("SignInfo duplicates tests", func(testCase TestCaseSignInfosStruct) {
		res_ := IsUniqueSignInfoList(testCase.signInfos)
		Expect(res_).To(Equal(testCase.isValid))
	},

		Entry(
			"Signatures are different",
			TestCaseSignInfosStruct{
				signInfos: []*SignInfo{
					{
						VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#method1",
						Signature:            []byte("aaa="),
					},
					{
						VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#method1",
						Signature:            []byte("bbb="),
					},
				},
				isValid: true,
			}),

		Entry(
			"All fields are different",
			TestCaseSignInfosStruct{
				signInfos: []*SignInfo{
					{
						VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#method1",
						Signature:            []byte("aaa="),
					},
					{
						VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#method1",
						Signature:            []byte("bbb="),
					},
				},
				isValid: true,
			}),

		Entry(
			"All fields are the same",
			TestCaseSignInfosStruct{
				signInfos: []*SignInfo{
					{
						VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#method1",
						Signature:            []byte("aaa="),
					},
					{
						VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#method1",
						Signature:            []byte("aaa="),
					},
				},
				isValid: false,
			}),

		Entry(
			"All fields are the same and more elments",
			TestCaseSignInfosStruct{
				signInfos: []*SignInfo{
					{
						VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#method1",
						Signature:            []byte("aaa="),
					},
					{
						VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#method1",
						Signature:            []byte("aaa="),
					},
					{
						VerificationMethodId: "did:cheqd:zABCDEFG123456789abcd#method1",
						Signature:            []byte("aaa="),
					},
				},
				isValid: false,
			}),
	)
})
