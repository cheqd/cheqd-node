package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/did/utils"
)

// GetDidNamespace get did namespace
func (k Keeper) GetDidNamespace(ctx *context.Context) (string, error) {
	store := k.storeService.OpenKVStore(*ctx)

	key := utils.StrBytes(types.DidNamespaceKey)
	value, err := store.Get(key)
	if err != nil {
		return "", err
	}

	return string(value), nil
}

// SetDidNamespace set did namespace
func (k Keeper) SetDidNamespace(ctx *context.Context, namespace string) {
	store := k.storeService.OpenKVStore(*ctx)

	key := utils.StrBytes(types.DidNamespaceKey)
	value := []byte(namespace)

	store.Set(key, value)
}
