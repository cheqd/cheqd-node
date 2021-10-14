package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMsgWriteRequestValidation(t *testing.T) {
	cases := []struct {
		valid  bool
		msg    *MsgWriteRequest
		errMsg string
	}{
		{true, NewMsgWriteRequest(&types.Any{TypeUrl: "1", Value: []byte{1}}, nil, map[string]string{"foo": "bar"}), ""},
		{false, NewMsgWriteRequest(&types.Any{TypeUrl: "1"}, nil, map[string]string{"foo": "bar"}), "Invalid Data: it cannot be empty: bad request"},
		{false, NewMsgWriteRequest(&types.Any{Value: []byte{1}}, nil, map[string]string{"foo": "bar"}), "Invalid Data: it cannot be empty: bad request"},
		{false, NewMsgWriteRequest(nil, nil, nil), "Invalid Data: it is required: bad request"},
		{false, NewMsgWriteRequest(&types.Any{TypeUrl: "1", Value: []byte{1}}, nil, nil), "Invalid Signatures: it is required: bad request"},
	}

	for _, tc := range cases {
		err := tc.msg.ValidateBasic()

		if tc.valid {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
			require.Equal(t, tc.errMsg, err.Error())
		}
	}
}
