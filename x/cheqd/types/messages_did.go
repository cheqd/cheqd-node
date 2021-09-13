package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateDid{}

func NewMsgCreateDid(creator string, verkey string, alias string) *MsgCreateDid {
	return &MsgCreateDid{
		Creator: creator,
		Verkey:  verkey,
		Alias:   alias,
	}
}

func (msg *MsgCreateDid) Route() string {
	return RouterKey
}

func (msg *MsgCreateDid) Type() string {
	return "CreateDid"
}

func (msg *MsgCreateDid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDid) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateDid{}

func NewMsgUpdateDid(creator string, id uint64, verkey string, alias string) *MsgUpdateDid {
	return &MsgUpdateDid{
		Id:      id,
		Creator: creator,
		Verkey:  verkey,
		Alias:   alias,
	}
}

func (msg *MsgUpdateDid) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDid) Type() string {
	return "UpdateDid"
}

func (msg *MsgUpdateDid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDid) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateDid{}

func NewMsgDeleteDid(creator string, id uint64) *MsgDeleteDid {
	return &MsgDeleteDid{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteDid) Route() string {
	return RouterKey
}

func (msg *MsgDeleteDid) Type() string {
	return "DeleteDid"
}

func (msg *MsgDeleteDid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteDid) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
