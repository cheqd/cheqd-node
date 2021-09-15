package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateCredDef{}

func NewMsgCreateCredDef(schema_id string, tag string, signature_type string, value string) *MsgCreateCredDef {
	return &MsgCreateCredDef{
		SchemaId:      schema_id,
		Tag:           tag,
		SignatureType: signature_type,
		Value:         value,
	}
}

func (msg *MsgCreateCredDef) Route() string {
	return RouterKey
}

func (msg *MsgCreateCredDef) Type() string {
	return "CreateCredDef"
}

func (msg *MsgCreateCredDef) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateCredDef) GetSignBytes() []byte {
	return []byte{}
}

func (msg *MsgCreateCredDef) ValidateBasic() error {
	return nil
}
