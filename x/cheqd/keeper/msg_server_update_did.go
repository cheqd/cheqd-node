package keeper

import (
	"context"
	"reflect"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	// Retrieve existing state value and did
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

	// Get sign bytes before modifying payload
	signBytes := msg.Payload.GetSignBytes()

	// Construct the new version of the DID and temporary rename it and its self references
	// in order to consider old and new versions different DIDs during signatures validation
	updatedDid := msg.Payload.ToDid()
	updatedDid.ReplaceIds(updatedDid.Id, updatedDid.Id+UpdatedPostfix)

	updatedMetadata := existingStateValue.Metadata
	updatedMetadata.Update(ctx)

	updatedStateValue, err := types.NewStateValue(&updatedDid, updatedMetadata)
	if err != nil {
		return nil, err
	}

	// Consider the new version of the DID a separate DID
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
	// Duplicate signatures that reference the old version, make them reference a new (in memory) version
	signers := GetSignerDIDsForDIDUpdate(*existingDid, updatedDid)
	extendedSignatures := DuplicateSignatures(msg.Signatures, existingDid.Id, updatedDid.Id)
	for _, signer := range signers {
		signaturesBySigner := types.FindSignInfosBySigner(extendedSignatures, signer)

		if len(signaturesBySigner) == 0 {
			return nil, types.ErrSignatureNotFound.Wrapf("there should be at least one signature by %s", signer)
		}

		found := false
		for _, signature := range signaturesBySigner {
			err := VerifySignature(&k.Keeper, &ctx, inMemoryDids, signBytes, signature)
			if err == nil {
				found = true
				break
			}
		}

		if !found {
			return nil, types.ErrSignatureNotFound.Wrapf("there should be at least one valid signature by %s", signer)
		}
	}

	// Apply changes: return original id and modify state
	updatedDid.ReplaceIds(updatedDid.Id, existingDid.Id)
	err = k.SetDid(&ctx, &updatedDid, &updatedMetadata)
	if err != nil {
		return nil, err
	}

	// Build and return response
	return &types.MsgUpdateDidResponse{
		Id: updatedDid.Id,
	}, nil
}

func DuplicateSignatures(signatures []*types.SignInfo, didToDuplicate string, newDid string) []*types.SignInfo {
	var result []*types.SignInfo

	for _, signature := range signatures {
		result = append(result, signature)

		did, path, query, fragment := utils.MustSplitDIDUrl(signature.VerificationMethodId)
		if did == didToDuplicate {
			duplicate := types.SignInfo{
				VerificationMethodId: utils.JoinDIDUrl(newDid, path, query, fragment),
				Signature:            signature.Signature,
			}

			result = append(result, &duplicate)
		}
	}

	return result
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
