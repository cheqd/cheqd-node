package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateSchema{}

func NewMsgCreateSchema(id string, name string, version string, attrNames []string) *MsgCreateSchema {
	return &MsgCreateSchema{
		Id:        id,
		Name:      name,
		Version:   version,
		AttrNames: attrNames,
	}
}

func (msg *MsgCreateSchema) Route() string {
	return RouterKey
}

func (msg *MsgCreateSchema) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateSchema) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateSchema) ValidateBasic() error {
	return nil
}
