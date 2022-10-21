package types_test

import (
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TxMsgCreateResourcePayload", func() {
	var msg *resourcetypes.MsgCreateResourcePayload
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
})
