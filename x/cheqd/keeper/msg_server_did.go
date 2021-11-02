package keeper

import (
	"context"
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
	"github.com/cheqd/cheqd-node/x/cheqd/utils/strings"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDid(goCtx context.Context, msg *v1.MsgCreateDid) (*v1.MsgCreateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	prefix := k.GetDidPrefix(ctx)

	didMsg := msg.GetPayload()
	if err := didMsg.Validate(prefix); err != nil {
		return nil, err
	}

	if err := k.VerifySignature(&ctx, didMsg, didMsg.GetSigners(), msg.GetSignatures()); err != nil {
		return nil, err
	}

	// Checks that the did doesn't exist
	if err := k.EnsureDidIsNotUsed(ctx, didMsg.Id); err != nil {
		return nil, err
	}

	var did = v1.Did{
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
		Context:              didMsg.Context,
	}

	metadata := v1.NewMetadata(ctx)
	id, err := k.AppendDid(ctx, did, &metadata)
	if err != nil {
		return nil, err
	}

	return &v1.MsgCreateDidResponse{
		Id: *id,
	}, nil
}

func (k msgServer) UpdateDid(goCtx context.Context, msg *v1.MsgUpdateDid) (*v1.MsgUpdateDidResponse, error) {
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

	if err := k.VerifySignatureOnDidUpdate(&ctx, oldDIDDoc, didMsg, msg.Signatures); err != nil {
		return nil, err
	}

	// replay protection
	if oldStateValue.Metadata.VersionId != didMsg.VersionId {
		errMsg := fmt.Sprintf("Ecpected %s with version %s. Got version %s", didMsg.Id, oldStateValue.Metadata.VersionId, didMsg.VersionId)
		return nil, sdkerrors.Wrap(v1.ErrUnexpectedDidVersion, errMsg)
	}

	var did = v1.Did{
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
		Context:              didMsg.Context,
	}

	metadata := v1.NewMetadata(ctx)
	metadata.Created = oldStateValue.Metadata.Created
	metadata.Deactivated = oldStateValue.Metadata.Deactivated

	if err = k.SetDid(ctx, did, &metadata); err != nil {
		return nil, err
	}

	return &v1.MsgUpdateDidResponse{
		Id: didMsg.Id,
	}, nil
}

func (k msgServer) VerifySignatureOnDidUpdate(ctx *sdk.Context, oldDIDDoc *v1.Did, newDIDDoc *v1.MsgUpdateDidPayload, signatures []*v1.SignInfo) error {
	var signers = newDIDDoc.GetSigners()

	// Get Old DID Doc controller if it's nil then assign self
	oldController := oldDIDDoc.Controller
	if len(oldController) == 0 {
		oldController = []string{oldDIDDoc.Id}
	}

	// Get New DID Doc controller if it's nil then assign self
	newController := newDIDDoc.Controller
	if len(newController) == 0 {
		newController = []string{newDIDDoc.Id}
	}

	// DID Doc controller has been changed
	if removedControllers := strings.Complement(oldController, newController); len(removedControllers) > 0 {
		for _, controller := range removedControllers {
			signers = append(signers, v1.Signer{Signer: controller})
		}
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

func AppendSignerIfNeed(signers []v1.Signer, controller string, msg *v1.MsgUpdateDidPayload) []v1.Signer {
	for _, signer := range signers {
		if signer.Signer == controller {
			return signers
		}
	}

	signer := v1.Signer{
		Signer: controller,
	}

	if controller == msg.Id {
		signer.VerificationMethod = msg.VerificationMethod
		signer.Authentication = msg.Authentication
	}

	return append(signers, signer)
}
