package keeper

import (
	"context"
	"fmt"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDid(goCtx context.Context, msg *types.MsgWriteRequest) (*types.MsgCreateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	didMsg, isMsgIdentity := msg.Data.GetCachedValue().(*types.MsgCreateDid)

	if !isMsgIdentity {
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}

	// Checks that the element exists
	if err := k.HasDidDoc(ctx, didMsg.Id); err != nil {
		return nil, err
	}

	k.AppendDid(
		ctx,
		didMsg.Id,
		didMsg.Controller,
		didMsg.VerificationMethod,
		didMsg.Authentication,
		didMsg.AssertionMethod,
		didMsg.CapabilityInvocation,
		didMsg.CapabilityDelegation,
		didMsg.KeyAgreement,
		didMsg.AlsoKnownAs,
		didMsg.Service,
	)

	return &types.MsgCreateDidResponse{
		Id: didMsg.Id,
	}, nil
}

func (k msgServer) UpdateDid(goCtx context.Context, msg *types.MsgWriteRequest) (*types.MsgUpdateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	didMsg, isMsgIdentity := msg.Data.GetCachedValue().(*types.MsgUpdateDid)

	if !isMsgIdentity {
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}

	// Checks that the element exists
	if !k.HasDid(ctx, didMsg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", didMsg.Id))
	}

	_, metadata, err := k.GetDid(ctx, didMsg.Id)
	if err != nil {
		return nil, err
	}

	versionId, exists := msg.Metadata["versionId"]
	if !exists {
		return nil, sdkerrors.Wrap(types.ErrUnexpectedDidVersion, "Metadata doesn't contain `versionId`")
	}

	if metadata.VersionId != versionId {
		errMsg := fmt.Sprintf("Ecpected %s with version %s. Got version %s", didMsg.Id, metadata.VersionId, versionId)
		return nil, sdkerrors.Wrap(types.ErrUnexpectedDidVersion, errMsg)
	}

	var did = types.Did{
		Id:                   didMsg.Id,
		Controller:           didMsg.Controller,
		VerificationMethod:   didMsg.VerificationMethod,
		Authentication:       didMsg.Authentication,
		AssertionMethod:      didMsg.AssertionMethod,
		CapabilityInvocation: didMsg.CapabilityInvocation,
		CapabilityDelegation: didMsg.CapabilityDelegation,
		KeyAgreement:         didMsg.KeyAgreement,
		AlsoKnownAs:          didMsg.AlsoKnownAs,
		Service:              didMsg.Service,
	}

	k.SetDid(ctx, did, metadata)

	return &types.MsgUpdateDidResponse{
		Id: didMsg.Id,
	}, nil
}
