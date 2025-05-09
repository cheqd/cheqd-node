package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/did/utils"
)

func (k MsgServer) CreateDidDoc(goCtx context.Context, msg *types.MsgCreateDidDoc) (*types.MsgCreateDidDocResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Get sign bytes before modifying payload
	signBytes := msg.Payload.GetSignBytes()

	// Normalize UUID identifiers
	msg.Normalize()

	hasDidDoc, err := k.HasDidDoc(goCtx, msg.Payload.Id)
	if err != nil {
		return nil, err
	}
	// Validate DID doesn't exist
	if hasDidDoc {
		return nil, types.ErrDidDocExists.Wrap(msg.Payload.Id)
	}

	// Validate namespaces
	namespace, err := k.GetDidNamespace(goCtx)
	if err != nil {
		return nil, err
	}
	err = msg.Validate([]string{namespace})
	if err != nil {
		return nil, types.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Build metadata and stateValue
	didDoc := msg.Payload.ToDidDoc()
	metadata := types.NewMetadataFromContext(goCtx, msg.Payload.VersionId)
	didDocWithMetadata := types.NewDidDocWithMetadata(&didDoc, &metadata)

	// Consider did that we are going to create during did resolutions
	inMemoryDids := map[string]types.DidDocWithMetadata{didDoc.Id: didDocWithMetadata}

	// Check controllers' existence
	controllers := didDoc.AllControllerDids()
	for _, controller := range controllers {
		_, err := MustFindDidDoc(&k.Keeper, goCtx, inMemoryDids, controller)
		if err != nil {
			return nil, err
		}
	}

	// Verify signatures
	signers := GetSignerDIDsForDIDCreation(didDoc)
	err = VerifyAllSignersHaveAllValidSignatures(&k.Keeper, goCtx, inMemoryDids, signBytes, signers, msg.Signatures)
	if err != nil {
		return nil, err
	}

	// Save first DIDDoc version
	err = k.AddNewDidDocVersion(goCtx, &didDocWithMetadata)
	if err != nil {
		return nil, types.ErrInternal.Wrap(err.Error())
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
