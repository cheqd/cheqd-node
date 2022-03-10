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

// TODO: High priority
// Define dynamic validation logic in ADR
// TODO: Medium priority
// Having multiple namespaces is impossible now. Do we need fix?
// Do we need check references existence? It's not necessary according to spec but was implemented.
// Validate keys in verification methods - Andrew N.
// Rename DID to DIDDoc
// TODO: Low priority
// Migrate old tests for static validation
// Check if signatures are checked in static validation most likely no

func (k msgServer) CreateDid(goCtx context.Context, msg *types.MsgCreateDid) (*types.MsgCreateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate namespaces
	namespace := k.GetDidNamespace(ctx)
	// FIXME: We can't allow empty namespace because currently DIDs are stored as did.id -> did.
	// FIXME: So 'did:cheqd:mainnet:abc' and 'did:cheqd:abc' will be counted as different DIDs.
	err := msg.Validate([]string{namespace})
	if err != nil {
		return nil, types.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Validate DID doesn't exist
	if k.HasDid(&ctx, msg.Payload.Id) {
		return nil, types.ErrDidDocExists.Wrap(msg.Payload.Id)
	}


	// Build metadata and stateValue
	did := msg.Payload.ToDid()
	metadata := types.NewMetadataFromContext(ctx)
	stateValue, err := types.NewStateValue(&did, &metadata)
	if err != nil {
		return nil, err
	}

	// Check there is at least one validation method
	if len(did.VerificationMethod) < 1 {
		return nil, types.ErrBadRequestInvalidVerMethod.Wrap("there should be at least one verification method")
	}

	// Check controllers' existence
	controllers := did.AggregateControllerDids()
	for _, controller := range controllers {
		_, found, err := FindDid(&k.Keeper, &ctx, controller, map[string]types.StateValue{did.Id: stateValue})
		if err != nil {
			return nil, err
		}

		if !found {
			return nil, types.ErrDidDocNotFound.Wrapf("controller not found: %s", controller)
		}
	}

	// Verify signatures
	if err := k.VerifySignature(&ctx, payload, payload.GetSigners(), msg.GetSignatures()); err != nil {
		return nil, err
	}

	// Write to state
	id, err := k.AppendDid(&ctx, &did, &metadata)
	if err != nil {
		return nil, err
	}

	// Build response
	return &types.MsgCreateDidResponse{
		Id: id,
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
