package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const DeactivatedPostfix string = "-deactivated"

func (k msgServer) DeactivateDid(goCtx context.Context, msg *types.MsgDeactivateDid) (*types.MsgDeactivateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate DID does exist
	if !k.HasDid(&ctx, msg.Payload.Id) {
		return nil, types.ErrDidDocNotFound.Wrap(msg.Payload.Id)
	}

	// Validate namespaces
	namespace := k.GetDidNamespace(ctx)
	err := msg.Validate([]string{namespace})
	if err != nil {
		return nil, types.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Retrieve existing state value and did
	existingStateValue, err := k.GetDid(&ctx, msg.Payload.Id)
	if err != nil {
		return nil, err
	}

	existingDid, err := existingStateValue.UnpackDataAsDid()
	if err != nil {
		return nil, err
	}

	

	updatedMetadata := *existingStateValue.Metadata
	updatedMetadata.Update(ctx)
	updatedMetadata.Deactivated = true

	updatedStateValue, err := types.NewStateValue(existingDid, &updatedMetadata)
	if err != nil {
		return nil, err
	}

	
	// Consider did that we are going to create during did resolutions
	inMemoryDids := map[string]types.StateValue{existingDid.Id: updatedStateValue}

	// Check controllers' existence
	controllers := existingDid.AllControllerDids()
	for _, controller := range controllers {
		_, err := MustFindDid(&k.Keeper, &ctx, inMemoryDids, controller)
		if err != nil {
			return nil, err
		}
	}

	// Verify signatures
	signers := GetSignerDIDsForDIDCreation(*existingDid)
	for _, signer := range signers {
		signature, found := types.FindSignInfoBySigner(msg.Signatures, signer)

		if !found {
			return nil, types.ErrSignatureNotFound.Wrapf("signer: %s", signer)
		}

		err := VerifySignature(&k.Keeper, &ctx, inMemoryDids, msg.Payload.GetSignBytes(), signature)
		if err != nil {
			return nil, err
		}
	}

	// Apply changes: return original id and modify state
	err = k.SetDid(&ctx, existingDid, &updatedMetadata)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgDeactivateDidResponse{
		Did: existingDid,
		Metadata: &updatedMetadata,
	}, nil
}





func GetSignerDIDsForDIDDeactivate(existingDid types.Did, updatedDid types.Did) []string { //
	signers := existingDid.GetControllersOrSubject()
	signers = append(signers, updatedDid.GetControllersOrSubject()...)

	existingVMMap := types.VerificationMethodListToMapByFragment(existingDid.VerificationMethod)
	updatedVMMap := types.VerificationMethodListToMapByFragment(updatedDid.VerificationMethod)

	for _, updatedVM := range updatedDid.VerificationMethod {
		_, _, _, fragment := utils.MustSplitDIDUrl(updatedVM.Id)
		existingVM, found := existingVMMap[fragment]

		// VM added
		if !found {
			signers = append(signers, updatedVM.Controller)
			continue
		}

		// VM updated
		// We don't compare ids because they will be different after replacing ids on the updated version of DID.
		// Fragments equality is checked above.
		if !types.CompareVerificationMethodsWithoutIds(existingVM, *updatedVM) {
			signers = append(signers, existingVM.Controller, updatedVM.Controller)
			continue
		}

		// VM not changed
	}

	for _, existingVM := range existingDid.VerificationMethod {
		_, _, _, fragment := utils.MustSplitDIDUrl(existingVM.Id)
		_, found := updatedVMMap[fragment]

		// VM removed
		if !found {
			signers = append(signers, existingVM.Controller)
			continue
		}
	}

	return utils.UniqueSorted(signers)
}
