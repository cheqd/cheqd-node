package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgCreateDid(
	id string,
	controller []string,
	verificationMethod []*VerificationMethod,
	authentication []string,
	assertionMethod []string,
	capabilityInvocation []string,
	capabilityDelegation []string,
	keyAgreement []string,
	alsoKnownAs []string,
	service []*DidService,
	context []string,
) *MsgCreateDid {
	return &MsgCreateDid{
		Id:                   id,
		Controller:           controller,
		VerificationMethod:   verificationMethod,
		Authentication:       authentication,
		AssertionMethod:      assertionMethod,
		CapabilityInvocation: capabilityInvocation,
		CapabilityDelegation: capabilityDelegation,
		KeyAgreement:         keyAgreement,
		AlsoKnownAs:          alsoKnownAs,
		Service:              service,
		Context:              context,
	}
}

var _ IdentityMsg = &MsgCreateDid{}

func (msg *MsgCreateDid) GetSigners() []Signer {
	if len(msg.Controller) > 0 {
		result := make([]Signer, len(msg.Controller))

		for i, signer := range msg.Controller {
			if signer == msg.Id {
				result[i] = Signer{
					Signer:             signer,
					Authentication:     msg.Authentication,
					VerificationMethod: msg.VerificationMethod,
				}
			} else {
				result[i] = Signer{
					Signer: signer,
				}
			}
		}

		return result
	}

	if len(msg.Authentication) > 0 {
		return []Signer{
			{
				Signer:             msg.Id,
				Authentication:     msg.Authentication,
				VerificationMethod: msg.VerificationMethod,
			},
		}
	}

	return []Signer{}
}

