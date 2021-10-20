package keeper

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k *Keeper) VerifySignature(ctx *sdk.Context, msg *types.MsgWriteRequest, signers []types.Signer) error {
	if len(signers) == 0 {
		return types.ErrInvalidSignature.Wrap("At least one signer should be present")
	}

	signingInput, err := BuildSigningInput(k.cdc, msg)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidSignature, err.Error())
	}

	for _, signer := range signers {
		if signer.VerificationMethod == nil {
			state, err := k.GetDid(ctx, signer.Signer)
			if err != nil {
				return types.ErrDidDocNotFound.Wrap(signer.Signer)
			}

			didDoc := state.GetDid()
			if didDoc == nil {
				return types.ErrDidDocNotFound.Wrap(signer.Signer)
			}

			signer.Authentication = state.GetDid().Authentication
			signer.VerificationMethod = state.GetDid().VerificationMethod
		}

		valid, err := VerifyIdentitySignature(signer, msg.Signatures, signingInput)
		if err != nil {
			return sdkerrors.Wrap(types.ErrInvalidSignature, err.Error())
		}

		if !valid {
			return sdkerrors.Wrap(types.ErrInvalidSignature, signer.Signer)
		}
	}

	return nil
}

func VerifyIdentitySignature(signer types.Signer, signatures map[string]string, signingInput []byte) (bool, error) {
	result := true
	foundOne := false

	for id, signature := range signatures {
		did, _ := utils.SplitDidUrlIntoDidAndFragment(id)
		if did == signer.Signer {
			pubKey, err := FindPublicKey(signer, id)
			if err != nil {
				return false, err
			}

			signature, err := base64.StdEncoding.DecodeString(signature)
			if err != nil {
				return false, err
			}

			result = result && ed25519.Verify(pubKey, signingInput, signature)
			foundOne = true
		}
	}

	if !foundOne {
		return false, fmt.Errorf("signature %s not found", signer.Signer)
	}

	return result, nil
}
