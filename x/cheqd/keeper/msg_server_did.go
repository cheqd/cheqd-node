package keeper

import (
	"context"
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/utils/strings"
	"reflect"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDid(goCtx context.Context, msg *types.MsgWriteRequest) (*types.MsgCreateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	didMsg, isMsgIdentity := msg.Data.GetCachedValue().(*types.MsgCreateDid)

	if !isMsgIdentity {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
	}

	if err := didMsg.ValidateBasic(); err != nil {
		return nil, err
	}

	if err := k.VerifySignature(&ctx, msg, didMsg.GetSigners()); err != nil {
		return nil, err
	}

	// Checks that the element exists
	if err := k.HasDidDoc(ctx, didMsg.Id); err != nil {
		return nil, err
	}

	k.AppendDid(
		ctx,
		didMsg.Id,
		didMsg.Controller,
		didMsg.VerificationMethod,
		didMsg.Authentication,
		didMsg.AssertionMethod,
		didMsg.CapabilityInvocation,
		didMsg.CapabilityDelegation,
		didMsg.KeyAgreement,
		didMsg.AlsoKnownAs,
		didMsg.Service,
		didMsg.Context,
	)

	return &types.MsgCreateDidResponse{
		Id: didMsg.Id,
	}, nil
}

func (k msgServer) UpdateDid(goCtx context.Context, msg *types.MsgWriteRequest) (*types.MsgUpdateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	didMsg, isMsgIdentity := msg.Data.GetCachedValue().(*types.MsgUpdateDid)

	if !isMsgIdentity {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
	}

	if err := didMsg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Checks that the element exists
	if !k.HasDid(ctx, didMsg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", didMsg.Id))
	}

	oldDIDDoc, metadata, err := k.GetDid(&ctx, didMsg.Id)
	if err != nil {
		return nil, err
	}

	if err := k.UpdateDidVerifySignature(&ctx, msg, oldDIDDoc, didMsg); err != nil {
		return nil, err
	}

	versionId, exists := msg.Metadata["versionId"]
	if !exists {
		return nil, sdkerrors.Wrap(types.ErrUnexpectedDidVersion, "Metadata doesn't contain `versionId`")
	}

	if metadata.VersionId != versionId {
		errMsg := fmt.Sprintf("Ecpected %s with version %s. Got version %s", didMsg.Id, metadata.VersionId, versionId)
		return nil, sdkerrors.Wrap(types.ErrUnexpectedDidVersion, errMsg)
	}

	var did = types.Did{
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

	k.SetDid(ctx, did, metadata)

	return &types.MsgUpdateDidResponse{
		Id: didMsg.Id,
	}, nil
}

func (k msgServer) UpdateDidVerifySignature(ctx *sdk.Context, msg *types.MsgWriteRequest, oldDIDDoc *types.Did, newDIDDoc *types.MsgUpdateDid) error {
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
			signers = append(signers, types.Signer{Signer: controller})
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

	if err := k.VerifySignature(ctx, msg, signers); err != nil {
		return err
	}

	return nil
}

func AppendSignerIfNeed(signers []types.Signer, controller string, msg *types.MsgUpdateDid) []types.Signer {
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
