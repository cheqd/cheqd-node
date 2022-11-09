package migration

import (
	cheqdtypes "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// AssertHandlers() is a function that tests the handlers for each module, and asserts that the handlers are working as expected.
// Idiomatically, it is called from the upgrade_suite_test.go file, in the BeforeSuite() function.
// This function is called before the Pre() function, during the upgrade test suite.
func AssertHandlers() error {
	// Here we will unit Test the migration handlers, to ensure that they are working as expected.
	// Essentially, we are matching the expected pre-generated payloads with the actual payloads that were generated during an actual migration scenario.

	By("Ensuring the ResourceChecksum migration scenario is successful")
	err := InitResourceChecksumScenario()
	Expect(err).To(BeNil())
	migrator := NewResourceMigrator([]ResourceMigrationScenario{ResourceChecksumScenario})
	err = migrator.Migrate()
	Expect(err).To(BeNil())

	// TODO: Add more migration scenarios here.

	return nil
}

// AssertMigration() is a function that runs after the upgrade test suite, and asserts that the upgrade migration was successful.
// Idiomatically, it is called from the upgrade_suite_test.go file, in the AfterSuite() function.
// This function is called after the Post() function, during the upgrade test suite.
func AssertMigration(preUpgradeDidDoc *cheqdtypes.Did, preUpgradeResource *resourcetypes.Resource) error {
	// Here we will utilize the Migrator scenario decorators per module, per handler, to assert that the migration was successful.
	// Essentially, we are matching the expected pre-generated payloads with the actual payloads that were generated during the migration.

	By("Ensuring the pre-upgrade DIDDoc is modified as expected against the post-upgrade DIDDoc")
	// TODO: Bundle Did migration scenarios and run them here.

	By("Ensuring the pre-upgrade Resource is modified as expected against the post-upgrade Resource")
	// TODO: Bundle Resource migration scenarios and run them here.

	return nil
}
