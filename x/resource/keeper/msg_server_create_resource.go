package keeper

import (
	"context"
	"crypto/sha256"
	"time"

	"github.com/cheqd/cheqd-node/x/resource/utils"

	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdutils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateResource(goCtx context.Context, msg *types.MsgCreateResource) (*types.MsgCreateResourceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Remember bytes before modifying payload
	signBytes := msg.Payload.GetSignBytes()

	msg.Normalize()

	// Validate corresponding DIDDoc exists
	namespace := k.cheqdKeeper.GetDidNamespace(&ctx)
	did := cheqdutils.JoinDID(cheqdtypes.DidMethod, namespace, msg.Payload.CollectionId)
	didDocStateValue, err := k.cheqdKeeper.GetDid(&ctx, did)
	if err != nil {
		return nil, err
	}

	// Validate DID is not deactivated
	if didDocStateValue.Metadata.Deactivated {
		return nil, cheqdtypes.ErrDIDDocDeactivated.Wrap(did)
	}

	// Validate Resource doesn't exist
	if k.HasResource(&ctx, msg.Payload.CollectionId, msg.Payload.Id) {
		return nil, types.ErrResourceExists.Wrap(msg.Payload.Id)
	}

	// Validate signatures
	didDoc, err := didDocStateValue.UnpackDataAsDid()
	if err != nil {
		return nil, err
	}

	// We can use the same signers as for DID creation because didDoc stays the same
	signers := cheqdkeeper.GetSignerDIDsForDIDCreation(*didDoc)
	err = cheqdkeeper.VerifyAllSignersHaveAllValidSignatures(&k.cheqdKeeper, &ctx, map[string]cheqdtypes.StateValue{},
		signBytes, signers, msg.Signatures)
	if err != nil {
		return nil, err
	}

	// Build Resource
	resource := msg.Payload.ToResource()
	checksum := sha256.Sum256([]byte(resource.Data))
	resource.Header.Checksum = checksum[:]
	resource.Header.Created = ctx.BlockTime().Format(time.RFC3339)
	resource.Header.MediaType = utils.DetectMediaType(resource.Data)

	// Find previous version and upgrade backward and forward version links
	previousResourceVersionHeader, found := k.GetLastResourceVersionHeader(&ctx, resource.Header.CollectionId, resource.Header.Name, resource.Header.ResourceType)
	if found {
		// Set links
		previousResourceVersionHeader.NextVersionId = resource.Header.Id
		resource.Header.PreviousVersionId = previousResourceVersionHeader.Id

		// Update previous version
		err := k.UpdateResourceHeader(&ctx, &previousResourceVersionHeader)
		if err != nil {
			return nil, err
		}
	}

	// Append backlink to didDoc
	didDocStateValue.Metadata.Resources = append(didDocStateValue.Metadata.Resources, resource.Header.Id)
	err = k.cheqdKeeper.SetDid(&ctx, &didDocStateValue)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Persist resource
	err = k.SetResource(&ctx, &resource)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgCreateResourceResponse{
		Resource: &resource,
	}, nil
}
