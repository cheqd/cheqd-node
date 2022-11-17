package migrations

import (
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MigrationContext struct {
	didStoreKey *storetypes.KVStoreKey
	resourceStoreKey *storetypes.KVStoreKey
	codec          codec.Codec
	didKeeper      didkeeper.Keeper
	resourceKeeper resourcekeeper.Keeper
}

func NewMigrationContext(
	didStoreKey *storetypes.KVStoreKey,
	resourceStoreKey *storetypes.KVStoreKey,
	codec codec.Codec,
	didKeeper didkeeper.Keeper,
	resourceKeeper resourcekeeper.Keeper,
) MigrationContext {
	return MigrationContext{
		didStoreKey:      didStoreKey,
		resourceStoreKey: resourceStoreKey,
		codec:            codec,
		didKeeper:        didKeeper,
		resourceKeeper:   resourceKeeper,
	}
}

type Migration func(sctx sdk.Context, mctx MigrationContext) error

type Migrator struct {
	migration Migration
	context   MigrationContext
}

func NewMigrator(
	didStoreKey *storetypes.KVStoreKey,
	resourceStoreKey *storetypes.KVStoreKey,
	codec codec.Codec,
	didKeeper didkeeper.Keeper,
	resourceKeeper resourcekeeper.Keeper,
	migration Migration,
) Migrator {
	return Migrator{
		migration: migration,
		context: NewMigrationContext(
			didStoreKey,
			resourceStoreKey,
			codec,
			didKeeper,
			resourceKeeper,
		),
	}
}

func (m *Migrator) Migrate(ctx sdk.Context) error {
	err := m.migration(ctx, m.context)
	if err != nil {
		return err
	}

	return nil
}
