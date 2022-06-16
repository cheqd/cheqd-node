package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgUpdateDidValidation(t *testing.T) {
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
				MimeType:     "text/plain",
				Data:         []byte {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			isValid: true,
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
