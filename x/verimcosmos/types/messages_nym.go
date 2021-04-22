package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateNym{}

func NewMsgCreateNym(creator string, alais string, verkey string, did string, role string) *MsgCreateNym {
	return &MsgCreateNym{
		Creator: creator,
		Alais:   alais,
		Verkey:  verkey,
		Did:     did,
		Role:    role,
	}
}

func (msg *MsgCreateNym) Route() string {
	return RouterKey
}

func (msg *MsgCreateNym) Type() string {
	return "CreateNym"
}

func (msg *MsgCreateNym) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateNym) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateNym) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateNym{}

func NewMsgUpdateNym(creator string, id uint64, alais string, verkey string, did string, role string) *MsgUpdateNym {
	return &MsgUpdateNym{
		Id:      id,
		Creator: creator,
		Alais:   alais,
		Verkey:  verkey,
		Did:     did,
		Role:    role,
	}
}

func (msg *MsgUpdateNym) Route() string {
	return RouterKey
}

func (msg *MsgUpdateNym) Type() string {
	return "UpdateNym"
}

func (msg *MsgUpdateNym) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateNym) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateNym) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateNym{}

func NewMsgDeleteNym(creator string, id uint64) *MsgDeleteNym {
	return &MsgDeleteNym{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteNym) Route() string {
	return RouterKey
}

func (msg *MsgDeleteNym) Type() string {
	return "DeleteNym"
}

func (msg *MsgDeleteNym) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteNym) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteNym) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
