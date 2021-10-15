package keeper

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k *Keeper) VerifySignature(ctx *sdk.Context, msg *types.MsgWriteRequest, signers []types.Signer) error {
	if len(signers) == 0 {
		return types.ErrInvalidSignature.Wrap("At least one signer should be present")
	}

	signingInput, err := BuildSigningInput(msg)
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

func BuildSigningInput(msg *types.MsgWriteRequest) ([]byte, error) {
	metadataBytes, err := json.Marshal(&msg.Metadata)
	if err != nil {
		return nil, types.ErrInvalidSignature.Wrap("An error has occurred during metadata marshalling")
	}

	dataBytes := msg.Data.Value
	signingInput := ([]byte)(base64.StdEncoding.EncodeToString(metadataBytes) + base64.StdEncoding.EncodeToString(dataBytes))
	return signingInput, nil
}

func VerifyIdentitySignature(signer types.Signer, signatures map[string]string, signingInput []byte) (bool, error) {
	result := true

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
		}
	}

	return result, nil
}

func FindPublicKey(signer types.Signer, id string) (ed25519.PublicKey, error) {
	for _, authentication := range signer.Authentication {
		if authentication == id {
			for _, vm := range signer.VerificationMethod {
				if vm.Id == id {
					return base58.Decode(vm.PublicKeyMultibase[1:]), nil
				}
			}

			msg := fmt.Sprintf("Verification Method %s not found", id)
			return nil, errors.New(msg)
		}
	}

	msg := fmt.Sprintf("Authentication %s not found", id)
	return nil, errors.New(msg)
}
