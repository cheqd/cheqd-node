package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-playground/validator/v10"
)

var _ sdk.Msg = &MsgCreateDid{}

func NewMsgCreateDid(payload *MsgCreateDidPayload, signatures []*SignInfo) *MsgCreateDid {
	return &MsgCreateDid{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgCreateDid) Route() string {
	return RouterKey
}

func (msg *MsgCreateDid) Type() string {
	return "MsgCreateDid"
}

func (msg *MsgCreateDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshal(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDid) ValidateBasic() error {
	validate := validator.New()

	if err := validate.Struct(msg); err != nil {
		return ErrBadRequest.Wrapf("basic validation failed: %s", err.Error())
	}

	return nil
}
