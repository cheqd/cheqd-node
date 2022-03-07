package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	JsonWebKey2020             = "JsonWebKey2020"
	Ed25519VerificationKey2020 = "Ed25519VerificationKey2020"
)

var SupportedMethods = []string{
	JsonWebKey2020,
	Ed25519VerificationKey2020,
}

var JwkMethods = []string{
	JsonWebKey2020,
}

var MultibaseMethods = []string{
	Ed25519VerificationKey2020,
}

func NewVerificationMethod(id string, type_ string, controller string, publicKeyJwk []*KeyValuePair, publicKeyMultibase string) *VerificationMethod {
	return &VerificationMethod{
		Id:                 id,
		Type:               type_,
		Controller:         controller,
		PublicKeyJwk:       publicKeyJwk,
		PublicKeyMultibase: publicKeyMultibase,
	}
}

func (vm VerificationMethod) Validate() error {
	return validation.ValidateStruct(&vm,
		validation.Field(&vm.Id, validation.Required, IsDID()),
		validation.Field(&vm.Controller, validation.Required, IsDID()),
		validation.Field(&vm.Type, validation.Required, validation.In(utils.ToInterfaces(SupportedMethods)...)),
		validation.Field(&vm.PublicKeyJwk,
			validation.When(utils.Contains(JwkMethods, vm.Type), validation.Required).Else(validation.Empty),
		),
		validation.Field(&vm.PublicKeyMultibase,
			validation.When(utils.Contains(MultibaseMethods, vm.Type), validation.Required, IsMultibase()).Else(validation.Empty),
		),
	)
}

func FindVerificationMethod(vms []VerificationMethod, id string) (VerificationMethod, bool) {
	for _, vm := range vms {
		if vm.Id == id {
			return vm, true
		}
	}

	return VerificationMethod{}, false
}
