package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateAttrib{}

func NewMsgCreateAttrib(did string, raw string) *MsgCreateAttrib {
	return &MsgCreateAttrib{
		Did: did,
		Raw: raw,
	}
}

func (msg *MsgCreateAttrib) Route() string {
	return RouterKey
}

func (msg *MsgCreateAttrib) Type() string {
	return "CreateAttrib"
}

func (msg *MsgCreateAttrib) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateAttrib) GetSignBytes() []byte {
	return []byte{}
}

func (msg *MsgCreateAttrib) ValidateBasic() error {
	return nil
}

var _ sdk.Msg = &MsgUpdateAttrib{}

func NewMsgUpdateAttrib(did string, raw string) *MsgUpdateAttrib {
	return &MsgUpdateAttrib{
		Did: did,
		Raw: raw,
	}
}

func (msg *MsgUpdateAttrib) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAttrib) Type() string {
	return "UpdateAttrib"
}

func (msg *MsgUpdateAttrib) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgUpdateAttrib) GetSignBytes() []byte {
	return []byte{}
}

func (msg *MsgUpdateAttrib) ValidateBasic() error {
	return nil
}
