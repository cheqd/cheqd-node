package migration

import (
	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/migration/setup"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type IDataSet interface {
	Load() error
	Prepare() error
	Validate() error
}

type MigrationScenario struct {
	name string
	handler func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error
}

func NewMigrationScenario(
	name string,
	handler func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error,
) MigrationScenario {
	return MigrationScenario{
		name:    name,
		handler: handler,
	}
}

func (m MigrationScenario) Name() string {
	return m.name
}

func (m MigrationScenario) Handler() func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error {
	return m.handler
}

type Migrator struct {
	migrations []MigrationScenario
	dataSet    IDataSet
	setup      migrationsetup.TestSetup
}

func NewMigrator(
	migrations []MigrationScenario,
	setup migrationsetup.TestSetup,
	dataSet IDataSet) Migrator {
	return Migrator{
		migrations: migrations,
		dataSet:    dataSet,
		setup:      setup,
	}
}

func (m Migrator) Migrate() error {
	migrationCtx := appmigrations.NewMigrationContext(
		m.setup.DidStoreKey,
		m.setup.ResourceStoreKey,
		m.setup.Cdc,
		m.setup.DidKeeper,
		m.setup.ResourceKeeper)
	for _, migration := range m.migrations {
		err := migration.Handler()(m.setup.SdkCtx, migrationCtx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Migrator) Prepare() error {
	return m.dataSet.Prepare()
}

func (m Migrator) Validate() error {
	return m.dataSet.Validate()
}