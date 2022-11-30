

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
		existingDataset.MustAddDidDocV2(JoinGenerated("payload", "checksum", "existing", "v2"), "diddoc")
		existingDataset.MustAddResourceV2(JoinGenerated("payload", "checksum", "existing", "v2"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(JoinGenerated("payload", "checksum", "expected", "v2"), "diddoc")
		expectedDataset.MustAddResourceV2(JoinGenerated("payload", "checksum", "expected", "v2"), "resource")

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
		existingDataset.MustAddDidDocV1(JoinGenerated("payload", "protobuf", "existing", "v1"), "diddoc")
		existingDataset.MustAddResourceV1(JoinGenerated("payload", "protobuf", "existing", "v1"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(JoinGenerated("payload", "protobuf", "expected", "v2"), "diddoc")
		expectedDataset.MustAddResourceV2(JoinGenerated("payload", "protobuf", "expected", "v2"), "resource")

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
		existingDataset.MustAddDidDocV2(JoinGenerated("payload", "indy_style", "existing", "v2"), "diddoc")
		existingDataset.MustAddResourceV2(JoinGenerated("payload", "indy_style", "existing", "v2"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(JoinGenerated("payload", "indy_style", "expected", "v2"), "diddoc")
		expectedDataset.MustAddResourceV2(JoinGenerated("payload", "indy_style", "expected", "v2"), "resource")

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

	It("checks that UUID migration works", func() {
		By("Ensuring the UUID migration handler is working as expected")
		// Init storages, keepers and setup the migration context.
		setup := Setup()

		// Existing dataset
		existingDataset := NewExistingDataset(setup)
		existingDataset.MustAddDidDocV2(JoinGenerated("payload", "uuid", "existing", "v2"), "diddoc")
		existingDataset.MustAddResourceV2(JoinGenerated("payload", "uuid", "existing", "v2"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(JoinGenerated("payload", "uuid", "expected", "v2"), "diddoc")
		expectedDataset.MustAddResourceV2(JoinGenerated("payload", "uuid", "expected", "v2"), "resource")

		// Migrator
		migrator := NewMigrator(
			setup,
			[]appmigrations.Migration{
				appmigrations.MigrateDidUUID,
				appmigrations.MigrateResourceUUID,
			},
			*existingDataset,
			*expectedDataset)

		// Run migration
		err := migrator.Run()
		Expect(err).To(BeNil())
	})
})
