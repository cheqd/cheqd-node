package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k MsgServer) DeactivateDidDoc(goCtx context.Context, msg *types.MsgDeactivateDidDoc) (*types.MsgDeactivateDidDocResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate DID does exist
	if !k.HasDidDoc(&ctx, msg.Payload.Id) {
		return nil, types.ErrDidDocNotFound.Wrap(msg.Payload.Id)
	}

	// Validate namespaces
	namespace := k.GetDidNamespace(&ctx)
	err := msg.Validate([]string{namespace})
	if err != nil {
		return nil, types.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Retrieve didDoc state value and did
	didDoc, err := k.GetDidDoc(&ctx, msg.Payload.Id)
	if err != nil {
		return nil, err
	}

	// Validate DID is not deactivated
	if didDoc.Metadata.Deactivated {
		return nil, types.ErrDIDDocDeactivated.Wrap(msg.Payload.Id)
	}

	// We neither create dids nor update
	inMemoryDids := map[string]types.DidDocWithMetadata{}

	// Verify signatures
	signers := GetSignerDIDsForDIDCreation(*didDoc.DidDoc)
	err = VerifyAllSignersHaveAllValidSignatures(&k.Keeper, &ctx, inMemoryDids, msg.Payload.GetSignBytes(), signers, msg.Signatures)
	if err != nil {
		return nil, err
	}

	// Update metadata
	didDoc.Metadata.Deactivated = true
	didDoc.Metadata.Update(ctx)

	// Apply changes
	err = k.SetDidDoc(&ctx, &didDoc)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgDeactivateDidDocResponse{
		Value: &didDoc,
	}, nil
}
