package validate

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

func ValidateVerificationMethods(namespace string, did string, vms []*VerificationMethod) error {
	for i, vm := range vms {
		if err := ValidateVerificationMethod(namespace, vm); err != nil {
			return ErrBadRequestInvalidVerMethod.Wrap(sdkerrors.Wrapf(err, "index %d, value %s", i, vm.Id).Error())
		}
	}

	for i, vm := range vms {
		if !strings.HasPrefix(vm.Id, did) {
			return ErrBadRequestInvalidVerMethod.Wrapf("%s not belong %s DID Doc", vm.Id, did)
		}

		if IncludeVerificationMethod(did, vms[i+1:], vm.Id) {
			return ErrBadRequestInvalidVerMethod.Wrapf("%s is duplicated", vm.Id)
		}
	}

	return nil
}

func ValidateVerificationMethod(namespace string, vm *VerificationMethod) error {
	if !utils.IsFullDidFragment(namespace, vm.Id) {
		return ErrBadRequestIsNotDidFragment.Wrap(vm.Id)
	}

	if len(vm.PublicKeyMultibase) != 0 && len(vm.PublicKeyJwk) != 0 {
		return ErrBadRequest.Wrap("contains multiple verification material properties")
	}

	switch utils.GetVerificationMethodType(vm.Type) {
	case utils.PublicKeyJwk:
		if len(vm.PublicKeyJwk) == 0 {
			return ErrBadRequest.Wrapf("%s: should contain `PublicKeyJwk` verification material property", vm.Type)
		}
	case utils.PublicKeyMultibase:
		if len(vm.PublicKeyMultibase) == 0 {
			return ErrBadRequest.Wrapf("%s: should contain `PublicKeyMultibase` verification material property", vm.Type)
		}
	default:
		return ErrBadRequest.Wrapf("%s: unsupported verification method type", vm.Type)
	}

	if len(vm.PublicKeyMultibase) == 0 && vm.PublicKeyJwk == nil {
		return ErrBadRequest.Wrap("The verification method must contain either a PublicKeyMultibase or a PublicKeyJwk")
	}

	if len(vm.Controller) == 0 {
		return ErrBadRequestIsRequired.Wrap("Controller")
	}

	return nil
}

func ValidateServices(namespace string, did string, services []*Service) error {
	for i, s := range services {
		if err := ValidateService(namespace, s); err != nil {
			return ErrBadRequestInvalidService.Wrap(sdkerrors.Wrapf(err, "index %d, value %s", i, s.Id).Error())
		}
	}

	for i, s := range services {
		if !strings.HasPrefix(utils.ResolveId(did, s.Id), did) {
			return ErrBadRequestInvalidService.Wrapf("%s not belong %s DID Doc", s.Id, did)
		}

		if IncludeService(did, services[i+1:], s.Id) {
			return ErrBadRequestInvalidService.Wrapf("%s is duplicated", s.Id)
		}
	}

	return nil
}

func ValidateService(namespace string, s *Service) error {
	if !utils.IsDidFragment(namespace, s.Id) {
		return ErrBadRequestIsNotDidFragment.Wrap(s.Id)
	}

	if !utils.IsValidDidServiceType(s.Type) {
		return ErrBadRequest.Wrapf("%s: unsupported service type", s.Type)
	}

	return nil
}

func IncludeVerificationMethod(did string, vms []*VerificationMethod, id string) bool {
	for _, vm := range vms {
		if vm.Id == utils.ResolveId(did, id) {
			return true
		}
	}

	return false
}

func IncludeService(did string, services []*Service, id string) bool {
	for _, s := range services {
		if utils.ResolveId(did, s.Id) == utils.ResolveId(did, id) {
			return true
		}
	}

	return false
}
