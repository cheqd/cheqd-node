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
	namespace := k.GetDidNamespace(&ctx)
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
	controllers := did.AllControllerDids()
	for _, controller := range controllers {
		_, err := MustFindDid(&k.Keeper, &ctx, inMemoryDids, controller)
		if err != nil {
			return nil, err
		}
	}

	// Verify signatures
	signers := GetSignerDIDsForDIDCreation(did)
	err = VerifyAllSignersHaveAllValidSignatures(&k.Keeper, &ctx, inMemoryDids, msg.Payload.GetSignBytes(), signers, msg.Signatures)
	if err != nil {
		return nil, err
	}

	// Build state value
	value, err := types.NewStateValue(&did, &metadata)
	if err != nil {
		return nil, err
	}

	err = k.SetDid(&ctx, &value)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgCreateDidResponse{
		Id: did.Id,
	}, nil
}

func GetSignerDIDsForDIDCreation(did types.Did) []string {
	res := did.GetControllersOrSubject()
	res = append(res, did.GetVerificationMethodControllers()...)

	return utils.UniqueSorted(res)
}
