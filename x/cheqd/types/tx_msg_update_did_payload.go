package types

import validation "github.com/go-ozzo/ozzo-validation/v4"

var _ IdentityMsg = &MsgUpdateDidPayload{}

func (msg *MsgUpdateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

func (msg *MsgUpdateDidPayload) ToDid() Did {
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
