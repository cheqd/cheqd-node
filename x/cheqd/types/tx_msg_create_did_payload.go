package types

import "github.com/cheqd/cheqd-node/x/cheqd/utils"

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
func (msg *MsgCreateDidPayload) Normalize() {
	msg.Id = utils.NormalizeDID(msg.Id)
	for _, vm := range msg.VerificationMethod {
		vm.Controller = utils.NormalizeDID(vm.Controller)
		vm.Id = utils.NormalizeDIDUrl(vm.Id)
	}
	for _, s := range msg.Service {
		s.Id = utils.NormalizeDIDUrl(s.Id)
	}
	msg.Controller = utils.NormalizeDIDList(msg.Controller)
	msg.Authentication = utils.NormalizeDIDUrlList(msg.Authentication)
	msg.AssertionMethod = utils.NormalizeDIDUrlList(msg.AssertionMethod)
	msg.CapabilityInvocation = utils.NormalizeDIDUrlList(msg.CapabilityInvocation)
	msg.CapabilityDelegation = utils.NormalizeDIDUrlList(msg.CapabilityDelegation)
	msg.KeyAgreement = utils.NormalizeDIDUrlList(msg.KeyAgreement)
}

func (msg *MsgCreateDid) Normalize() {
	msg.Payload.Normalize()
	NormalizeSignatureUUIDIdentifiers(msg.Signatures)
}
