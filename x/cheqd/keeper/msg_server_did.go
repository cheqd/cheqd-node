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

	// Checks that the element exists
	if !k.HasDid(ctx, didMsg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", didMsg.Id))
	}

	k.SetDid(ctx, did)

	return &types.MsgUpdateDidResponse{
		Id: didMsg.Id,
	}, nil
}
