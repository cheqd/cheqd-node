package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ sdk.Msg = &MsgDeactivateDidDoc{}

func NewMsgDeactivateDid(payload *MsgDeactivateDidDocPayload, signatures []*SignInfo) *MsgDeactivateDidDoc {
	return &MsgDeactivateDidDoc{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgDeactivateDidDoc) ValidateBasic() error {
	err := msg.Validate(nil)
	if err != nil {
		return ErrBasicValidation.Wrap(err.Error())
	}

	return nil
}

// Validate

func (msg MsgDeactivateDidDoc) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.Payload, validation.Required, ValidMsgDeactivateDidPayloadRule(allowedNamespaces)),
		validation.Field(&msg.Signatures, IsUniqueSignInfoListRule(), validation.Each(ValidSignInfoRule(allowedNamespaces))),
	)
}

// Normalize

func (msg *MsgDeactivateDidDoc) Normalize() {
	msg.Payload.Normalize()
	NormalizeSignInfoList(msg.Signatures)
}
