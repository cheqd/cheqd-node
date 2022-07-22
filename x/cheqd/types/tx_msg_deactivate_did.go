package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ sdk.Msg = &MsgDeactivateDid{}

func NewMsgDeactivateDid(payload *MsgDeactivateDidPayload, signatures []*SignInfo) *MsgDeactivateDid {
	return &MsgDeactivateDid{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgDeactivateDid) Route() string {
	return RouterKey
}

func (msg *MsgDeactivateDid) Type() string {
	return "WriteRequest"
}

func (msg *MsgDeactivateDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgDeactivateDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshal(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeactivateDid) ValidateBasic() error {
	err := msg.Validate(nil)
	if err != nil {
		return ErrBasicValidation.Wrap(err.Error())
	}

	return nil
}

// Validate

func (msg MsgDeactivateDid) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.Payload, validation.Required, ValidMsgDeactivateDidPayloadRule(allowedNamespaces)),
		validation.Field(&msg.Signatures, IsUniqueSignInfoListRule(), validation.Each(ValidSignInfoRule(allowedNamespaces))),
	)
}
