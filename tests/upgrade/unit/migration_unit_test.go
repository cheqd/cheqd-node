package unit

import (
	. "github.com/cheqd/cheqd-node/tests/upgrade/unit/scenarios"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/unit/setup"
)

var _ = Describe("Migration - Unit", func() {
	It("checks that Checksum migration handler works", func() {
		By("Ensuring the Checksum migration scenario is successful")

		// Init storages, keepers and setup the migration context.
		setup := migrationsetup.Setup()

		builder := NewChecksumBuilder(setup)

		dataSet, err := builder.BuildDataSet(setup)
		Expect(err).To(BeNil())

		resourceChecksumScenario := []appmigrations.Migration{
			appmigrations.MigrateResourceChecksum,
		}

		// Init Migrator structure
		migrator := NewMigrator(
			resourceChecksumScenario,
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

		migrations := []appmigrations.Migration{
			appmigrations.MigrateDidProtobuf,
			appmigrations.MigrateResourceProtobuf,
		}

		// Init Migrator structure
		migrator := NewMigrator(
			migrations,
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

		indyStyleScenario := []appmigrations.Migration{
			appmigrations.MigrateDidIndyStyle,
			appmigrations.MigrateResourceIndyStyle,
		}
			
		// Init Migrator structure
		migrator := NewMigrator(
			indyStyleScenario,
			setup,
			&dataSet)

		// Run migration scenario
		err = migrator.Run()
		Expect(err).To(BeNil())
	})
})
