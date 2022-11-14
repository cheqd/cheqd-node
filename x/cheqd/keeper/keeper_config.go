package keeper

import (
	"github.com/canow-co/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetDidNamespace get did namespace
func (k Keeper) GetDidNamespace(ctx *sdk.Context) string {
	return k.GetFromState(ctx, types.DidNamespaceKey)
}

// SetDidNamespace set did namespace
func (k Keeper) SetDidNamespace(ctx *sdk.Context, namespace string) {
	k.SetToState(ctx, types.DidNamespaceKey, []byte(namespace))
}

// GetFromState - get State value
func (k Keeper) GetFromState(ctx *sdk.Context, stateKey string) string {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(stateKey)
	bz := store.Get(byteKey)

	// Parse bytes
	namespace := string(bz)
	return namespace
}

// SetToState - set State value
func (k Keeper) SetToState(ctx *sdk.Context, stateKey string, stateValue []byte) {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(stateKey)
	store.Set(byteKey, stateValue)
}

// DeleteFromState - remove value from State by key
func (k Keeper) DeleteFromState(ctx *sdk.Context, stateKey string) {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(stateKey)
	store.Delete(byteKey)
}
