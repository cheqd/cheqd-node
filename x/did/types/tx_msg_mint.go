package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgMint{}

func NewMsgMint(authority string, toAddr string, coins sdk.Coins) *MsgMint {
	return &MsgMint{
		Authority: authority,
		ToAddress: toAddr,
		Amount:    coins,
	}
}

func (msg *MsgMint) Route() string {
	return RouterKey
}

func (msg *MsgMint) Type() string {
	return "MsgMint"
}

func (msg *MsgMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgMint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMint) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	// Check if the 'toAddress' is a valid Bech32 address
	_, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid recipient address")
	}

	// Validate that the 'amount' is a valid coin denomination and positive value
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid coin denomination or amount")
	}

	return nil
}
