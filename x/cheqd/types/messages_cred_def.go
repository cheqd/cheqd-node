package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateCredDef{}

func NewMsgCreateCredDef(creator string, schema_id string, tag string, signature_type string, value string) *MsgCreateCredDef {
	return &MsgCreateCredDef{
		Creator:        creator,
		Schema_id:      schema_id,
		Tag:            tag,
		Signature_type: signature_type,
		Value:          value,
	}
}

func (msg *MsgCreateCredDef) Route() string {
	return RouterKey
}

func (msg *MsgCreateCredDef) Type() string {
	return "CreateCredDef"
}

func (msg *MsgCreateCredDef) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateCredDef) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateCredDef) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateCredDef{}

func NewMsgUpdateCredDef(creator string, id uint64, schema_id string, tag string, signature_type string, value string) *MsgUpdateCredDef {
	return &MsgUpdateCredDef{
		Id:             id,
		Creator:        creator,
		Schema_id:      schema_id,
		Tag:            tag,
		Signature_type: signature_type,
		Value:          value,
	}
}

func (msg *MsgUpdateCredDef) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCredDef) Type() string {
	return "UpdateCredDef"
}

func (msg *MsgUpdateCredDef) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCredDef) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCredDef) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateCredDef{}

func NewMsgDeleteCredDef(creator string, id uint64) *MsgDeleteCredDef {
	return &MsgDeleteCredDef{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteCredDef) Route() string {
	return RouterKey
}

func (msg *MsgDeleteCredDef) Type() string {
	return "DeleteCredDef"
}

func (msg *MsgDeleteCredDef) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteCredDef) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteCredDef) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
