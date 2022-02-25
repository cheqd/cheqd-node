package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}


func AppendSignerIfNeed(signers []types.Signer, controller string, msg *types.MsgUpdateDidPayload) []types.Signer {
	for _, signer := range signers {
		if signer.Signer == controller {
			return signers
		}
	}

	signer := types.Signer{
		Signer: controller,
	}

	if controller == msg.Id {
		signer.VerificationMethod = msg.VerificationMethod
		signer.Authentication = msg.Authentication
	}

	return append(signers, signer)
}

func (k msgServer) ValidateDidControllers(ctx *sdk.Context, id string, controllers []string, verMethods []*types.VerificationMethod) error {

	for _, verificationMethod := range verMethods {
		if err := k.ValidateController(ctx, id, verificationMethod.Controller); err != nil {
			return err
		}
	}

	for _, didController := range controllers {
		if err := k.ValidateController(ctx, id, didController); err != nil {
			return err
		}
	}
	return nil
}
