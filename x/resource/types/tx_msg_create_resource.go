package types

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ sdk.Msg = &MsgCreateResource{}

func NewMsgCreateResource(payload *MsgCreateResourcePayload, signatures []*didtypes.SignInfo) *MsgCreateResource {
	return &MsgCreateResource{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgCreateResource) Type() string {
	return "MsgCreateResource"
}

// Validate

func (msg MsgCreateResource) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.Payload, validation.Required, ValidMsgCreateResourcePayload()),
		validation.Field(&msg.Signatures, didtypes.IsUniqueSignInfoListRule(), validation.Each(didtypes.ValidSignInfoRule(allowedNamespaces))),
	)
}

// Normalize

func (msg *MsgCreateResource) Normalize() {
	msg.Payload.Normalize()
	didtypes.NormalizeSignInfoList(msg.Signatures)
}
