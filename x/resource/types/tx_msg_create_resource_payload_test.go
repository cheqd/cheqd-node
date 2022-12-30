package types_test

import (
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TxMsgCreateResourcePayload", func() {
	var msg *resourcetypes.MsgCreateResourcePayload
	type TestCaseUUIDDidStruct struct {
		inputCollectionID    string
		inputID              string
		expectedID           string
		expectedCollectionID string
	}

	Describe("Validate", func() {
		Context("Valid: MsgCreateResourcePayload", func() {
			It("should return nil if the message is valid", func() {
				msg = &resourcetypes.MsgCreateResourcePayload{
					CollectionId: "zABCDEFG123456789abcd",
					Id:           "ba62c728-cb15-498b-8e9e-9259cc242186",
					Name:         "Test Resource",
					Version:      "1.0",
					ResourceType: "CL-Schema",
					Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				}
				Expect(msg.Validate()).To(BeNil())
			})
		})

		Context("Invalid: MsgCreateResourcePayload", func() {
			It("should return error if the resource type is empty", func() {
				msg = &resourcetypes.MsgCreateResourcePayload{
					CollectionId: "zABCDEFG123456789abcd",
					Id:           "ba62c728-cb15-498b-8e9e-9259cc242186",
					Name:         "Test Resource",
					Version:      "1.0",
					ResourceType: "",
					Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				}
				Expect(msg.Validate().Error()).To(Equal("resource_type: cannot be blank."))
			})
		})
	})

	DescribeTable("UUID validation tests", func(testCase TestCaseUUIDDidStruct) {
		inputMsg := resourcetypes.MsgCreateResourcePayload{
			CollectionId: testCase.inputCollectionID,
			Id:           testCase.inputID,
			Name:         "Test Resource",
			ResourceType: "",
			Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
		expectedMsg := resourcetypes.MsgCreateResourcePayload{
			CollectionId: testCase.expectedCollectionID,
			Id:           testCase.expectedID,
			Name:         "Test Resource",
			ResourceType: "",
			Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
		inputMsg.Normalize()
		Expect(inputMsg).To(Equal(expectedMsg))
	},

		Entry(
			"base58 identifier - not changed",
			TestCaseUUIDDidStruct{
				inputCollectionID:    "zABCDEFG123456789abcd",
				inputID:              "8c614475-ec20-4ff2-bcf3-e4f28b849dbc",
				expectedCollectionID: "zABCDEFG123456789abcd",
				expectedID:           "8c614475-ec20-4ff2-bcf3-e4f28b849dbc",
			}),

		Entry(
			"Mixed case UUID",
			TestCaseUUIDDidStruct{
				inputCollectionID:    "BAbbba14-f294-458a-9b9c-474d188680fd",
				inputID:              "BAbbba14-f294-458a-9b9c-474d188680fd",
				expectedCollectionID: "babbba14-f294-458a-9b9c-474d188680fd",
				expectedID:           "babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Low case UUID",
			TestCaseUUIDDidStruct{
				inputCollectionID:    "babbba14-f294-458a-9b9c-474d188680fd",
				inputID:              "babbba14-f294-458a-9b9c-474d188680fd",
				expectedCollectionID: "babbba14-f294-458a-9b9c-474d188680fd",
				expectedID:           "babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Upper case UUID",
			TestCaseUUIDDidStruct{
				inputCollectionID:    "A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
				inputID:              "A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
				expectedCollectionID: "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
				expectedID:           "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			}),
	)
})
