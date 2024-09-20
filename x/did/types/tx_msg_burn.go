package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgBurn{}

// NewMsgBurn creates a message to burn tokens
func NewMsgBurn(sender string, amount sdk.Coins) *MsgBurn {
	return &MsgBurn{
		FromAddress: sender,
		Amount:      amount,
	}
}
func (m MsgBurn) Route() string { return RouterKey }

func (m MsgBurn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !m.Amount.IsValid() || m.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

func (m MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgBurn) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.FromAddress)
	return []sdk.AccAddress{sender}
}
