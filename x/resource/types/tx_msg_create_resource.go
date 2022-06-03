package types

import (
	cheqd_types "github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateResource{}

func NewMsgCreateResource(payload *MsgCreateResourcePayload, signatures []*cheqd_types.SignInfo) *MsgCreateResource {
	return &MsgCreateResource{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgCreateResource) Route() string {
	return RouterKey
}

func (msg *MsgCreateResource) Type() string {
	return "MsgCreateResource"
}

func (msg *MsgCreateResource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateResource) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshal(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateResource) ValidateBasic() error {
	err := msg.Validate()
	if err != nil {
		return ErrBasicValidation.Wrap(err.Error())
	}

	return nil
}

// Validate

func (msg MsgCreateResource) Validate() error {
	//TODO: implement
	return nil
}
