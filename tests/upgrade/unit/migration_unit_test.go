//go:build upgrade_unit

package unit

import (
	. "github.com/cheqd/cheqd-node/tests/upgrade/unit/setup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
)

var _ = Describe("Migration - Unit", func() {
	It("checks that Checksum migration handler works", func() {
		By("Ensuring the Checksum migration scenario is successful")

		// Init storages, keepers and setup the migration context.
		setup := Setup()

		// Existing dataset
		existingDataset := NewExistingDataset(setup)
		existingDataset.MustAddDidDocV2(JoinGenerated("payload", "existing", "v2", "checksum"), "diddoc")
		existingDataset.MustAddResourceV2(JoinGenerated("payload", "existing", "v2", "checksum"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(JoinGenerated("payload", "expected", "v2", "checksum"), "diddoc")
		expectedDataset.MustAddResourceV2(JoinGenerated("payload", "expected", "v2", "checksum"), "resource")

		// Migrator
		migrator := NewMigrator(
			setup,
			[]appmigrations.Migration{
				appmigrations.MigrateResourceChecksum,
			},
			*existingDataset,
			*expectedDataset)

		// Run migration
		err := migrator.Run()
		Expect(err).To(BeNil())
	})

	It("checks that Protobuf migration handler works", func() {
		By("Ensuring the Protobuf migration handler is working as expected")

		// Init storages, keepers and setup the migration context.
		setup := Setup()

		// Existing dataset
		existingDataset := NewExistingDataset(setup)
		existingDataset.MustAddDidDocV1(JoinGenerated("payload", "existing", "v1", "protobuf"), "diddoc")
		existingDataset.MustAddResourceV1(JoinGenerated("payload", "existing", "v1", "protobuf"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(JoinGenerated("payload", "expected", "v2", "protobuf"), "diddoc")
		expectedDataset.MustAddResourceV2(JoinGenerated("payload", "expected", "v2", "protobuf"), "resource")

		// Migrator
		migrator := NewMigrator(
			setup,
			[]appmigrations.Migration{
				appmigrations.MigrateDidProtobuf,
				appmigrations.MigrateResourceProtobuf,
			},
			*existingDataset,
			*expectedDataset)

		// Run migration
		err := migrator.Run()
		Expect(err).To(BeNil())
	})

	It("checks IndyStyle Migration", func() {
		By("Ensuring the IndyStyle migration handler is working as expected")

		// Init storages, keepers and setup the migration context.
		setup := Setup()

		// Existing dataset
		existingDataset := NewExistingDataset(setup)
		existingDataset.MustAddDidDocV2(JoinGenerated("payload", "existing", "v2", "indy_style"), "diddoc")
		existingDataset.MustAddResourceV2(JoinGenerated("payload", "existing", "v2", "indy_style"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(JoinGenerated("payload", "expected", "v2", "indy_style"), "diddoc")
		expectedDataset.MustAddResourceV2(JoinGenerated("payload", "expected", "v2", "indy_style"), "resource")

		// Migrator
		migrator := NewMigrator(
			setup,
			[]appmigrations.Migration{
				appmigrations.MigrateDidIndyStyle,
				appmigrations.MigrateResourceIndyStyle,
			},
			*existingDataset,
			*expectedDataset)

		// Run migration
		err := migrator.Run()
		Expect(err).To(BeNil())
	})
})
