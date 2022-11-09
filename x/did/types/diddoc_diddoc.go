package types

import (
	"github.com/cheqd/cheqd-node/x/did/utils"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func NewDidDoc(context []string, id string, controller []string, verificationMethod []*VerificationMethod,
	authentication []string, assertionMethod []string, capabilityInvocation []string, capabilityDelegation []string,
	keyAgreement []string, service []*Service, alsoKnownAs []string,
) *DidDoc {
	return &DidDoc{
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
func (didDoc *DidDoc) AllControllerDids() []string {
	result := didDoc.Controller
	result = append(result, didDoc.GetVerificationMethodControllers()...)

	return utils.UniqueSorted(result)
}

// ReplaceDids replaces ids in all controller and id fields
func (didDoc *DidDoc) ReplaceDids(old, new string) {
	// Controllers
	utils.ReplaceInSlice(didDoc.Controller, old, new)

	// Id
	if didDoc.Id == old {
		didDoc.Id = new
	}

	for _, vm := range didDoc.VerificationMethod {
		vm.ReplaceDids(old, new)
	}
}

func (didDoc *DidDoc) GetControllersOrSubject() []string {
	result := didDoc.Controller

	if len(result) == 0 {
		result = append(result, didDoc.Id)
	}

	return result
}

func (didDoc *DidDoc) GetVerificationMethodControllers() []string {
	var result []string

	for _, vm := range didDoc.VerificationMethod {
		result = append(result, vm.Controller)
	}

	return result
}

func NewDidDocFromV1(didV1 didtypesv1.Did) *DidDoc {
	vms = []*VerificationMethod{}
	for _, vm := range didV1.VerificationMethod {
		didDoc.VerificationMethod = append(
			vms, 
			VerificationMethod{
				Id: vm.Id,
				Type: vm.Type,
				Controller: vm.Controller,
				PublicKeyBase58: vm.PublicKeyBase58,
			})
	}
	srvs := []*Service{}
	for _, srv := range didV1.Service {
		srvs = append(
			srvs,
			Service{
				Id: srv.Id,
				Type: srv.Type,
				ServiceEndpoint: []string{srv.ServiceEndpoint},
			})
	}
	return &DidDoc{
		Context:              didV1.Context,
		Id:                   didV1.Id,
		Controller:           didV1.Controller,
		VerificationMethod:   vms,
		Authentication:       didV1.Authentication,
		AssertionMethod:      didV1.AssertionMethod,
		CapabilityInvocation: didV1.CapabilityInvocation,
		CapabilityDelegation: didV1.CapabilityDelegation,
		KeyAgreement:         didV1.KeyAgreement,
		Service:              didV1.Service,
		AlsoKnownAs:          didV1.AlsoKnownAs,
	}
}

// Validation

func (didDoc DidDoc) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&didDoc,
		validation.Field(&didDoc.Id, validation.Required, IsDID(allowedNamespaces)),
		validation.Field(&didDoc.Controller, IsUniqueStrList(), validation.Each(IsDID(allowedNamespaces))),
		validation.Field(&didDoc.VerificationMethod,
			IsUniqueVerificationMethodListByIdRule(), validation.Each(ValidVerificationMethodRule(didDoc.Id, allowedNamespaces)),
		),

		validation.Field(&didDoc.Authentication,
			IsUniqueStrList(), validation.Each(IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(didDoc.Id)),
		),
		validation.Field(&didDoc.AssertionMethod,
			IsUniqueStrList(), validation.Each(IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(didDoc.Id)),
		),
		validation.Field(&didDoc.CapabilityInvocation,
			IsUniqueStrList(), validation.Each(IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(didDoc.Id)),
		),
		validation.Field(&didDoc.CapabilityDelegation,
			IsUniqueStrList(), validation.Each(IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(didDoc.Id)),
		),
		validation.Field(&didDoc.KeyAgreement,
			IsUniqueStrList(), validation.Each(IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(didDoc.Id)),
		),

		validation.Field(&didDoc.Service, IsUniqueServiceListByIdRule(), validation.Each(ValidServiceRule(didDoc.Id, allowedNamespaces))),
		validation.Field(&didDoc.AlsoKnownAs, IsUniqueStrList(), validation.Each(IsURI())),
	)
}