func (msg *MsgCreateDid) ValidateBasic() error {
	if utils.IsNotDid(msg.Id) {
		return ErrBadRequestIsNotDid.Wrap("Id")
	}

	if notValid, i := utils.ArrayContainsNotDid(msg.Controller); notValid {
		return ErrBadRequestIsNotDid.Wrapf("Controller item %s at position %d", msg.Controller[i], i)
	}

	if err := ValidateVerificationMethods(msg.Id, msg.VerificationMethod); err != nil {
		return err
	}

	if err := ValidateServices(msg.Id, msg.Service); err != nil {
		return err
	}

	if notValid, i := utils.ArrayContainsNotDidFragment(msg.Authentication); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("Authentication item %s", msg.Authentication[i])
	}

	if notValid, i := utils.ArrayContainsNotDidFragment(msg.CapabilityInvocation); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("CapabilityInvocation item %s", msg.CapabilityInvocation[i])
	}

	if notValid, i := utils.ArrayContainsNotDidFragment(msg.CapabilityDelegation); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("CapabilityDelegation item %s", msg.CapabilityDelegation[i])
	}

	if notValid, i := utils.ArrayContainsNotDidFragment(msg.KeyAgreement); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("KeyAgreement item %s", msg.KeyAgreement[i])
	}

	if len(msg.Authentication) == 0 && len(msg.Controller) == 0 {
		return ErrBadRequest.Wrap("The message must contain either a Controller or a Authentication")
	}

	for _, i := range msg.Authentication {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	for _, i := range msg.KeyAgreement {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	for _, i := range msg.CapabilityDelegation {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	for _, i := range msg.CapabilityInvocation {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	return nil
}

var _ IdentityMsg = &MsgUpdateDid{}

func NewMsgUpdateDid(
	id string,
	controller []string,
	verificationMethod []*VerificationMethod,
	authentication []string,
	assertionMethod []string,
	capabilityInvocation []string,
	capabilityDelegation []string,
	keyAgreement []string,
	alsoKnownAs []string,
	service []*DidService,
	context []string,
) *MsgUpdateDid {
	return &MsgUpdateDid{
		Id:                   id,
		Controller:           controller,
		VerificationMethod:   verificationMethod,
		Authentication:       authentication,
		AssertionMethod:      assertionMethod,
		CapabilityInvocation: capabilityInvocation,
		CapabilityDelegation: capabilityDelegation,
		KeyAgreement:         keyAgreement,
		AlsoKnownAs:          alsoKnownAs,
		Service:              service,
		Context:              context,
	}
}

func (msg *MsgUpdateDid) GetSigners() []Signer {
	if len(msg.Controller) > 0 {
		result := make([]Signer, len(msg.Controller))

		for i, signer := range msg.Controller {
			if signer == msg.Id {
				result[i] = Signer{
					Signer:             signer,
					Authentication:     msg.Authentication,
					VerificationMethod: msg.VerificationMethod,
				}
			} else {
				result[i] = Signer{
					Signer: signer,
				}
			}
		}

		return result
	}

	if len(msg.Authentication) > 0 {
		return []Signer{
			{
				Signer:             msg.Id,
				Authentication:     msg.Authentication,
				VerificationMethod: msg.VerificationMethod,
			},
		}
	}

	return []Signer{}
}

func (msg *MsgUpdateDid) ValidateBasic() error {
	if utils.IsNotDid(msg.Id) {
		return ErrBadRequestIsNotDid.Wrap("Id")
	}

	if notValid, i := utils.ArrayContainsNotDid(msg.Controller); notValid {
		return ErrBadRequestIsNotDid.Wrapf("Controller item %s at position %d", msg.Controller[i], i)
	}

	if err := ValidateVerificationMethods(msg.Id, msg.VerificationMethod); err != nil {
		return err
	}

	if err := ValidateServices(msg.Id, msg.Service); err != nil {
		return err
	}

	if notValid, i := utils.ArrayContainsNotDidFragment(msg.Authentication); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("Authentication item %s", msg.Authentication[i])
	}

	if notValid, i := utils.ArrayContainsNotDidFragment(msg.CapabilityInvocation); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("CapabilityInvocation item %s", msg.CapabilityInvocation[i])
	}

	if notValid, i := utils.ArrayContainsNotDidFragment(msg.CapabilityDelegation); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("CapabilityDelegation item %s", msg.CapabilityDelegation[i])
	}

	if notValid, i := utils.ArrayContainsNotDidFragment(msg.KeyAgreement); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("KeyAgreement item %s", msg.KeyAgreement[i])
	}

	if len(msg.Authentication) == 0 && len(msg.Controller) == 0 {
		return ErrBadRequest.Wrap("The message must contain either a Controller or a Authentication")
	}

	for _, i := range msg.Authentication {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	for _, i := range msg.KeyAgreement {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	for _, i := range msg.CapabilityDelegation {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	for _, i := range msg.CapabilityInvocation {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	return nil
}

func ValidateVerificationMethods(did string, vms []*VerificationMethod) error {
	for i, vm := range vms {
		if err := ValidateVerificationMethod(vm); err != nil {
			return ErrBadRequestInvalidVerMethod.Wrap(sdkerrors.Wrapf(err, "index %d, value %s", i, vm.Id).Error())
		}
	}

	for i, vm := range vms {
		if IncludeVerificationMethod(did, vms[i+1:], vm.Id) {
			return ErrBadRequestInvalidVerMethod.Wrapf("%s is duplicated", vm.Id)
		}
	}

	return nil
}

func ValidateVerificationMethod(vm *VerificationMethod) error {
	if !utils.IsFullDidFragment(vm.Id) {
		return ErrBadRequestIsNotDidFragment.Wrap(vm.Id)
	}

	if !utils.IsVerificationMethodType(vm.Type) {
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

func ValidateServices(did string, services []*DidService) error {
	for i, s := range services {
		if err := ValidateService(s); err != nil {
			return ErrBadRequestInvalidService.Wrap(sdkerrors.Wrapf(err, "index %d, value %s", i, s.Id).Error())
		}
	}

	for i, s := range services {
		if IncludeService(did, services[i+1:], s.Id) {
			return ErrBadRequestInvalidService.Wrapf("%s is duplicated", s.Id)
		}
	}

	return nil
}

func ValidateService(s *DidService) error {
	if !utils.IsDidFragment(s.Id) {
		return ErrBadRequestIsNotDidFragment.Wrap(s.Id)
	}

	if !utils.IsDidServiceType(s.Type) {
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

func IncludeService(did string, services []*DidService, id string) bool {
	for _, s := range services {
		if utils.ResolveId(did, s.Id) == utils.ResolveId(did, id) {
			return true
		}
	}

	return false
}
