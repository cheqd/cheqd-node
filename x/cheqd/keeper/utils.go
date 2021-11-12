package keeper

import (
	"crypto/ed25519"
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetDidPrefix(ctx sdk.Context) string {
	prefix := v1.DidPrefix + ":" + v1.DidMethod + ":"
	namespace := k.GetDidNamespace(ctx)
	if len(namespace) > 0 {
		prefix = prefix + namespace + ":"
	}
	return prefix
}

func FindPublicKey(signer v1.Signer, id string) (ed25519.PublicKey, error) {
	for _, authentication := range signer.Authentication {
		if authentication == id {
			vm := FindVerificationMethod(signer.VerificationMethod, id)
			if vm == nil {
				return nil, v1.ErrVerificationMethodNotFound.Wrap(id)
			}
			return vm.GetPublicKey()
		}
	}

	return nil, v1.ErrVerificationMethodNotFound.Wrap(id)
}

func FindVerificationMethod(vms []*v1.VerificationMethod, id string) *v1.VerificationMethod {
	for _, vm := range vms {
		if vm.Id == id {
			return vm
		}
	}

	return nil
}
