package keeper

import (
	"context"

	"github.com/Workiva/go-datastructures/threadsafe/err"
	cheqd_types "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/cheqd-node/x/resource/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateResource(goCtx context.Context, msg *types.MsgCreateResource) (*types.MsgCreateResourceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate Resource doesn't exist
	if k.HasResource(&ctx, msg.Payload.CollectionId, msg.Payload.Id) {
		return nil, types.ErrResourceExists.Wrap(msg.Payload.Id)
	}

	// Build Resource
	resource := msg.Payload.ToResource()

	// Consider resource that we are going to create during resource resolutions
	inMemoryResources := map[string]types.Resource{resource.Id: resource}

	// Verify signatures
	signers := GetSignerDIDsForDIDCreation(resource)
	for _, signer := range signers {
		signature, found := cheqd_types.FindSignInfoBySigner(msg.Signatures, signer)

		if !found {
			return nil, cheqd_types.ErrSignatureNotFound.Wrapf("signer: %s", signer)
		}

		err := VerifySignature(&k.Keeper, &ctx, inMemoryResources, msg.Payload.GetSignBytes(), signature)
		if err != nil {
			return nil, err
		}
	}

	// Apply changes
	err = k.AppendResource(&ctx, &resource)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgCreateResourceResponse{
		Id: resource.Id,
	}, nil
}

func GetSignerDIDsForDIDCreation(resource types.Resource) []string {
	res := resource.GetControllersOrSubject()
	res = append(res, resource.GetVerificationMethodControllers()...)

	return utils.UniqueSorted(res)
}
