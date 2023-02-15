package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ sdk.Msg = &MsgCreateDidDoc{}

func NewMsgCreateDid(payload *MsgCreateDidDocPayload, signatures []*SignInfo) *MsgCreateDidDoc {
	return &MsgCreateDidDoc{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgCreateDidDoc) Route() string {
	return RouterKey
}

func (msg *MsgCreateDidDoc) Type() string {
	return "MsgCreateDidDoc"
}

func (msg *MsgCreateDidDoc) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateDidDoc) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDidDoc) ValidateBasic() error {
	err := msg.Validate(nil)
	if err != nil {
		return ErrBasicValidation.Wrap(err.Error())
	}

	return nil
}

// Validate

func (msg MsgCreateDidDoc) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.Payload, validation.Required, ValidMsgCreateDidPayloadRule(allowedNamespaces)),
		validation.Field(&msg.Signatures, IsUniqueSignInfoListByIDRule(), validation.Each(ValidSignInfoRule(allowedNamespaces))),
	)
}

// Normalize

func (msg *MsgCreateDidDoc) Normalize() {
	msg.Payload.Normalize()
	NormalizeSignInfoList(msg.Signatures)
}
