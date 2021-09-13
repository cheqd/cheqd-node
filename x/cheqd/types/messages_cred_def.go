package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateCred_def{}

func NewMsgCreateCred_def(creator string, schema_id string, tag string, signature_type string, value string) *MsgCreateCred_def {
	return &MsgCreateCred_def{
		Creator:        creator,
		Schema_id:      schema_id,
		Tag:            tag,
		Signature_type: signature_type,
		Value:          value,
	}
}

func (msg *MsgCreateCred_def) Route() string {
	return RouterKey
}

func (msg *MsgCreateCred_def) Type() string {
	return "CreateCred_def"
}

func (msg *MsgCreateCred_def) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateCred_def) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateCred_def) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateCred_def{}

func NewMsgUpdateCred_def(creator string, id uint64, schema_id string, tag string, signature_type string, value string) *MsgUpdateCred_def {
	return &MsgUpdateCred_def{
		Id:             id,
		Creator:        creator,
		Schema_id:      schema_id,
		Tag:            tag,
		Signature_type: signature_type,
		Value:          value,
	}
}

func (msg *MsgUpdateCred_def) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCred_def) Type() string {
	return "UpdateCred_def"
}

func (msg *MsgUpdateCred_def) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCred_def) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCred_def) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateCred_def{}

func NewMsgDeleteCred_def(creator string, id uint64) *MsgDeleteCred_def {
	return &MsgDeleteCred_def{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteCred_def) Route() string {
	return RouterKey
}

func (msg *MsgDeleteCred_def) Type() string {
	return "DeleteCred_def"
}

func (msg *MsgDeleteCred_def) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteCred_def) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteCred_def) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
