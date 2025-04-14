package keeper

import (
	"context"
)

// GetDidNamespace get did namespace
func (k Keeper) GetDidNamespace(ctx *context.Context) (string, error) {
	return k.DidNamespace.Get(*ctx)
}

// SetDidNamespace set did namespace
func (k Keeper) SetDidNamespace(ctx *context.Context, namespace string) error {
	return k.DidNamespace.Set(*ctx, namespace)
}
