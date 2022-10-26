//go:build upgrade
package upgrade

import (
	migration "github.com/cheqd/cheqd-node/tests/upgrade/migration"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Post() is a function that runs after the upgrade test suite.
// Idiomatically, it is called from the upgrade_suite_test.go file, in the BeforeSuite() function.
// We will keep both AfterSuite() and Post() callback here for easiness of conceptual understanding.
var _ = AfterSuite(func() {
	err := Post()
	Expect(err).To(BeNil())

	err = migration.AssertMigration()
})

func Post() error {
	// TODO: Implement this function
	Expect(true).To(BeTrue())

	return nil
}