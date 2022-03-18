package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate2to3 migrates x/cheqd storage from version 2 to 3.
func (m *Migrator) Migrate2to3(ctx sdk.Context) error {

	store := ctx.KVStore(m.keeper.storeKey)
	byteKey := types.KeyPrefix(types.DidNamespaceKey)

	oldKey := "testnettestnet"
	namespace := m.keeper.GetFromState(ctx, oldKey)
	store.Delete(byteKey)
	m.keeper.SetDidNamespace(ctx, namespace)
	return nil
}
