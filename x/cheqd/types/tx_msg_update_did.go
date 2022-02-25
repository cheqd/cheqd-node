package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateDid{}

func NewMsgUpdateDid(payload *MsgUpdateDidPayload, signatures []*SignInfo) *MsgUpdateDid {
	return &MsgUpdateDid{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgUpdateDid) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDid) Type() string {
	return "WriteRequest"
}

func (msg *MsgUpdateDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgUpdateDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshal(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDid) ValidateBasic() error {
	validate, err := BuildValidator(nil)
	if err != nil {
		return ErrInternal.Wrapf("can't init validator: %s", err.Error())
	}

	if err := validate.Struct(msg); err != nil {
		return ErrBadRequest.Wrapf("basic validation failed: %s", err.Error())
	}

	return nil
}
