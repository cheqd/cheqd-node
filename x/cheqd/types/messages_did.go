package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateDid{}

func NewMsgCreateDid(id string, verkey string, alias string) *MsgCreateDid {
	return &MsgCreateDid{
		Id:     id,
		Verkey: verkey,
		Alias:  alias,
	}
}

func (msg *MsgCreateDid) Route() string {
	return RouterKey
}

func (msg *MsgCreateDid) Type() string {
	return "CreateDid"
}

func (msg *MsgCreateDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateDid) GetSignBytes() []byte {
	return []byte{}
}

func (msg *MsgCreateDid) ValidateBasic() error {
	return nil
}

var _ sdk.Msg = &MsgUpdateDid{}

func NewMsgUpdateDid(id string, verkey string, alias string) *MsgUpdateDid {
	return &MsgUpdateDid{
		Id:     id,
		Verkey: verkey,
		Alias:  alias,
	}
}

func (msg *MsgUpdateDid) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDid) Type() string {
	return "UpdateDid"
}

func (msg *MsgUpdateDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgUpdateDid) GetSignBytes() []byte {
	return []byte{}
}

func (msg *MsgUpdateDid) ValidateBasic() error {
	return nil
}
