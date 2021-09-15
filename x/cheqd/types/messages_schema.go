package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateSchema{}

func NewMsgCreateSchema(name string, version string, attr_names string) *MsgCreateSchema {
	return &MsgCreateSchema{
		Name:      name,
		Version:   version,
		AttrNames: attr_names,
	}
}

func (msg *MsgCreateSchema) Route() string {
	return RouterKey
}

func (msg *MsgCreateSchema) Type() string {
	return "CreateSchema"
}

func (msg *MsgCreateSchema) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateSchema) GetSignBytes() []byte {
	return []byte{}
}

func (msg *MsgCreateSchema) ValidateBasic() error {
	return nil
}
