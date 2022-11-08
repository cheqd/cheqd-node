package migrations

import (
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/codec"
)

type MigrationContext struct {
	codec codec.Codec

	didKeeper      didkeeper.Keeper
	resourceKeeper resourcekeeper.Keeper
}

type CheqdMigration func(ctx MigrationContext) error

type CheqdMigrator struct {
	migrations []CheqdMigration
	context    MigrationContext
}

func NewCheqdMigrator(codec codec.Codec, didKeeper didkeeper.Keeper, resourceKeeper resourcekeeper.Keeper, migrations ...CheqdMigration) CheqdMigrator {
	return CheqdMigrator{
		migrations: migrations,
		context: MigrationContext{
			didKeeper:      didKeeper,
			resourceKeeper: resourceKeeper,

			codec: codec,
		},
	}
}

func (m *CheqdMigrator) Migrate(ctx sdk.Context) error {
	for _, migration := range m.migrations {
		err := migration(m.context)
		if err != nil {
			return err
		}
	}
	return nil
}
