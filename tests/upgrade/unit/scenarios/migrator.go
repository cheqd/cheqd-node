package scenarios

import (
	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/unit/setup"
)

type Migrator struct {
	migrations []MigrationScenario
	dataSet    IDataSet
	setup      migrationsetup.TestSetup
}

func NewMigrator(
	migrations []MigrationScenario,
	setup migrationsetup.TestSetup,
	dataSet IDataSet,
) Migrator {

	return Migrator{
		migrations: migrations,
		dataSet:    dataSet,
		setup:      setup,
	}
}

func (m Migrator) Migrate() error {
	migrationCtx := appmigrations.NewMigrationContext(
		m.setup.Cdc,
		m.setup.DidStoreKey,
		m.setup.ResourceStoreKey,
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

func (m Migrator) Run() error {
	// err := m.LoadDataSet()
	// if err != nil {
	// 	return err
	// }
	err := m.Prepare()
	if err != nil {
		return err
	}
	err = m.Migrate()
	if err != nil {
		return err
	}
	err = m.Validate()
	if err != nil {
		return err
	}
	return nil
}
