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
				ResourceType: "",
				Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			isValid:  false,
			errorMsg: "resource_type: cannot be blank.",
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
