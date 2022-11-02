package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ sdk.Msg = &MsgUpdateDidDoc{}

func NewMsgUpdateDid(payload *MsgUpdateDidDocPayload, signatures []*SignInfo) *MsgUpdateDidDoc {
	return &MsgUpdateDidDoc{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgUpdateDidDoc) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDidDoc) Type() string {
	return "MsgUpdateDidDoc"
}

func (msg *MsgUpdateDidDoc) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgUpdateDidDoc) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDidDoc) ValidateBasic() error {
	err := msg.Validate(nil)
	if err != nil {
		return ErrBasicValidation.Wrap(err.Error())
	}

	return nil
}

// Validate

func (msg MsgUpdateDidDoc) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.Payload, validation.Required, ValidMsgUpdateDidPayloadRule(allowedNamespaces)),
		validation.Field(&msg.Signatures, IsUniqueSignInfoListRule(), validation.Each(ValidSignInfoRule(allowedNamespaces))),
	)
}

// Normalize

func (msg *MsgUpdateDidDoc) Normalize() {
	msg.Payload.Normalize()
	NormalizeSignInfoList(msg.Signatures)
}
