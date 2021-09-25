package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgWriteRequest{}

func NewMsgWriteRequest(data *types.Any, authors []string, signatures map[string]string) *MsgWriteRequest {
	return &MsgWriteRequest{
		Data:       data,
		Authors:    authors,
		Signatures: signatures,
	}
}

func (msg *MsgWriteRequest) Route() string {
	return RouterKey
}

func (msg *MsgWriteRequest) Type() string {
	return "WriteRequest"
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
