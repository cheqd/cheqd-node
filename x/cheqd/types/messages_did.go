package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

var _ sdk.Msg = &MsgCreateDid{}

func NewMsgCreateDid(payload *MsgCreateDidPayload, signatures []*SignInfo) *MsgCreateDid {
	return &MsgCreateDid{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgCreateDid) Route() string {
	return RouterKey
}

func (msg *MsgCreateDid) Type() string {
	return "MsgCreateDid"
}

func (msg *MsgCreateDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshal(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDid) ValidateBasic() error {
	if msg.Payload == nil {
		return ErrBadRequestIsRequired.Wrap("Payload")
	}

	if len(msg.Signatures) == 0 {
		return ErrBadRequestAreRequired.Wrap("Signatures")
	}

	return nil
}

var _ sdk.Msg = &MsgUpdateDid{}

func NewMsgUpdateDid(payload *MsgUpdateDidPayload, signatures []*SignInfo) *MsgUpdateDid {
	return &MsgUpdateDid{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg *MsgUpdateDid) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDid) Type() string {
	return "WriteRequest"
}

func (msg *MsgUpdateDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgUpdateDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshal(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDid) ValidateBasic() error {
	if msg.Payload == nil {
		return ErrBadRequestIsRequired.Wrap("Payload")
	}

	if len(msg.Signatures) == 0 {
		return ErrBadRequestAreRequired.Wrap("Signatures")
	}

	return nil
}

func NewMsgCreateDidPayloadPayload(
	context []string,
	id string,
	controller []string,
	verificationMethod []*VerificationMethod,
	authentication []string,
	assertionMethod []string,
	capabilityInvocation []string,
	capabilityDelegation []string,
	keyAgreement []string,
	alsoKnownAs []string,
	service []*Service,
) *MsgCreateDidPayload {
	return &MsgCreateDidPayload{
		Context:              context,
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
	}
}

var _ IdentityMsg = &MsgCreateDidPayload{}

func (msg *MsgCreateDidPayload) GetSigners() []Signer {
	if len(msg.Controller) > 0 {
		result := make([]Signer, len(msg.Controller))

		for i, controller := range msg.Controller {
			if controller == msg.Id {
				result[i] = Signer{
					Signer:             controller,
					Authentication:     msg.Authentication,
					VerificationMethod: msg.VerificationMethod,
				}
			} else {
				result[i] = Signer{
					Signer: controller,
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

func (msg *MsgCreateDidPayload) Validate(namespace string) error {
	if !utils.IsValidDid(namespace, msg.Id) {
		return ErrBadRequestIsNotDid.Wrap("Id")
	}

	if notValid, i := utils.IsNotValidDIDArray(namespace, msg.Controller); notValid {
		return ErrBadRequestIsNotDid.Wrapf("Controller item %s at position %d", msg.Controller[i], i)
	}

	if err := ValidateVerificationMethods(namespace, msg.Id, msg.VerificationMethod); err != nil {
		return err
	}

	if err := ValidateServices(namespace, msg.Id, msg.Service); err != nil {
		return err
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.Authentication); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("Authentication item %s", msg.Authentication[i])
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.CapabilityInvocation); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("CapabilityInvocation item %s", msg.CapabilityInvocation[i])
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.CapabilityDelegation); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("CapabilityDelegation item %s", msg.CapabilityDelegation[i])
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.KeyAgreement); notValid {
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

func (msg *MsgCreateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

var _ IdentityMsg = &MsgUpdateDidPayload{}

func NewMsgUpdateDidPayloadPayload(
	context []string,
	id string,
	controller []string,
	verificationMethod []*VerificationMethod,
	authentication []string,
	assertionMethod []string,
	capabilityInvocation []string,
	capabilityDelegation []string,
	keyAgreement []string,
	alsoKnownAs []string,
	service []*Service,
) *MsgUpdateDidPayload {
	return &MsgUpdateDidPayload{
		Context:              context,
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
	}
}

func (msg *MsgUpdateDidPayload) GetSigners() []Signer {
	if len(msg.Controller) > 0 {
		result := make([]Signer, len(msg.Controller))

		for i, controller := range msg.Controller {
			if controller == msg.Id {
				result[i] = Signer{
					Signer:             controller,
					Authentication:     msg.Authentication,
					VerificationMethod: msg.VerificationMethod,
				}
			} else {
				result[i] = Signer{
					Signer: controller,
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

func (msg *MsgUpdateDidPayload) Validate(namespace string) error {
	if !utils.IsValidDid(namespace, msg.Id) {
		return ErrBadRequestIsNotDid.Wrap("Id")
	}

	if notValid, i := utils.IsNotValidDIDArray(namespace, msg.Controller); notValid {
		return ErrBadRequestIsNotDid.Wrapf("Controller item %s at position %d", msg.Controller[i], i)
	}

	if err := ValidateVerificationMethods(namespace, msg.Id, msg.VerificationMethod); err != nil {
		return err
	}

	if err := ValidateServices(namespace, msg.Id, msg.Service); err != nil {
		return err
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.Authentication); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("Authentication item %s", msg.Authentication[i])
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.CapabilityInvocation); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("CapabilityInvocation item %s", msg.CapabilityInvocation[i])
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.CapabilityDelegation); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("CapabilityDelegation item %s", msg.CapabilityDelegation[i])
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.KeyAgreement); notValid {
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

func (msg *MsgUpdateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

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
