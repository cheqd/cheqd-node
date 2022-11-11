//go:build upgrade

package upgrade

import (
	"fmt"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"
	// migration "github.com/cheqd/cheqd-node/tests/upgrade/migration"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Post() is a function that runs after the upgrade test suite.
// Idiomatically, it is called from the upgrade_suite_test.go file, in the BeforeSuite() function.
// We will keep both AfterSuite() and Post() callback here for easiness of conceptual understanding.
var _ = AfterSuite(func() {
	/* err := Post()
	Expect(err).To(BeNil())

	err = migration.AssertMigration(&QueriedDidDoc, &QueriedResource)
	Expect(err).To(BeNil()) */

	// TODO: Add localnet volume mount cleanup & cli binary cleanup
	// This allows for a clean run of the upgrade test suite, even if the run is done locally.
})

func Post() error {
	By("Ensuring the QueryDid query is successful for the existing DID")
	res, err := cli.QueryDid(DidDoc.Id, cli.VALIDATOR0)
	Expect(err).To(BeNil())
	Expect(res.Value.DidDoc.Id).To(BeEquivalentTo(DidDoc.Id))

	By("Ensuring the QueryResource query is successful for the existing Resource")
	res_, err := cli.QueryResource(ResourcePayload.CollectionId, ResourcePayload.Id, cli.VALIDATOR0)
	Expect(err).To(BeNil())
	Expect(res_.Resource.Metadata.Id).To(BeEquivalentTo(ResourcePayload.Id))

	fmt.Printf("%sPost-Upgrade successful.\n", cli.GREEN)

	return nil
}
