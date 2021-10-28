package keeper

import (
	"crypto/ed25519"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

func FindPublicKey(signer types.Signer, id string) (ed25519.PublicKey, error) {
	for _, authentication := range signer.Authentication {
		if authentication == id {
			for _, vm := range signer.VerificationMethod {
				if vm.Id == id {
					return vm.GetPublicKey()
				}
			}

			return nil, types.ErrVerificationMethodNotFound.Wrap(id)
		}
	}

	return nil, types.ErrVerificationMethodNotFound.Wrap(id)
}

func FindVerificationMethod(vms []*types.VerificationMethod, id string) *types.VerificationMethod {
	for _, vm := range vms {
		if vm.Id == id {
			return vm
		}
	}

	return nil
}
