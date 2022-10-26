package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeactivateDid(goCtx context.Context, msg *types.MsgDeactivateDid) (*types.MsgDeactivateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate DID does exist
	if !k.HasDid(&ctx, msg.Payload.Id) {
		return nil, types.ErrDidDocNotFound.Wrap(msg.Payload.Id)
	}

	// Validate namespaces
	namespace := k.GetDidNamespace(&ctx)
	err := msg.Validate([]string{namespace})
	if err != nil {
		return nil, types.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Retrieve existing state value and did
	existingStateValue, err := k.GetDid(&ctx, msg.Payload.Id)
	if err != nil {
		return nil, err
	}

	// Validate DID is not deactivated
	if existingStateValue.Metadata.Deactivated {
		return nil, types.ErrDIDDocDeactivated.Wrap(msg.Payload.Id)
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

	// We neither create dids nor update
	inMemoryDids := map[string]types.StateValue{}

	// Verify signatures
	signers := GetSignerDIDsForDIDCreation(*existingDid)
	err = VerifyAllSignersHaveAllValidSignatures(&k.Keeper, &ctx, inMemoryDids, msg.Payload.GetSignBytes(), signers, msg.Signatures)
	if err != nil {
		return nil, err
	}

	// Build new state value
	updatedStateValue, err = types.NewStateValue(existingDid, &updatedMetadata)
	if err != nil {
		return nil, err
	}

	// Modify state
	err = k.SetDid(&ctx, &updatedStateValue)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgDeactivateDidResponse{
		Did:      existingDid,
		Metadata: &updatedMetadata,
	}, nil
}
