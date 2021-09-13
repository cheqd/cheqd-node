package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateAttrib{}

func NewMsgCreateAttrib(creator string, did string, raw string) *MsgCreateAttrib {
	return &MsgCreateAttrib{
		Creator: creator,
		Did:     did,
		Raw:     raw,
	}
}

func (msg *MsgCreateAttrib) Route() string {
	return RouterKey
}

func (msg *MsgCreateAttrib) Type() string {
	return "CreateAttrib"
}

func (msg *MsgCreateAttrib) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateAttrib) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateAttrib) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateAttrib{}

func NewMsgUpdateAttrib(creator string, id uint64, did string, raw string) *MsgUpdateAttrib {
	return &MsgUpdateAttrib{
		Id:      id,
		Creator: creator,
		Did:     did,
		Raw:     raw,
	}
}

func (msg *MsgUpdateAttrib) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAttrib) Type() string {
	return "UpdateAttrib"
}

func (msg *MsgUpdateAttrib) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAttrib) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAttrib) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateAttrib{}

func NewMsgDeleteAttrib(creator string, id uint64) *MsgDeleteAttrib {
	return &MsgDeleteAttrib{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteAttrib) Route() string {
	return RouterKey
}

func (msg *MsgDeleteAttrib) Type() string {
	return "DeleteAttrib"
}

func (msg *MsgDeleteAttrib) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteAttrib) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteAttrib) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
