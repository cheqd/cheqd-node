package types

import validation "github.com/go-ozzo/ozzo-validation/v4"

var _ IdentityMsg = &MsgUpdateDidPayload{}

func (msg *MsgUpdateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

func (msg *MsgUpdateDidPayload) ToDid() Did {
	return Did{
		Context:              msg.Context,
		Id:                   msg.Id,
		Controller:           msg.Controller,
		VerificationMethod:   msg.VerificationMethod,
		Authentication:       msg.Authentication,
		AssertionMethod:      msg.AssertionMethod,
		CapabilityInvocation: msg.CapabilityInvocation,
		CapabilityDelegation: msg.CapabilityDelegation,
		KeyAgreement:         msg.KeyAgreement,
		AlsoKnownAs:          msg.AlsoKnownAs,
		Service:              msg.Service,
	}
}

// Validation

func (msg MsgUpdateDidPayload) Validate() error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.VersionId, validation.Required),
	)
}
