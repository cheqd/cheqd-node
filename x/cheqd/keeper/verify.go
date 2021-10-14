package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k *Keeper) VerifySignature(ctx *sdk.Context, msg *types.MsgWriteRequest, signers []types.Signer) error {
	signingInput, err := utils.BuildSigningInput(msg)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidSignature, err.Error())
	}

	for _, signer := range signers {
		if signer.VerificationMethod == nil {
			didDoc, _, err := k.GetDid(ctx, signer.Signer)
			if err != nil {
				return sdkerrors.Wrap(types.ErrInvalidSignature, err.Error())
			}

			signer.Authentication = didDoc.Authentication
			signer.VerificationMethod = didDoc.VerificationMethod
		}

		valid, err := utils.VerifyIdentitySignature(signer, msg.Signatures, signingInput)
		if err != nil {
			return sdkerrors.Wrap(types.ErrInvalidSignature, err.Error())
		}

		if !valid {
			return sdkerrors.Wrap(types.ErrInvalidSignature, signer.Signer)
		}
	}

	return nil
}
