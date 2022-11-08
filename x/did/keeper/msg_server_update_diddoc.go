package keeper

import (
	"context"
	"reflect"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/did/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const UpdatedPostfix string = "-updated"

func (k MsgServer) UpdateDidDoc(goCtx context.Context, msg *types.MsgUpdateDidDoc) (*types.MsgUpdateDidDocResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get sign bytes before modifying payload
	signBytes := msg.Payload.GetSignBytes()

	// Normalize UUID identifiers
	msg.Normalize()

	// Validate namespaces
	namespace := k.GetDidNamespace(&ctx)
	err := msg.Validate([]string{namespace})
	if err != nil {
		return nil, types.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Retrieve existing state value and did
	existingDidDocWithMetadata, err := k.GetDidDoc(&ctx, msg.Payload.Id)
	if err != nil {
		return nil, err
	}

	existingDidDoc := existingDidDocWithMetadata.DidDoc

	// Validate DID is not deactivated
	if existingDidDocWithMetadata.Metadata.Deactivated {
		return nil, types.ErrDIDDocDeactivated.Wrap(msg.Payload.Id)
	}

	// Check version id
	if msg.Payload.VersionId != existingDidDocWithMetadata.Metadata.VersionId {
		return nil, types.ErrUnexpectedDidVersion.Wrapf("got: %s, must be: %s", msg.Payload.VersionId, existingDidDocWithMetadata.Metadata.VersionId)
	}

	// Construct the new version of the DID and temporary rename it and its self references
	// in order to consider old and new versions different DIDs during signatures validation
	updatedDidDoc := msg.Payload.ToDidDoc()
	updatedDidDoc.ReplaceDids(updatedDidDoc.Id, updatedDidDoc.Id+UpdatedPostfix)

	updatedMetadata := *existingDidDocWithMetadata.Metadata
	updatedMetadata.Update(ctx)

	updatedDidDocWithMetadata := types.NewDidDocWithMetadata(&updatedDidDoc, &updatedMetadata)

	// Consider the new version of the DID a separate DID
	inMemoryDids := map[string]types.DidDocWithMetadata{updatedDidDoc.Id: updatedDidDocWithMetadata}

	// Check controllers existence
	controllers := updatedDidDoc.AllControllerDids()
	for _, controller := range controllers {
		_, err := MustFindDidDoc(&k.Keeper, &ctx, inMemoryDids, controller)
		if err != nil {
			return nil, err
		}
	}

	// Verify signatures
	// Duplicate signatures that reference the old version, make them reference a new (in memory) version
	// We can't use VerifySignatures because we can't uniquely identify a verification method corresponding to a given signInfo.
	// In other words if a signature belongs to the did being updated, there is no way to know which did version it belongs to: old or new.
	// To eliminate this problem we have to add pubkey to the signInfo in future.
	signers := GetSignerDIDsForDIDUpdate(*existingDidDoc, updatedDidDoc)
	extendedSignatures := DuplicateSignatures(msg.Signatures, existingDidDocWithMetadata.DidDoc.Id, updatedDidDoc.Id)
	err = VerifyAllSignersHaveAtLeastOneValidSignature(&k.Keeper, &ctx, inMemoryDids, signBytes, signers, extendedSignatures, existingDidDoc.Id, updatedDidDoc.Id)
	if err != nil {
		return nil, err
	}

	// Return original id
	updatedDidDoc.ReplaceDids(updatedDidDoc.Id, existingDidDoc.Id)

	// Update state
	err = k.SetDidDoc(&ctx, &updatedDidDocWithMetadata)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgUpdateDidDocResponse{
		Value: &updatedDidDocWithMetadata,
	}, nil
}

func GetSignerIdForErrorMessage(signerId string, existingVersionId string, updatedVersionId string) interface{} {
	if signerId == existingVersionId { // oldDid->id
		return existingVersionId + " (old version)"
	}

	if signerId == updatedVersionId { // oldDid->id + UpdatedPrefix
		return existingVersionId + " (new version)"
	}

	return signerId
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

func GetSignerDIDsForDIDUpdate(existingDidDoc types.DidDoc, updatedDidDoc types.DidDoc) []string {
	signers := existingDidDoc.GetControllersOrSubject()
	signers = append(signers, updatedDidDoc.GetControllersOrSubject()...)

	existingVMMap := types.VerificationMethodListToMapByFragment(existingDidDoc.VerificationMethod)
	updatedVMMap := types.VerificationMethodListToMapByFragment(updatedDidDoc.VerificationMethod)

	for _, updatedVM := range updatedDidDoc.VerificationMethod {
		_, _, _, fragment := utils.MustSplitDIDUrl(updatedVM.Id)
		existingVM, found := existingVMMap[fragment]

		// VM added
		if !found {
			signers = append(signers, updatedVM.Controller)
			continue
		}

		// VM updated
		// We have to revert renaming before comparing veriifcation methods.
		// Otherwise we will detect id and controller change
		// for non changed VMs because of `-updated` postfix.
		originalUpdatedVM := *updatedVM
		originalUpdatedVM.ReplaceDids(updatedDidDoc.Id, existingDidDoc.Id)

		if !reflect.DeepEqual(existingVM, originalUpdatedVM) {
			signers = append(signers, existingVM.Controller, updatedVM.Controller)
			continue
		}

		// VM not changed
	}

	for _, existingVM := range existingDidDoc.VerificationMethod {
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
