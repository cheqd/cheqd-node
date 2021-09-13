package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateSchema{}

func NewMsgCreateSchema(creator string, name string, version string, attr_names string) *MsgCreateSchema {
	return &MsgCreateSchema{
		Creator:    creator,
		Name:       name,
		Version:    version,
		Attr_names: attr_names,
	}
}

func (msg *MsgCreateSchema) Route() string {
	return RouterKey
}

func (msg *MsgCreateSchema) Type() string {
	return "CreateSchema"
}

func (msg *MsgCreateSchema) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateSchema) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateSchema) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateSchema{}

func NewMsgUpdateSchema(creator string, id uint64, name string, version string, attr_names string) *MsgUpdateSchema {
	return &MsgUpdateSchema{
		Id:         id,
		Creator:    creator,
		Name:       name,
		Version:    version,
		Attr_names: attr_names,
	}
}

func (msg *MsgUpdateSchema) Route() string {
	return RouterKey
}

func (msg *MsgUpdateSchema) Type() string {
	return "UpdateSchema"
}

func (msg *MsgUpdateSchema) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateSchema) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateSchema) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateSchema{}

func NewMsgDeleteSchema(creator string, id uint64) *MsgDeleteSchema {
	return &MsgDeleteSchema{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteSchema) Route() string {
	return RouterKey
}

func (msg *MsgDeleteSchema) Type() string {
	return "DeleteSchema"
}

func (msg *MsgDeleteSchema) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteSchema) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteSchema) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
