package keeper

import (
	"context"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
)

const UpdatedPostfix string = "-updated"

func (k msgServer) UpdateDid(goCtx context.Context, msg *types.MsgUpdateDid) (*types.MsgUpdateDidResponse, error) {
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

	// Retrieve existing state value
	existingStateValue, err := k.GetDid(&ctx, msg.Payload.Id)
	if err != nil {
		return nil, err
	}

	existingDid, err := existingStateValue.UnpackDataAsDid()
	if err != nil {
		return nil, err
	}

	// Check version id
	if msg.Payload.VersionId != existingStateValue.Metadata.VersionId {
		return nil, types.ErrUnexpectedDidVersion.Wrapf("got: %s, must be: %s", msg.Payload.VersionId, existingStateValue.Metadata.VersionId)
	}

	// Construct updated did
	updatedDid := msg.Payload.ToDid()
	updatedDid.ReplaceId(updatedDid.Id, updatedDid.Id+UpdatedPostfix) // Temporary replace id

	updatedMetadata := types.NewMetadataFromContext(ctx)
	updatedMetadata.Created = existingStateValue.Metadata.Created
	updatedMetadata.Updated = ctx.BlockTime().String()

	updatedStateValue, err := types.NewStateValue(&updatedDid, &updatedMetadata)
	if err != nil {
		return nil, err
	}

	// Consider did that we are going to update with during did resolutions
	inMemoryDids := map[string]types.StateValue{updatedDid.Id: updatedStateValue}

	// Check controllers existence
	controllers := updatedDid.AllControllerDids()
	for _, controller := range controllers {
		_, err := MustFindDid(&k.Keeper, &ctx, inMemoryDids, controller)

		if err != nil {
			return nil, err
		}
	}

	// Verify signatures
	signers := GetSignerDIDsForDIDUpdate(*existingDid, updatedDid)
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

	// Apply changes
	updatedDid.ReplaceId(updatedDid.Id, existingDid.Id) // Return original id
	err = k.SetDid(&ctx, &updatedDid, &updatedMetadata)
	if err != nil {
		return nil, err
	}

	// Build and return response
	return &types.MsgUpdateDidResponse{
		Id: updatedDid.Id,
	}, nil
}

func GetSignerDIDsForDIDUpdate(existingDid types.Did, updatedDid types.Did) []string {
	signers := existingDid.GetControllersOrSubject()
	signers = append(signers, updatedDid.GetControllersOrSubject()...)

	existingVMMap := types.VerificationMethodListToMap(existingDid.VerificationMethod)
	updatedVMMap := types.VerificationMethodListToMap(updatedDid.VerificationMethod)

	for _, updatedVM := range updatedDid.VerificationMethod {
		existingVM, found := existingVMMap[updatedVM.Id]

		// VM added
		if !found {
			signers = append(signers, updatedVM.Controller)
			break
 		}

 		// VM updated
 		if !reflect.DeepEqual(existingVM, updatedVM) {
			signers = append(signers, existingVM.Controller, updatedVM.Controller)
 			break
		}

		// VM not changed
	}

	for _, existingVM := range existingDid.VerificationMethod {
		_, found := updatedVMMap[existingVM.Id]

		// VM removed
		if !found {
			signers = append(signers, existingVM.Controller)
			break
		}
	}

	return utils.Unique(signers)
}
