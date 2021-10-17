package keeper

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

func BuildSigningInput(msg *types.MsgWriteRequest) ([]byte, error) {
	metadataBytes, err := json.Marshal(&msg.Metadata)
	if err != nil {
		return nil, types.ErrInvalidSignature.Wrap("An error has occurred during metadata marshalling")
	}

	dataBytes := msg.Data.Value
	signingInput := ([]byte)(base64.StdEncoding.EncodeToString(metadataBytes) + base64.StdEncoding.EncodeToString(dataBytes))
	return signingInput, nil
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

func FindVerificationMethod(vms []*types.VerificationMethod, id string) *types.VerificationMethod {
	for _, vm := range vms {
		if vm.Id == id {
			return vm
		}
	}

	return nil
}
