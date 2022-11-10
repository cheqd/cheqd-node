package migrations

import (
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
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

type DidMigration func(sdk.Context, didkeeper.Keeper) error

type DidMigrator struct {
	migrations  []DidMigration
	cheqdKeeper didkeeper.Keeper
}

func NewDidMigrator(cheqdKeeper didkeeper.Keeper, migrations ...DidMigration) DidMigrator {
	return DidMigrator{
		cheqdKeeper: cheqdKeeper,
		migrations:  migrations,
	}
}

func (m *DidMigrator) Migrate(ctx sdk.Context) error {
	for _, migration := range m.migrations {
		err := migration(ctx, m.cheqdKeeper)
		if err != nil {
			return err
		}
	}
	return nil
}

type ResourceMigration func(sdk.Context, didkeeper.Keeper, resourcekeeper.Keeper) error

type ResourceMigrator struct {
	migrations     []ResourceMigration
	cheqdKeeper    didkeeper.Keeper
	resourceKeeper resourcekeeper.Keeper
}

func NewResourceMigrator(cheqdKeeper didkeeper.Keeper, resourceKeeper resourcekeeper.Keeper, migrations ...ResourceMigration) ResourceMigrator {
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
