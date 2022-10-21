package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ StateValueData = &Did{}

func NewDid(context []string, id string, controller []string, verificationMethod []*VerificationMethod,
	authentication []string, assertionMethod []string, capabilityInvocation []string, capabilityDelegation []string,
	keyAgreement []string, service []*Service, alsoKnownAs []string,
) *Did {
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

// AllControllerDids returns controller DIDs used in both did.controllers and did.verification_method.controller
func (did *Did) AllControllerDids() []string {
	result := did.Controller
	result = append(result, did.GetVerificationMethodControllers()...)

	return utils.UniqueSorted(result)
}

// ReplaceIds replaces ids in all controller and id fields
func (did *Did) ReplaceIds(old, new string) {
	// Controllers
	utils.ReplaceInSlice(did.Controller, old, new)

	// Id
	if did.Id == old {
		did.Id = new
	}

	for _, vm := range did.VerificationMethod {
		vm.ReplaceIds(old, new)
	}
}

func (did *Did) GetControllersOrSubject() []string {
	result := did.Controller

	if len(result) == 0 {
		result = append(result, did.Id)
	}

	return result
}

func (did *Did) GetVerificationMethodControllers() []string {
	var result []string

	for _, vm := range did.VerificationMethod {
		result = append(result, vm.Controller)
	}

	return result
}

// Validation

func (did Did) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&did,
		validation.Field(&did.Id, validation.Required, IsDID(allowedNamespaces)),
		validation.Field(&did.Controller, IsUniqueStrList(), validation.Each(IsDID(allowedNamespaces))),
		validation.Field(&did.VerificationMethod,
			IsUniqueVerificationMethodListByIdRule(), validation.Each(ValidVerificationMethodRule(did.Id, allowedNamespaces)),
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

		validation.Field(&did.Service, IsUniqueServiceListByIdRule(), validation.Each(ValidServiceRule(did.Id, allowedNamespaces))),
		validation.Field(&did.AlsoKnownAs, IsUniqueStrList(), validation.Each(IsURI())),
	)
}
