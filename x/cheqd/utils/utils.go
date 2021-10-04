package utils

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"strings"
)

func CompareOwners(authors []string, controllers []string) bool {
	type void struct{}
	var member void

	controllerSet := make(map[string]void)
	for _, author := range authors {
		controllerSet[author] = member
	}
	result := true
	for _, controller := range controllers {
		_, exists := controllerSet[controller]
		result = result && exists
	}

	return result
}

func SplitDidUrlIntoDidAndFragment(didUrl string) (string, string) {
	fragments := strings.Split(didUrl, "#")
	return fragments[0], fragments[1]
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
