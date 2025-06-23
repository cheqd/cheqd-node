package keeper

import (
	"github.com/cheqd/cheqd-node/x/resource/exported"
	v4 "github.com/cheqd/cheqd-node/x/resource/migration/v4"
	v5 "github.com/cheqd/cheqd-node/x/resource/migration/v5"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place state migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
}

// NewMigrator returns Migrator instance for the state migration.
func NewMigrator(k Keeper, ss exported.Subspace) Migrator {
	return Migrator{
		keeper:         k,
		legacySubspace: ss,
	}
}

func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v4.MigrateStore(ctx, m.keeper.storeService, m.legacySubspace, m.keeper.cdc,
		m.keeper.ResourceCount, m.keeper.ResourceMetadata, m.keeper.ResourceData, m.keeper.LatestResourceVersion)
}

func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return v5.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc)
}
