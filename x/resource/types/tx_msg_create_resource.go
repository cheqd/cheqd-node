package types

import (
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ sdk.Msg = &MsgCreateResource{}

func NewMsgCreateResource(payload *MsgCreateResourcePayload, signatures []*cheqdtypes.SignInfo) *MsgCreateResource {
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
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateResource) ValidateBasic() error {
	err := msg.Validate([]string{})
	if err != nil {
		return ErrBasicValidation.Wrap(err.Error())
	}

	return nil
}

// Validate

func (msg MsgCreateResource) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.Payload, validation.Required, ValidMsgCreateResourcePayload()),
		validation.Field(&msg.Signatures, cheqdtypes.IsUniqueSignInfoListRule(), validation.Each(cheqdtypes.ValidSignInfoRule(allowedNamespaces))),
	)
}

// Normalize

func (msg *MsgCreateResource) Normalize() {
	msg.Payload.Normalize()
	cheqdtypes.NormalizeSignInfoList(msg.Signatures)
}
