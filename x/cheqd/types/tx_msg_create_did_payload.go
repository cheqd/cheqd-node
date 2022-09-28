package types

var _ IdentityMsg = &MsgCreateDidPayload{}

func (msg *MsgCreateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

func (msg *MsgCreateDidPayload) ToDid() Did {
	did := &Did{
		Context:              msg.Context,
		Id:                   msg.Id,
		Controller:           msg.Controller,
		Authentication:       msg.Authentication,
		AssertionMethod:      msg.AssertionMethod,
		CapabilityInvocation: msg.CapabilityInvocation,
		CapabilityDelegation: msg.CapabilityDelegation,
		KeyAgreement:         msg.KeyAgreement,
		AlsoKnownAs:          msg.AlsoKnownAs,
	}
	for _, vm := range msg.VerificationMethod {
		newVM := VerificationMethod(*vm)
		did.VerificationMethod = append(did.VerificationMethod, &newVM)
	}
	for _, s := range msg.Service {
		newS := Service(*s)
		did.Service = append(did.Service, &newS)
	}
	NormalizeDID(did)
	return *did
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
