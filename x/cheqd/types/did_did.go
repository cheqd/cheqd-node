package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
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

// AggregateControllerDids returns controller DIDs used in both did.controllers and did.verification_method.controller
func (did *Did) AggregateControllerDids() []string {
	result := did.Controller

	for _, vm := range did.VerificationMethod {
		result = append(result, vm.Controller)
	}

	return utils.Unique(result)
}
