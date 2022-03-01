package keeper

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/cheqd/cheqd-node/x/cheqd/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDid(goCtx context.Context, msg *types.MsgCreateDid) (*types.MsgCreateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation + DID namespaces
	namespace := k.GetDidNamespace(ctx)
	validator, err := types.BuildValidator(types.DidMethod, []string{namespace, ""})
	if err != nil {
		return nil, types.ErrValidatorInitialisation.Wrap(err.Error())
	}

	err = validator.Struct(msg)
	if err != nil {
		return nil, types.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Static
	// Verification methods -> no duplicates
	// Verification methods -> valid types, corresponding properties are set
	// Verification relationships -> no duplicates
	// Verification relationships -> valid references
	// Services -> valid types

	// Dynamic

	// Validate DID doesn't exist
	if k.HasDid(ctx, msg.Payload.Id) {
		return nil, types.ErrDidDocExists.Wrap(msg.Payload.Id)
	}

	// Verify that all controllers have at least one authentication key of supported type
	controllers := msg.Payload.AggregateControllerDids()
	for _, v := range controllers {
		// ValidateDIDHasAtLeaseOneSupportedAuthKey
	}

	//if err := k.ValidateDidControllers(&ctx, payload.Id, payload.Controller, payload.VerificationMethod); err != nil {
	//	return nil, err
	//}


	// Verify signatures
	if err := k.VerifySignature(&ctx, payload, payload.GetSigners(), msg.GetSignatures()); err != nil {
		return nil, err
	}

	// Build DID and metadata
	did := msg.Payload.ToDid()
	metadata := types.NewMetadataFromContext(ctx)

	// Write to state
	id, err := k.AppendDid(ctx, did, &metadata)
	if err != nil {
		return nil, err
	}

	// Build response
	return &types.MsgCreateDidResponse{
		Id: *id,
	}, nil
}

func (msg *MsgCreateDidPayload) ValidateDynamic(namespace string) error {
	if err := ValidateVerificationMethods(namespace, msg.Id, msg.VerificationMethod); err != nil {
		return err
	}

	if err := ValidateServices(namespace, msg.Id, msg.Service); err != nil {
		return err
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


func ValidateVerificationMethods(namespace string, did string, vms []*VerificationMethod) error {
	for i, vm := range vms {
		if err := ValidateVerificationMethod(namespace, vm); err != nil {
			return ErrBadRequestInvalidVerMethod.Wrap(sdkerrors.Wrapf(err, "index %d, value %s", i, vm.Id).Error())
		}
	}
}

func ValidateVerificationMethod(namespace string, vm *VerificationMethod) error {
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
