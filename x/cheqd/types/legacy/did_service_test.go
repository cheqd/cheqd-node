package legacy_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("Service tests", func() {
	type TestCaseServiceStruct struct {
		service           *Service
		baseDid           string
		allowedNamespaces []string
		isValid           bool
		errorMsg          string
	}

	DescribeTable("Service Validation tests", func(testCase TestCaseServiceStruct) {
		err := testCase.service.Validate(testCase.baseDid, testCase.allowedNamespaces)

		if testCase.isValid {
			Expect(err).To(BeNil())
		} else {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(testCase.errorMsg))
		}
	},

		Entry(
			"Positive case",
			TestCaseServiceStruct{
				service: &Service{
					Id:              "did:cheqd:aABCDEFG123456789abcd#service1",
					Type:            "DIDCommMessaging",
					ServiceEndpoint: "endpoint",
				},
				baseDid:           "did:cheqd:aABCDEFG123456789abcd",
				allowedNamespaces: []string{""},
				isValid:           true,
				errorMsg:          "",
			}),

		Entry(
			"Namespace is not allowed",
			TestCaseServiceStruct{
				service: &Service{
					Id:              "did:cheqd:zABCDEFG123456789abcd#service1",
					Type:            "DIDCommMessaging",
					ServiceEndpoint: "endpoint",
				},
				allowedNamespaces: []string{"mainnet"},
				isValid:           false,
				errorMsg:          "id: did namespace must be one of: mainnet.",
			}),

		Entry(
			"Base DID is not the same as in id",
			TestCaseServiceStruct{
				service: &Service{
					Id:              "did:cheqd:zABCDEFG123456789abcd#service1",
					Type:            "DIDCommMessaging",
					ServiceEndpoint: "endpoint",
				},
				baseDid:  "did:cheqd:zABCDEFG987654321abcd",
				isValid:  false,
				errorMsg: "id: must have prefix: did:cheqd:zABCDEFG987654321abcd",
			}),
	)
})
