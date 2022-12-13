package migrations

import (
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didkeeperv1 "github.com/cheqd/cheqd-node/x/did/keeper/v1"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	resourcekeeperv1 "github.com/cheqd/cheqd-node/x/resource/keeper/v1"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
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

	didStoreKey  *storetypes.KVStoreKey
	didKeeperOld *didkeeperv1.Keeper
	didKeeperNew *didkeeper.Keeper

	resourceStoreKey  *storetypes.KVStoreKey
	resourceKeeperOld *resourcekeeperv1.Keeper
	resourceKeeperNew *resourcekeeper.Keeper
}

func NewMigrationContext(
	codec codec.Codec,
	didStoreKey *storetypes.KVStoreKey,
	didSubspace didtypes.ParamSubspace,
	resourceStoreKey *storetypes.KVStoreKey,
	resourceSubspace resourcetypes.ParamSubspace,
) MigrationContext {
	return MigrationContext{
		codec: codec,

		didStoreKey:  didStoreKey,
		didKeeperOld: didkeeperv1.NewKeeper(codec, didStoreKey),
		didKeeperNew: didkeeper.NewKeeper(codec, didStoreKey, didSubspace),

		resourceStoreKey:  resourceStoreKey,
		resourceKeeperOld: resourcekeeperv1.NewKeeper(codec, resourceStoreKey),
		resourceKeeperNew: resourcekeeper.NewKeeper(codec, resourceStoreKey, resourceSubspace),
	}
}
