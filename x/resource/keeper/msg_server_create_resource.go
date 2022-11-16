package keeper

import (
	"context"
	"crypto/sha256"
	"time"

	"github.com/cheqd/cheqd-node/x/resource/utils"

	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateResource(goCtx context.Context, msg *types.MsgCreateResource) (*types.MsgCreateResourceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Remember bytes before modifying payload
	signBytes := msg.Payload.GetSignBytes()

	msg.Normalize()

	// Validate corresponding DIDDoc exists
	namespace := k.didKeeper.GetDidNamespace(&ctx)
	did := didutils.JoinDID(didtypes.DidMethod, namespace, msg.Payload.CollectionId)
	didDoc, err := k.didKeeper.GetLatestDidDoc(&ctx, did)
	if err != nil {
		return nil, err
	}

	// Validate namespaces
	err = msg.Validate([]string{namespace})
	if err != nil {
		return nil, didtypes.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Validate DID is not deactivated
	if didDoc.Metadata.Deactivated {
		return nil, didtypes.ErrDIDDocDeactivated.Wrap(did)
	}

	// Validate Resource doesn't exist
	if k.HasResource(&ctx, msg.Payload.CollectionId, msg.Payload.Id) {
		return nil, types.ErrResourceExists.Wrap(msg.Payload.Id)
	}

	// We can use the same signers as for DID creation because didDoc stays the same
	signers := didkeeper.GetSignerDIDsForDIDCreation(*didDoc.DidDoc)
	err = didkeeper.VerifyAllSignersHaveAllValidSignatures(&k.didKeeper, &ctx, map[string]didtypes.DidDocWithMetadata{},
		signBytes, signers, msg.Signatures)
	if err != nil {
		return nil, err
	}

	// Build Resource
	resource := msg.Payload.ToResource()
	checksum := sha256.Sum256([]byte(resource.Resource.Data))
	resource.Metadata.Checksum = checksum[:]
	resource.Metadata.Created = ctx.BlockTime().Format(time.RFC3339)
	resource.Metadata.MediaType = utils.DetectMediaType(resource.Resource.Data)

	// Find previous version and upgrade backward and forward version links
	previousResourceVersionHeader, found := k.GetLastResourceVersionMetadata(&ctx, resource.Metadata.CollectionId, resource.Metadata.Name, resource.Metadata.ResourceType)
	if found {
		// Set links
		previousResourceVersionHeader.NextVersionId = resource.Metadata.Id
		resource.Metadata.PreviousVersionId = previousResourceVersionHeader.Id

		// Update previous version
		err := k.UpdateResourceMetadata(&ctx, &previousResourceVersionHeader)
		if err != nil {
			return nil, err
		}
	}

	// Persist resource
	err = k.SetResource(&ctx, &resource)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgCreateResourceResponse{
		Resource: resource.Metadata,
	}, nil
}
