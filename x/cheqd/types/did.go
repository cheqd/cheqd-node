package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ StateValueData = &Did{}

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
		validation.Field(&vm.Id, validation.Required, NewDIDRule("", nil)),
		validation.Field(&vm.PublicKeyJwk, validation.When(vm.Type == "jwk", validation.Required.Error("must be set when type is jwk"))),
	)
}

// AggregateControllerDids returns controller DIDs used in both did.controllers and did.verification_method.controller
func (did *Did) AggregateControllerDids() []string {
	result := did.Controller

	for _, vm := range did.VerificationMethod {
		result = append(result, vm.Controller)
	}

	return utils.Unique(result)
}

//func (did *Did) FindVerificationMethod(id string) (VerificationMethod, bool) {
//	for _, vm := range vms {
//		if vm.Id == id {
//			return vm
//		}
//	}
//
//	return nil, true
//}

func FindVerificationMethod(vms []VerificationMethod, id string) (VerificationMethod, bool) {
	for _, vm := range vms {
		if vm.Id ==id {
			return vm, true
		}
	}

	return VerificationMethod{}, false
}

//func FilterVerificationMethods(vms []VerificationMethod, func()) {
//
//}