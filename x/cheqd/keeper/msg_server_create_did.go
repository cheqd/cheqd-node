package keeper

import (
	"context"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateDid(goCtx context.Context, msg *types.MsgCreateDid) (*types.MsgCreateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate DID doesn't exist
	if k.HasDid(&ctx, msg.Payload.Id) {
		return nil, types.ErrDidDocExists.Wrap(msg.Payload.Id)
	}

	// Validate namespaces
	namespace := k.GetDidNamespace(ctx)
	err := msg.Validate([]string{namespace})
	if err != nil {
		return nil, types.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Build metadata and stateValue
	did := msg.Payload.ToDid()
	metadata := types.NewMetadataFromContext(ctx)
	stateValue, err := types.NewStateValue(&did, &metadata)
	if err != nil {
		return nil, err
	}

	// Consider did that we are going to create during did resolutions
	inMemoryDids := map[string]types.StateValue{did.Id: stateValue}

	// Check controllers' existence
	controllers := did.AggregateControllerDids()
	for _, controller := range controllers {
		_, err := MustFindDid(&k.Keeper, &ctx, inMemoryDids, controller)

		if err != nil {
			return nil, err
		}
	}

	// Verify signatures
	signers := GetSignerDIDsForDIDCreation(did)
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
	return &types.MsgCreateDidResponse{
		Id: id,
	}, nil
}

func GetSignerDIDsForDIDCreation(did types.Did) []string {
	res := did.AggregateControllerDids()

	if len(did.Controller) == 0 {
		res = append(res, did.Id)
	}

	return utils.Unique(res)
}
