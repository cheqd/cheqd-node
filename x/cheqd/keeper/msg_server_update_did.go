package keeper

import (
	"context"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

	// Build updated did
	updatedDid := msg.Payload.ToDid()

	// Temporary update id to distinguish between links to existing and updated versions of the did
	updatedDid.ReplaceAllControllerDids(updatedDid.Id, updatedDid.Id + "-updated")

	metadata := types.NewMetadataFromContext(ctx)
	updatedStateValue, err := types.NewStateValue(&updatedDid, &metadata)
	if err != nil {
		return nil, err
	}

	// Consider did that we are going to create during did resolutions
	inMemoryDids := map[string]types.StateValue{updatedDid.Id: updatedStateValue}

	// Check controllers' existence
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
	id, err := k.AppendDid(&ctx, &did, &metadata)
	if err != nil {
		return nil, err
	}

	// Build and return response
	return &types.MsgUpdateDidResponse{
		Id: id,
	}, nil
}

func GetSignerDIDsForDIDUpdate(existing types.Did, updated types.Did) []string {
	res := existing.AllControllerDids()
	res = append(res, updated.AllControllerDids()...)

	if len(did.Controller) == 0 {
		res = append(res, did.Id)
	}

	return utils.Unique(res)
}
