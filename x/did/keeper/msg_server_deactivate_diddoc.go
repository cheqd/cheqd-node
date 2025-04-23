package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
)

func (k MsgServer) DeactivateDidDoc(goCtx context.Context, msg *types.MsgDeactivateDidDoc) (*types.MsgDeactivateDidDocResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	// Get sign bytes before modifying payload
	signBytes := msg.Payload.GetSignBytes()

	// Normalize UUID identifiers
	msg.Normalize()

	// Validate DID does exist
	hasDidDoc, err := k.HasDidDoc(&goCtx, msg.Payload.Id)
	if err != nil {
		return nil, err
	}
	if !hasDidDoc {
		return nil, types.ErrDidDocNotFound.Wrap(msg.Payload.Id)
	}

	// Validate namespaces
	namespace, err := k.GetDidNamespace(&goCtx)
	if err != nil {
		return nil, err
	}
	err = msg.Validate([]string{namespace})
	if err != nil {
		return nil, types.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Retrieve didDoc state value and did
	didDoc, err := k.GetLatestDidDoc(&goCtx, msg.Payload.Id)
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
	err = VerifyAllSignersHaveAllValidSignatures(&k.Keeper, &goCtx, inMemoryDids, signBytes, signers, msg.Signatures)
	if err != nil {
		return nil, err
	}

	// Update metadata
	didDoc.Metadata.Deactivated = true
	didDoc.Metadata.Update(goCtx, msg.Payload.VersionId)

	// Apply changes. We create a new version on deactivation to track deactivation time
	err = k.AddNewDidDocVersion(&goCtx, &didDoc)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Deactivate all previous versions
	var iterationErr error
	k.IterateDidDocVersions(&goCtx, msg.Payload.Id, func(didDocWithMetadata types.DidDocWithMetadata) bool {
		didDocWithMetadata.Metadata.Deactivated = true

		err := k.SetDidDocVersion(&goCtx, &didDocWithMetadata, true)
		if err != nil {
			iterationErr = err
			return false
		}

		return true
	})

	if iterationErr != nil {
		return nil, types.ErrInternal.Wrapf(iterationErr.Error())
	}

	// Build and return response
	return &types.MsgDeactivateDidDocResponse{
		Value: &didDoc,
	}, nil
}
