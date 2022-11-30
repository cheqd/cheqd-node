package types

import (
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	"github.com/cheqd/cheqd-node/x/did/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ IdentityMsg = &MsgCreateDidDocPayload{}

func (msg *MsgCreateDidDocPayload) GetSignBytes() []byte {
	bytes, err := msg.Marshal()
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgCreateDidDocPayload) ToDidDoc() DidDoc {
	return DidDoc{
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

func (msg MsgCreateDidDocPayload) Validate(allowedNamespaces []string) error {
	err := msg.ToDidDoc().Validate(allowedNamespaces)
	if err != nil {
		return err
	}

	return validation.ValidateStruct(&msg,
		validation.Field(&msg.VersionId, validation.Required),
	)
}

func ValidMsgCreateDidPayloadRule(allowedNamespaces []string) *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(*MsgCreateDidDocPayload)
		if !ok {
			panic("ValidMsgCreateDidPayloadRule must be only applied on MsgCreateDidPayload properties")
		}

		return casted.Validate(allowedNamespaces)
	})
}

// Normalize

func (msg *MsgCreateDidDocPayload) Normalize() {
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

	msg.VersionId = utils.NormalizeUUID(msg.VersionId)
}

// Creator

func MsgCreateDidPayloadFromDid(did *didtypesv1.Did) didtypesv1.MsgCreateDidPayload {
	return didtypesv1.MsgCreateDidPayload{
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
