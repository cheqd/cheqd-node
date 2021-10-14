package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgWriteRequest{}

func NewMsgWriteRequest(data *types.Any, metadata map[string]string, signatures map[string]string) *MsgWriteRequest {
	return &MsgWriteRequest{
		Data:       data,
		Metadata:   metadata,
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
	bz := ModuleCdc.MustMarshal(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWriteRequest) ValidateBasic() error {
	if msg.Data == nil {
		return ErrBadRequest.Wrap("Invalid Data: it is required")
	}

	if len(msg.Data.TypeUrl) == 0 || len(msg.Data.Value) == 0 {
		return ErrBadRequest.Wrap("Invalid Data: it cannot be empty")
	}

	if len(msg.Signatures) == 0 {
		return ErrBadRequest.Wrap("Invalid Signatures: it is required")
	}

	return nil
}
