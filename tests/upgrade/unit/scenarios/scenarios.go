package scenarios

import (
	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/unit/setup"

	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RunChecksumScenario() error {
	// Init storages, keepers and setup the migration context.
	setup := migrationsetup.Setup()

	builder := NewChecksumBuilder(setup)

	dataSet, err := builder.BuildDataSet(setup)
	if err != nil {
		return err
	}

	resourceChecksumScenario := NewMigrationScenario(
		"ResourceChecksum",
		func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error {
			return appmigrations.MigrateResourceChecksumV2(ctx, migrationCtx)
		},
	)
	// Init Migrator structure
	migrator := NewMigrator(
		[]MigrationScenario{resourceChecksumScenario},
		setup,
		&dataSet)

	// Run migration scenario
	err = migrator.Run()
	return err
}

func RunProtobufScenario() error {
	// Init storages, keepers and setup the migration context.
	setup := migrationsetup.Setup()

	builder := NewProtobufBuilder(setup)

	dataSet, err := builder.BuildDataSet(setup)
	if err != nil {
		return err
	}

	resourceProtobufScenario := NewMigrationScenario(
		"ResourceChecksum",
		func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error {
			return appmigrations.MigrateDidProtobufV1(ctx, migrationCtx)
		},
	)
	// Init Migrator structure
	migrator := NewMigrator(
		[]MigrationScenario{resourceProtobufScenario},
		setup,
		&dataSet)

	// Run migration scenario
	err = migrator.Run()

	return err
}
