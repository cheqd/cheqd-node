package keeper

import (
	"context"
	"fmt"
	"reflect"

	"github.com/cheqd/cheqd-node/x/cheqd/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDid(goCtx context.Context, msg *types.MsgCreateDid) (*types.MsgCreateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	prefix := k.GetDidPrefix(ctx)

	didMsg := msg.GetPayload()
	if err := didMsg.Validate(prefix); err != nil {
		return nil, err
	}

	if err := k.ValidateDidControllers(&ctx, didMsg.Id, didMsg.Controller, didMsg.VerificationMethod); err != nil {
		return nil, err
	}

	if err := k.VerifySignature(&ctx, didMsg, didMsg.GetSigners(), msg.GetSignatures()); err != nil {
		return nil, err
	}

	// Checks that the did doesn't exist
	if k.HasDid(ctx, didMsg.Id) {
		return nil, sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID is already used by DIDDoc %s", didMsg.Id))
	}

	var did = types.Did{
		Context:              didMsg.Context,
		Id:                   didMsg.Id,
		Controller:           didMsg.Controller,
		VerificationMethod:   didMsg.VerificationMethod,
		Authentication:       didMsg.Authentication,
		AssertionMethod:      didMsg.AssertionMethod,
		CapabilityInvocation: didMsg.CapabilityInvocation,
		CapabilityDelegation: didMsg.CapabilityDelegation,
		KeyAgreement:         didMsg.KeyAgreement,
		AlsoKnownAs:          didMsg.AlsoKnownAs,
		Service:              didMsg.Service,
	}

	metadata := types.NewMetadata(ctx)
	id, err := k.AppendDid(ctx, did, &metadata)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateDidResponse{
		Id: *id,
	}, nil
}

func (k msgServer) UpdateDid(goCtx context.Context, msg *types.MsgUpdateDid) (*types.MsgUpdateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	prefix := k.GetDidPrefix(ctx)

	didMsg := msg.GetPayload()
	if err := didMsg.Validate(prefix); err != nil {
		return nil, err
	}

	// Checks that the did doesn't exist
	if !k.HasDid(ctx, didMsg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", didMsg.Id))
	}

	oldStateValue, err := k.GetDid(&ctx, didMsg.Id)
	if err != nil {
		return nil, err
	}

	oldDIDDoc, err := oldStateValue.GetDid()
	if err != nil {
		return nil, err
	}

	if err := k.ValidateDidControllers(&ctx, didMsg.Id, didMsg.Controller, didMsg.VerificationMethod); err != nil {
		return nil, err
	}

	if err := k.VerifySignatureOnDidUpdate(&ctx, oldDIDDoc, didMsg, msg.Signatures); err != nil {
		return nil, err
	}

	// replay protection
	if oldStateValue.Metadata.VersionId != didMsg.VersionId {
		errMsg := fmt.Sprintf("Expected %s with version %s. Got version %s", didMsg.Id, oldStateValue.Metadata.VersionId, didMsg.VersionId)
		return nil, sdkerrors.Wrap(types.ErrUnexpectedDidVersion, errMsg)
	}

	var did = types.Did{
		Context:              didMsg.Context,
		Id:                   didMsg.Id,
		Controller:           didMsg.Controller,
		VerificationMethod:   didMsg.VerificationMethod,
		Authentication:       didMsg.Authentication,
		AssertionMethod:      didMsg.AssertionMethod,
		CapabilityInvocation: didMsg.CapabilityInvocation,
		CapabilityDelegation: didMsg.CapabilityDelegation,
		KeyAgreement:         didMsg.KeyAgreement,
		AlsoKnownAs:          didMsg.AlsoKnownAs,
		Service:              didMsg.Service,
	}

	metadata := types.NewMetadata(ctx)
	metadata.Created = oldStateValue.Metadata.Created
	metadata.Deactivated = oldStateValue.Metadata.Deactivated

	if err = k.SetDid(ctx, did, &metadata); err != nil {
		return nil, err
	}

	return &types.MsgUpdateDidResponse{
		Id: didMsg.Id,
	}, nil
}

func (k msgServer) VerifySignatureOnDidUpdate(ctx *sdk.Context, oldDIDDoc *types.Did, newDIDDoc *types.MsgUpdateDidPayload, signatures []*types.SignInfo) error {
	var signers []types.Signer

	// Get Old DID Doc controller if it's nil then assign self
	oldController := oldDIDDoc.Controller
	if len(oldController) == 0 {
		oldController = []string{oldDIDDoc.Id}
	}

	for _, controller := range oldController {
		signers = append(signers, types.Signer{Signer: controller})
	}

	for _, oldVM := range oldDIDDoc.VerificationMethod {
		newVM := FindVerificationMethod(newDIDDoc.VerificationMethod, oldVM.Id)

		// Verification Method has been deleted
		if newVM == nil {
			signers = AppendSignerIfNeed(signers, oldVM.Controller, newDIDDoc)
			continue
		}

		// Verification Method has been changed
		if !reflect.DeepEqual(oldVM, newVM) {
			signers = AppendSignerIfNeed(signers, newVM.Controller, newDIDDoc)
		}

		// Verification Method Controller has been changed, need to add old controller
		if newVM.Controller != oldVM.Controller {
			signers = AppendSignerIfNeed(signers, oldVM.Controller, newDIDDoc)
		}
	}

	if err := k.VerifySignature(ctx, newDIDDoc, signers, signatures); err != nil {
		return err
	}

	return nil
}

func AppendSignerIfNeed(signers []types.Signer, controller string, msg *types.MsgUpdateDidPayload) []types.Signer {
	for _, signer := range signers {
		if signer.Signer == controller {
			return signers
		}
	}

	signer := types.Signer{
		Signer: controller,
	}

	if controller == msg.Id {
		signer.VerificationMethod = msg.VerificationMethod
		signer.Authentication = msg.Authentication
	}

	return append(signers, signer)
}

func (k msgServer) ValidateDidControllers(ctx *sdk.Context, id string, controllers []string, verMethods []*types.VerificationMethod) error {

	for _, verificationMethod := range verMethods {
		if err := k.ValidateController(ctx, id, verificationMethod.Controller); err != nil {
			return err
		}
	}

	for _, didController := range controllers {
		if err := k.ValidateController(ctx, id, didController); err != nil {
			return err
		}
	}
	return nil
}
