package types

var _ IdentityMsg = &MsgCreateDidPayload{}

func (msg *MsgCreateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

func (msg *MsgCreateDidPayload) ToDid() Did {
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

func (msg *MsgCreateDidPayload) ToMsg(did *Did) *MsgCreateDidPayload {
	return &MsgCreateDidPayload{
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

func (msg MsgCreateDidPayload) Validate(allowedNamespaces []string) error {
	return msg.ToDid().Validate(allowedNamespaces)
}

func ValidMsgCreateDidPayloadRule(allowedNamespaces []string) *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(*MsgCreateDidPayload)
		if !ok {
			panic("ValidMsgCreateDidPayloadRule must be only applied on MsgCreateDidPayload properties")
		}

		return casted.Validate(allowedNamespaces)
	})
}

// Normalize
func (msg MsgCreateDidPayload) Normalize() *MsgCreateDidPayload {
	did := msg.ToDid()
	normilizedDid := NormalizeDID(&did)

	return msg.ToMsg(normilizedDid)
}

func (msg MsgCreateDid) Normalize() *MsgCreateDid {
	NormalizeSignatureUUIDIdentifiers(msg.Signatures)
	return &MsgCreateDid{
		Payload:    msg.Payload.Normalize(),
		Signatures: msg.Signatures,
	}
}
