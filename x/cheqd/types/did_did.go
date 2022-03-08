package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ StateValueData = &Did{}

func NewDid(context []string, id string, controller []string, verificationMethod []*VerificationMethod,
	authentication []string, assertionMethod []string, capabilityInvocation []string, capabilityDelegation []string,
	keyAgreement []string, service []*Service, alsoKnownAs []string) *Did {

	return &Did{
		Context:              context,
		Id:                   id,
		Controller:           controller,
		VerificationMethod:   verificationMethod,
		Authentication:       authentication,
		AssertionMethod:      assertionMethod,
		CapabilityInvocation: capabilityInvocation,
		CapabilityDelegation: capabilityDelegation,
		KeyAgreement:         keyAgreement,
		Service:              service,
		AlsoKnownAs:          alsoKnownAs,
	}
}

// Helpers

// AggregateControllerDids returns controller DIDs used in both did.controllers and did.verification_method.controller
func (did *Did) AggregateControllerDids() []string {
	result := did.Controller

	for _, vm := range did.VerificationMethod {
		result = append(result, vm.Controller)
	}

	return utils.Unique(result)
}

// Validation

func (did Did) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&did,
		validation.Field(&did.Id, validation.Required, IsDID(allowedNamespaces)),
		validation.Field(&did.Controller, IsUniqueStrList(), validation.Each(IsDID(allowedNamespaces))),
		validation.Field(&did.VerificationMethod,
			IsUniqueVerificationMethodList(), validation.Each(ValidVerificationMethod(did.Id, allowedNamespaces)),
		),

		validation.Field(&did.Authentication,
			IsUniqueStrList(), validation.Each(IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(did.Id)),
		),
		validation.Field(&did.AssertionMethod,
			IsUniqueStrList(), validation.Each(IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(did.Id)),
		),
		validation.Field(&did.CapabilityInvocation,
			IsUniqueStrList(), validation.Each(IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(did.Id)),
		),
		validation.Field(&did.CapabilityDelegation,
			IsUniqueStrList(), validation.Each(IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(did.Id)),
		),
		validation.Field(&did.KeyAgreement,
			IsUniqueStrList(), validation.Each(IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(did.Id)),
		),

		validation.Field(&did.Service, IsUniqueServiceList(), validation.Each(ValidService(did.Id, allowedNamespaces))),
		validation.Field(&did.AlsoKnownAs, IsUniqueStrList(), validation.Each(IsDID(allowedNamespaces))),
	)
}
