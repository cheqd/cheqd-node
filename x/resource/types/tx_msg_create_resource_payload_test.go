package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgCreateResourcePayloadValidation(t *testing.T) {
	cases := []struct {
		name     string
		struct_  *MsgCreateResourcePayload
		isValid  bool
		errorMsg string
	}{
		{
			name: "positive",
			struct_: &MsgCreateResourcePayload{
				CollectionId: "123456789abcdefg",
				Id:           "ba62c728-cb15-498b-8e9e-9259cc242186",
				Name:         "Test Resource",
				ResourceType: "CL-Schema",
				MediaType:    "application/json",
				Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			isValid: true,
		},
		{
			name: "negative resource type",
			struct_: &MsgCreateResourcePayload{
				CollectionId: "123456789abcdefg",
				Id:           "ba62c728-cb15-498b-8e9e-9259cc242186",
				Name:         "Test Resource",
				ResourceType: "Not-CL-Schema",
				MediaType:    "image/png",
				Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			isValid:  false,
			errorMsg: "resource_type: Not-CL-Schema resource type is not allowed. Only CL-Schema, JSONSchema2020.",
		},
		{
			name: "negative mime type",
			struct_: &MsgCreateResourcePayload{
				CollectionId: "123456789abcdefg",
				Id:           "ba62c728-cb15-498b-8e9e-9259cc242186",
				Name:         "Test Resource",
				ResourceType: "CL-Schema",
				MediaType:    "text/data",
				Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			isValid:  false,
			errorMsg: "media_type: text/data mime type is not allowed. Only application/json, application/octet-stream, text/plain, image/apng, image/avif, image/gif, image/jpeg, image/png, image/svg+xml, image/webp.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.struct_.Validate()

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, err.Error(), tc.errorMsg)
			}
		})
	}
}
