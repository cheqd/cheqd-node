package migration

import (
	. "github.com/cheqd/cheqd-node/tests/upgrade/unit/scenarios"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/unit/setup"

	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ = Describe("Migration - Unit", func() {
	It("checks that Checksum migration handler works", func() {
		By("Ensuring the Checksum migration scenario is successful")

		// Init storages, keepers and setup the migration context.
		setup := migrationsetup.Setup()

		builder := NewChecksumBuilder(setup)

		dataSet, err := builder.BuildDataSet(setup)
		Expect(err).To(BeNil())

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
		Expect(err).To(BeNil())
	})

	It("checks that Protobuf migration handler works", func() {
		By("Ensuring the Protobuf migration handler is working as expected")

		// Init storages, keepers and setup the migration context.
		setup := migrationsetup.Setup()

		builder := NewProtobufBuilder(setup)

		dataSet, err := builder.BuildDataSet(setup)
		Expect(err).To(BeNil())

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
		Expect(err).To(BeNil())
	})

	It("checks IndyStyle Migration", func() {
		By("Ensuring the IndyStyle migration handler is working as expected")

		// Run IndyStyle migration
		// Init storages, keepers and setup the migration context.
		setup := migrationsetup.Setup()

		builder := NewIndyStyleBuilder(setup)

		dataSet, err := builder.BuildDataSet(setup)
		Expect(err).To(BeNil())

		indyStyleScenario := NewMigrationScenario(
			"IndyStyle Migration",
			func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error {
				return appmigrations.MigrateDidIndyStyleIdsV1(ctx, migrationCtx)
			},
		)
		// Init Migrator structure
		migrator := NewMigrator(
			[]MigrationScenario{indyStyleScenario},
			setup,
			&dataSet)

		// Run migration scenario
		err = migrator.Run()
		Expect(err).To(BeNil())
	})
})
