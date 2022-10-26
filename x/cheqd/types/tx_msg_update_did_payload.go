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

func (msg *MsgUpdateDidPayload) ToMsg(did *Did) *MsgUpdateDidPayload {
	return &MsgUpdateDidPayload{
		Context:              did.Context,
		Id:                   did.Id,
		Controller:           did.Controller,
		VerificationMethod:   did.VerificationMethod,
		Authentication:       did.Authentication,
		AssertionMethod:      did.AssertionMethod,
		CapabilityInvocation: did.CapabilityInvocation,
		CapabilityDelegation: did.CapabilityDelegation,
		KeyAgreement:         did.KeyAgreement,
		AlsoKnownAs:          did.AlsoKnownAs,
		Service:              did.Service,
	}
}

// Validation

func (msg MsgUpdateDidPayload) Validate(allowedNamespaces []string) error {
	err := msg.ToDid().Validate(allowedNamespaces)
	if err != nil {
		return err
	}

	return validation.ValidateStruct(&msg,
		validation.Field(&msg.VersionId, validation.Required),
	)
}

func ValidMsgUpdateDidPayloadRule(allowedNamespaces []string) *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(*MsgUpdateDidPayload)
		if !ok {
			panic("ValidMsgUpdateDidPayloadRule must be only applied on MsgUpdateDidPayload properties")
		}

		return casted.Validate(allowedNamespaces)
	})
}

// Normalize
func (msg MsgUpdateDidPayload) Normalize() *MsgUpdateDidPayload {
	did := msg.ToDid()
	normilizedDid := NormalizeDID(&did)
	normalized := msg.ToMsg(normilizedDid)
	normalized.VersionId = msg.VersionId
	return normalized
}

func (msg MsgUpdateDid) Normalize() *MsgUpdateDid {
	NormalizeSignatureUUIDIdentifiers(msg.Signatures)
	return &MsgUpdateDid{
		Payload:    msg.Payload.Normalize(),
		Signatures: msg.Signatures,
	}
}
