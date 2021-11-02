package keeper

import (
	"crypto/ed25519"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetDidPrefix(ctx sdk.Context) string {
	prefix := types.DidPrefix + ":" + types.DidMethod + ":"
	namespace := k.GetDidNamespace(ctx)
	if len(namespace) > 0 {
		prefix = prefix + namespace + ":"
	}
	return prefix
}

func FindPublicKey(signer types.Signer, id string) (ed25519.PublicKey, error) {
	for _, authentication := range signer.Authentication {
		if authentication == id {
			vm := FindVerificationMethod(signer.VerificationMethod, id)
			if vm == nil {
				return nil, types.ErrVerificationMethodNotFound.Wrap(id)
			}
			return vm.GetPublicKey()
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
