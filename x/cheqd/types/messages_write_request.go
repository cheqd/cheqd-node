package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgWriteRequest{}

func NewMsgWriteRequest(data *types.Any, author string, signature string) *MsgWriteRequest {
	return &MsgWriteRequest{
		Data:      data,
		Author:    author,
		Signature: signature,
	}
}

func (msg *MsgWriteRequest) Route() string {
	return RouterKey
}

func (msg *MsgWriteRequest) Type() string {
	return "WriteRequestDef"
}

func (msg *MsgWriteRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgWriteRequest) GetSignBytes() []byte {
	return []byte{}
}

func (msg *MsgWriteRequest) ValidateBasic() error {
	return nil
}
