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

func (k *Keeper) VerifySignature(ctx *sdk.Context, msg types.IdentityMsg, signers []types.Signer, signatures []*types.SignInfo) error {
	if len(signers) == 0 {
		return types.ErrInvalidSignature.Wrap("At least one signer should be present")
	}

	if len(signatures) == 0 {
		return types.ErrInvalidSignature.Wrap("At least one signature should be present")
	}

	signingInput := msg.GetSignBytes()

	for _, signer := range signers {
		if signer.VerificationMethod == nil {
			state, err := k.GetDid(ctx, signer.Signer)
			if err != nil {
				return types.ErrDidDocNotFound.Wrap(signer.Signer)
			}

			didDoc, err := state.GetDid()
			if err != nil {
				return types.ErrDidDocNotFound.Wrap(signer.Signer)
			}

			signer.Authentication = didDoc.Authentication
			signer.VerificationMethod = didDoc.VerificationMethod
		}

		valid, err := VerifyIdentitySignature(signer, signatures, signingInput)
		if err != nil {
			return sdkerrors.Wrap(types.ErrInvalidSignature, err.Error())
		}

		if !valid {
			return sdkerrors.Wrap(types.ErrInvalidSignature, signer.Signer)
		}
	}

	return nil
}

func (k *Keeper) ValidateController(ctx *sdk.Context, id string, controller string) error {
	if id == controller {
		return nil
	}
	state, err := k.GetDid(ctx, controller)
	if err != nil {
		return types.ErrDidDocNotFound.Wrap(controller)
	}
	didDoc, err := state.GetDid()
	if err != nil {
		return types.ErrDidDocNotFound.Wrap(controller)
	}
	if len(didDoc.Authentication) < 0 {
		return types.ErrBadRequestInvalidVerMethod.Wrap(
			fmt.Sprintf("Verificatition method controller %s doesn't have an authentication keys", controller))
	}
	return nil
}

func VerifyIdentitySignature(signer types.Signer, signatures []*types.SignInfo, signingInput []byte) (bool, error) {
	result := true
	foundOne := false

	for _, info := range signatures {
		did, _ := utils.SplitDidUrlIntoDidAndFragment(info.VerificationMethodId)
		if did == signer.Signer {
			pubKey, err := FindPublicKey(signer, info.VerificationMethodId)
			if err != nil {
				return false, err
			}

			signature, err := base64.StdEncoding.DecodeString(info.Signature)
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
