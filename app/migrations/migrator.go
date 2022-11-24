package migrations

import (
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	context    MigrationContext
	migrations []Migration
}

func NewMigrator(
	context MigrationContext,
	migrations []Migration,
) Migrator {
	return Migrator{
		context:    context,
		migrations: migrations,
	}
}

func (m *Migrator) Migrate(ctx sdk.Context) error {
	for _, migration := range m.migrations {
		err := migration(ctx, m.context)
		if err != nil {
			return err
		}
	}

	return nil
}

type Migration func(sctx sdk.Context, mctx MigrationContext) error

type MigrationContext struct {
	codec codec.Codec

	didStoreKey      *storetypes.KVStoreKey
	resourceStoreKey *storetypes.KVStoreKey

	didKeeper      didkeeper.Keeper
	resourceKeeper resourcekeeper.Keeper
}

func NewMigrationContext(
	codec codec.Codec,
	didStoreKey *storetypes.KVStoreKey,
	resourceStoreKey *storetypes.KVStoreKey,
	didKeeper didkeeper.Keeper,
	resourceKeeper resourcekeeper.Keeper,
) MigrationContext {
	return MigrationContext{
		codec: codec,

		didStoreKey:      didStoreKey,
		resourceStoreKey: resourceStoreKey,

		didKeeper:      didKeeper,
		resourceKeeper: resourceKeeper,
	}
}
