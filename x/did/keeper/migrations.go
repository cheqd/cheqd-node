package keeper

import (
	"github.com/cheqd/cheqd-node/x/did/exported"
	v5 "github.com/cheqd/cheqd-node/x/did/migrations/v5"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, legacySubspace exported.Subspace) Migrator {
	return Migrator{keeper: keeper, legacySubspace: legacySubspace}
}

// module state.
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return v5.MigrateStore(ctx, m.keeper.storeService, m.legacySubspace, m.keeper.cdc,
		m.keeper.DidCount, m.keeper.DidDocuments)
}
