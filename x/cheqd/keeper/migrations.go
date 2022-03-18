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

const OldDidNamespaceKey = "testnettestnet"

// Migrate2to3 migrates x/cheqd storage from version 2 to 3.
func (m *Migrator) Migrate2to3(ctx sdk.Context) error {
	store := ctx.KVStore(m.keeper.storeKey)

	namespace := store.Get([]byte(OldDidNamespaceKey))
	store.Delete([]byte(OldDidNamespaceKey))

	store.Set([]byte(types.DidNamespaceKey), namespace)

	return nil
}
