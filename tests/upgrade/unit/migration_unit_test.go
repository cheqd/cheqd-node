// nogo:build upgrade_unit

package unit

import (
	"path/filepath"

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
		existingDataset.MustAddDidDocV2(filepath.Join("payload", "checksum", "existing", "v2"), "diddoc")
		existingDataset.MustAddResourceV2(filepath.Join("payload", "checksum", "existing", "v2"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(filepath.Join("payload", "checksum", "expected", "v2"), "diddoc")
		expectedDataset.MustAddResourceV2(filepath.Join("payload", "checksum", "expected", "v2"), "resource")

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
		existingDataset.MustAddDidDocV1(filepath.Join("payload", "protobuf", "existing", "v1"), "diddoc")
		existingDataset.MustAddResourceV1(filepath.Join("payload", "protobuf", "existing", "v1"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(filepath.Join("payload", "protobuf", "expected", "v2"), "diddoc")
		expectedDataset.MustAddResourceV2(filepath.Join("payload", "protobuf", "expected", "v2"), "resource")

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
		existingDataset.MustAddDidDocV2(filepath.Join("payload", "indy_style", "existing", "v2"), "diddoc")
		existingDataset.MustAddResourceV2(filepath.Join("payload", "indy_style", "existing", "v2"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(filepath.Join("payload", "indy_style", "expected", "v2"), "diddoc")
		expectedDataset.MustAddResourceV2(filepath.Join("payload", "indy_style", "expected", "v2"), "resource")

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
		existingDataset.MustAddDidDocV2(filepath.Join("payload", "uuid", "existing", "v2"), "diddoc")
		existingDataset.MustAddResourceV2(filepath.Join("payload", "uuid", "existing", "v2"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(filepath.Join("payload", "uuid", "expected", "v2"), "diddoc")
		expectedDataset.MustAddResourceV2(filepath.Join("payload", "uuid", "expected", "v2"), "resource")

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

	It("checks that Version Id migration works", func() {
		By("Ensuring the Version Id migration handler is working as expected")
		// Init storages, keepers and setup the migration context.
		setup := Setup()

		// Existing dataset
		existingDataset := NewExistingDataset(setup)
		existingDataset.MustAddDidDocV2(filepath.Join("payload", "version_id", "existing", "v2"), "diddoc")
		existingDataset.MustAddResourceV2(filepath.Join("payload", "version_id", "existing", "v2"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(filepath.Join("payload", "version_id", "expected", "v2"), "diddoc")
		expectedDataset.MustAddResourceV2(filepath.Join("payload", "version_id", "expected", "v2"), "resource")

		// Migrator
		migrator := NewMigrator(
			setup,
			[]appmigrations.Migration{
				appmigrations.MigrateDidVersionId,
			},
			*existingDataset,
			*expectedDataset)

		// Run migration
		err := migrator.Run()
		Expect(err).To(BeNil())
	})

	It("checks that Resource Version Links migration works", func() {
		By("Ensuring the Resource Version Links migration handler is working as expected")
		// Init storages, keepers and setup the migration context.
		setup := Setup()

		// Existing dataset
		existingDataset := NewExistingDataset(setup)
		existingDataset.MustAddDidDocV2(filepath.Join("payload", "resource_links", "existing", "v2"), "diddoc")
		existingDataset.MustAddResourceV2(filepath.Join("payload", "resource_links", "existing", "v2"), "resource")

		// Expected dataset
		expectedDataset := NewExpectedDataset(setup)
		expectedDataset.MustAddDidDocV2(filepath.Join("payload", "resource_links", "expected", "v2"), "diddoc")
		expectedDataset.MustAddResourceV2(filepath.Join("payload", "resource_links", "expected", "v2"), "resource")

		// Migrator
		migrator := NewMigrator(
			setup,
			[]appmigrations.Migration{
				appmigrations.MigrateResourceVersionLinks,
			},
			*existingDataset,
			*expectedDataset)

		// Run migration
		err := migrator.Run()
		Expect(err).To(BeNil())
	})
})
