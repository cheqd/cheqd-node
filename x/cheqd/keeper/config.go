package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetDidNamespace get the total number of did
func (k Keeper) GetDidNamespace(ctx sdk.Context) string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), v1.KeyPrefix(v1.DidNamespace))
	byteKey := v1.KeyPrefix(v1.DidNamespace)
	bz := store.Get(byteKey)

	// Parse bytes
	namespace := string(bz)
	return namespace
}

// SetDidNamespace set did namespace
func (k Keeper) SetDidNamespace(ctx sdk.Context, namespace string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), v1.KeyPrefix(v1.DidNamespace))
	byteKey := v1.KeyPrefix(v1.DidNamespace)
	bz := []byte(namespace)
	store.Set(byteKey, bz)
}
