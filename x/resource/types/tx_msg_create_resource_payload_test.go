package types_test

import (
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TxMsgCreateResourcePayload", func() {
	var msg *resourcetypes.MsgCreateResourcePayload
	type TestCaseUUIDDidStruct struct {
		inputCollectionId    string
		inputId              string
		expectedId           string
		expectedCollectionId string
	}

	Describe("Validate", func() {
		Context("Valid: MsgCreateResourcePayload", func() {
			It("should return nil if the message is valid", func() {
				msg = &resourcetypes.MsgCreateResourcePayload{
					CollectionId: "123456789abcdefg",
					Id:           "ba62c728-cb15-498b-8e9e-9259cc242186",
					Name:         "Test Resource",
					ResourceType: "CL-Schema",
					Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				}
				Expect(msg.Validate()).To(BeNil())
			})
		})

		Context("Invalid: MsgCreateResourcePayload", func() {
			It("should return error if the resource type is empty", func() {
				msg = &resourcetypes.MsgCreateResourcePayload{
					CollectionId: "123456789abcdefg",
					Id:           "ba62c728-cb15-498b-8e9e-9259cc242186",
					Name:         "Test Resource",
					ResourceType: "",
					Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				}
				Expect(msg.Validate().Error()).To(Equal("resource_type: cannot be blank."))
			})
		})
	})

	DescribeTable("UUID validation tests", func(testCase TestCaseUUIDDidStruct) {
		inputMsg := resourcetypes.MsgCreateResourcePayload{
			CollectionId: testCase.inputCollectionId,
			Id:           testCase.inputId,
			Name:         "Test Resource",
			ResourceType: "",
			Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
		expectedMsg := resourcetypes.MsgCreateResourcePayload{
			CollectionId: testCase.expectedCollectionId,
			Id:           testCase.expectedId,
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
				inputCollectionId:    "aaaaaaaaaaaaaaaa",
				inputId:              "aaaaaaaaaaaaaaaa",
				expectedCollectionId: "aaaaaaaaaaaaaaaa",
				expectedId:           "aaaaaaaaaaaaaaaa",
			}),

		Entry(
			"Mixed case UUID",
			TestCaseUUIDDidStruct{
				inputCollectionId:    "BAbbba14-f294-458a-9b9c-474d188680fd",
				inputId:              "BAbbba14-f294-458a-9b9c-474d188680fd",
				expectedCollectionId: "babbba14-f294-458a-9b9c-474d188680fd",
				expectedId:           "babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Low case UUID",
			TestCaseUUIDDidStruct{
				inputCollectionId:    "babbba14-f294-458a-9b9c-474d188680fd",
				inputId:              "babbba14-f294-458a-9b9c-474d188680fd",
				expectedCollectionId: "babbba14-f294-458a-9b9c-474d188680fd",
				expectedId:           "babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Upper case UUID",
			TestCaseUUIDDidStruct{
				inputCollectionId:    "A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
				inputId:              "A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
				expectedCollectionId: "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
				expectedId:           "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			}),
	)
})
