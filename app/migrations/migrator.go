package migrations

import (
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type MigrationContext struct {
	codec codec.Codec

	didKeeper      didkeeper.Keeper
	resourceKeeper resourcekeeper.Keeper
}

type CheqdMigration func(sctx sdk.Context, mctx MigrationContext) error

type CheqdMigrator struct {
	migration  CheqdMigration
	context    MigrationContext
}

func NewCheqdMigrator(
	codec codec.Codec, 
	didKeeper didkeeper.Keeper, 
	resourceKeeper resourcekeeper.Keeper, 
	migration CheqdMigration) CheqdMigrator {
	return CheqdMigrator{
		migration: migration,
		context: MigrationContext{
			codec: codec,

			didKeeper:      didKeeper,
			resourceKeeper: resourceKeeper,
		},
	}
}

func (m *CheqdMigrator) Migrate(ctx sdk.Context) error {
	err := m.migration(ctx, m.context)
	if err != nil {
		return err
	}

	return nil
}
