package setup

import (
	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
)

type Migrator struct {
	setup TestSetup

	migrations []appmigrations.Migration

	existingDataset ExistingDataset
	expectedDataset ExpectedDataset
}

func NewMigrator(
	setup TestSetup,
	migrations []appmigrations.Migration,
	existingDataset ExistingDataset,
	expectedDataset ExpectedDataset,
) Migrator {
	return Migrator{
		setup: setup,

		migrations: migrations,

		existingDataset: existingDataset,
		expectedDataset: expectedDataset,
	}
}

func (m Migrator) Run() error {
	err := m.existingDataset.FillStore()
	if err != nil {
		return err
	}

	migrationCtx := appmigrations.NewMigrationContext(
		m.setup.Cdc,
		m.setup.DidStoreKey,
		getSubspace(didtypes.ModuleName, m.setup.ParamsKeeper),
		m.setup.ResourceStoreKey,
		getSubspace(resourcetypes.ModuleName, m.setup.ParamsKeeper),
	)

	for _, migration := range m.migrations {
		err := migration(m.setup.SdkCtx, migrationCtx)
		if err != nil {
			return err
		}
	}

	err = m.expectedDataset.CheckStore()
	if err != nil {
		return err
	}

	return nil
}
