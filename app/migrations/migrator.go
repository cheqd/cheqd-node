package migrations

import (
	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migration func(sdk.Context, ...interface{}) error

type Migrator struct {
	migrations []Migration
	keeper     interface{}
}

func NewMigrator(keeper interface{}, migrations ...Migration) Migrator {
	return Migrator{
		migrations: migrations,
		keeper:     keeper,
	}
}

type CheqdMigration func(sdk.Context, cheqdkeeper.Keeper) error

type CheqdMigrator struct {
	migrations  []CheqdMigration
	cheqdKeeper cheqdkeeper.Keeper
}

func NewCheqdMigrator(cheqdKeeper cheqdkeeper.Keeper, migrations ...CheqdMigration) CheqdMigrator {
	return CheqdMigrator{
		cheqdKeeper: cheqdKeeper,
		migrations:  migrations,
	}
}

func (m *CheqdMigrator) Migrate(ctx sdk.Context) error {
	for _, migration := range m.migrations {
		err := migration(ctx, m.cheqdKeeper)
		if err != nil {
			return err
		}
	}
	return nil
}

type ResourceMigration func(sdk.Context, cheqdkeeper.Keeper, resourcekeeper.Keeper) error

type ResourceMigrator struct {
	migrations     []ResourceMigration
	cheqdKeeper    cheqdkeeper.Keeper
	resourceKeeper resourcekeeper.Keeper
}

func NewResourceMigrator(cheqdKeeper cheqdkeeper.Keeper, resourceKeeper resourcekeeper.Keeper, migrations ...ResourceMigration) ResourceMigrator {
	return ResourceMigrator{
		migrations:     migrations,
		cheqdKeeper:    cheqdKeeper,
		resourceKeeper: resourceKeeper,
	}
}

func (m *ResourceMigrator) Migrate(ctx sdk.Context) error {
	for _, migration := range m.migrations {
		err := migration(ctx, m.cheqdKeeper, m.resourceKeeper)
		if err != nil {
			return err
		}
	}
	return nil
}
