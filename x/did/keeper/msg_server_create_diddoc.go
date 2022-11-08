package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/did/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k MsgServer) CreateDidDoc(goCtx context.Context, msg *types.MsgCreateDidDoc) (*types.MsgCreateDidDocResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get sign bytes before modifying payload
	signBytes := msg.Payload.GetSignBytes()

	// Normalize UUID identifiers
	msg.Normalize()

	// Validate DID doesn't exist
	if k.HasDidDoc(&ctx, msg.Payload.Id) {
		return nil, types.ErrDidDocExists.Wrap(msg.Payload.Id)
	}

	// Validate namespaces
	namespace := k.GetDidNamespace(&ctx)
	err := msg.Validate([]string{namespace})
	if err != nil {
		return nil, types.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Build metadata and stateValue
	didDoc := msg.Payload.ToDidDoc()
	metadata := types.NewMetadataFromContext(ctx, msg.Payload.VersionId)
	didDocWithMetadata := types.NewDidDocWithMetadata(&didDoc, &metadata)

	// Consider did that we are going to create during did resolutions
	inMemoryDids := map[string]types.DidDocWithMetadata{didDoc.Id: didDocWithMetadata}

	// Check controllers' existence
	controllers := didDoc.AllControllerDids()
	for _, controller := range controllers {
		_, err := MustFindDidDoc(&k.Keeper, &ctx, inMemoryDids, controller)
		if err != nil {
			return nil, err
		}
	}

	// Verify signatures
	signers := GetSignerDIDsForDIDCreation(didDoc)
	err = VerifyAllSignersHaveAllValidSignatures(&k.Keeper, &ctx, inMemoryDids, signBytes, signers, msg.Signatures)
	if err != nil {
		return nil, err
	}

	// Save first DIDDoc version
	err = k.AddNewDidDocVersion(&ctx, &didDocWithMetadata)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgCreateDidDocResponse{
		Value: &didDocWithMetadata,
	}, nil
}

func GetSignerDIDsForDIDCreation(did types.DidDoc) []string {
	res := did.GetControllersOrSubject()
	res = append(res, did.GetVerificationMethodControllers()...)

	return utils.UniqueSorted(res)
}
