package utils

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"strings"
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

func SplitDidUrlIntoDidAndFragment(didUrl string) (string, string) {
	fragments := strings.Split(didUrl, "#")
	return fragments[0], fragments[1]
}

func VerifyIdentitySignature(controller string, authentication []string, verificationMethods []*types.VerificationMethod, signatures map[string]string, signingInput []byte) (bool, error) {
	result := true

	for signer, signature := range signatures {
		did, _ := SplitDidUrlIntoDidAndFragment(signer)
		if did == controller {
			pubKey, err := FindPublicKey(authentication, verificationMethods, signer)
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

func FindPublicKey(authentication []string, verificationMethods []*types.VerificationMethod, id string) (ed25519.PublicKey, error) {
	for _, authentication := range authentication {
		if authentication == id {
			for _, vm := range verificationMethods {
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
