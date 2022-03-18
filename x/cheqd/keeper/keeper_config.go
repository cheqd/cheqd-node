package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetDidNamespace get did namespace
func (k Keeper) GetDidNamespace(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(types.DidNamespaceKey)
	bz := store.Get(byteKey)

	// Parse bytes
	namespace := string(bz)
	return namespace
}

// SetDidNamespace set did namespace
func (k Keeper) SetDidNamespace(ctx sdk.Context, namespace string) {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(types.DidNamespaceKey)

	bz := []byte(namespace)
	store.Set(byteKey, bz)
}
