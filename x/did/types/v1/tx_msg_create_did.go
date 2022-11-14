package v1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ sdk.Msg = &MsgCreateDid{}

func NewMsgCreateDid(payload *MsgCreateDidPayload, signatures []*SignInfo) *MsgCreateDid {
	return &MsgCreateDid{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgCreateDid) Route() string {
	return RouterKey
}

func (msg *MsgCreateDid) Type() string {
	return "MsgCreateDid"
}

func (msg *MsgCreateDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshal(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDid) ValidateBasic() error {
	err := msg.Validate(nil)
	if err != nil {
		return err
	}

	return nil
}

// Validate

func (msg MsgCreateDid) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.Payload, validation.Required, ValidMsgCreateDidPayloadRule(allowedNamespaces)),
		validation.Field(&msg.Signatures, IsUniqueSignInfoListByIdRule(), validation.Each(ValidSignInfoRule(allowedNamespaces))),
	)
}
