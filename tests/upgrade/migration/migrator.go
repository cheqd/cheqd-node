package migration

import (
	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/migration/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DataChunk struct {
	existingResources []resourcetypesv1.Resource
	existingDIDDocs   []didtypesv1.StateValue
	expectedResources []resourcetypes.Resource
	expectedDIDDocs   []didtypes.DidDocWithMetadata
}

type MigrationScenario struct {
	name string
	// setup func() migrationsetup.ExtendedTestSetup
	// dataChunk DataChunk
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
	dataChunk  DataChunk
	setup      migrationsetup.ExtendedTestSetup
}

func NewMigrator(migrations []MigrationScenario, setup migrationsetup.ExtendedTestSetup, dataChunk DataChunk) Migrator {
	return Migrator{
		migrations: migrations,
		dataChunk:  dataChunk,
		setup:      setup,
	}
}

func (m Migrator) Prepare() error {
	// for _, resource := range m.dataChunk.existingResources {
	// 	err := m.setup.ResourceKeeperV1.SetResource(m.setup.Ctx, resource)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// for _, didDoc := range m.dataChunk.existingDIDDocs {
	// 	err := m.setup.DidKeeper.SetDidDoc(m.setup.Ctx, didDoc)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (m Migrator) Migrate() error {
	migrationCtx := appmigrations.NewMigrationContext(
		m.setup.DidStoreKey,
		m.setup.ResourceStoreKey,
		m.setup.Cdc,
		m.setup.Keeper,
		m.setup.ResourceKeeper)
	for _, migration := range m.migrations {
		err := migration.Handler()(m.setup.SdkCtx, migrationCtx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Migrator) Validate() error {
	// Add check for resources and diddoc
	// Iterate over all the resources and check if they are present in the store
	// and equal to the expectedResources and expectedDIDDocs
	return nil
}
