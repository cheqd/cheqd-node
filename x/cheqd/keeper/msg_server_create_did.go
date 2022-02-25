package keeper

import (
	"context"
	"fmt"
	"reflect"

	"github.com/cheqd/cheqd-node/x/cheqd/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDid(goCtx context.Context, msg *types.MsgCreateDid) (*types.MsgCreateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	prefix := k.GetDidPrefix(ctx)

	payload := msg.GetPayload()
	if err := payload.Validate(prefix); err != nil {
		return nil, err
	}

	if err := k.ValidateDidControllers(&ctx, payload.Id, payload.Controller, payload.VerificationMethod); err != nil {
		return nil, err
	}

	if err := k.VerifySignature(&ctx, payload, payload.GetSigners(), msg.GetSignatures()); err != nil {
		return nil, err
	}

	// Checks that the did doesn't exist
	if k.HasDid(ctx, payload.Id) {
		return nil, sdkerrors.Wrap(types.ErrDidDocExists, fmt.Sprintf("DID is already used by DIDDoc %s", payload.Id))
	}

	var did = types.Did{
		Context:              payload.Context,
		Id:                   payload.Id,
		Controller:           payload.Controller,
		VerificationMethod:   payload.VerificationMethod,
		Authentication:       payload.Authentication,
		AssertionMethod:      payload.AssertionMethod,
		CapabilityInvocation: payload.CapabilityInvocation,
		CapabilityDelegation: payload.CapabilityDelegation,
		KeyAgreement:         payload.KeyAgreement,
		AlsoKnownAs:          payload.AlsoKnownAs,
		Service:              payload.Service,
	}

	metadata := types.NewMetadataFromContext(ctx)
	id, err := k.AppendDid(ctx, did, &metadata)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateDidResponse{
		Id: *id,
	}, nil
}
