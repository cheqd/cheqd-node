//go:build upgrade

package upgrade

import (
	"fmt"

	integrationtestdata "github.com/cheqd/cheqd-node/tests/integration/testdata"
	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"
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

	err = migration.AssertMigration(&QueriedDidDoc, &QueriedResource)
	Expect(err).To(BeNil())

	// TODO: Add localnet volume mount cleanup & cli binary cleanup
	// This allows for a clean run of the upgrade test suite, even if the run is done locally.
})

func Post() error {
	By("Ensuring the QueryDid query is successful for the existing DID")
	res, err := cli.QueryDid(DidDoc.Id, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res.Did.Id).To(BeEquivalentTo(DidDoc.Id))

	By("Ensuring the QueryResource query is successful for the existing Resource")
	res_, err := cli.QueryResource(ResourcePayload.CollectionId, ResourcePayload.Id, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res_.Resource.Header.Id).To(BeEquivalentTo(ResourcePayload.Id))

	By("Ensuring the CreateDid Tx is successful for a new DID")
	PostDidDoc, PostSignInputs, PostErr = GenerateDidDocWithSignInputs()
	Expect(PostErr).To(BeNil())
	resp, err := cli.CreateDid(PostDidDoc, PostSignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(resp.Code).To(BeEquivalentTo(0))

	By("Ensuring the CreateResource Tx is successful for a new Resource")
	PostResourcePayload, PostResourceErr = GenerateResource(PostDidDoc)
	Expect(PostResourceErr).To(BeNil())

	By("Ensuring the PostResourceFileErr in memory is nil")
	PostResourceFile, PostResourceFileErr = integrationtestdata.CreateTestJson(GinkgoT().TempDir())
	Expect(PostResourceFileErr).To(BeNil())

	By("Ensuring the PostResourceFile is copied to the localnet container")
	_, err = cli.LocalnetExecCopyAbsoluteWithPermissions(PostResourceFile, cli.DOCKER_HOME, cli.VALIDATOR1)
	Expect(err).To(BeNil())

	By("Ensuring CreateResource Tx is successful")
	resp, err = cli.CreateResource(PostResourcePayload.CollectionId, PostResourcePayload.Id, PostResourcePayload.Name, PostResourcePayload.ResourceType, integrationtestdata.JSON_FILE_NAME, PostSignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(resp.Code).To(BeEquivalentTo(0))

	By("Ensuring the UpdateDid Tx is successful for a new DID")
	PostRotatedKeysDidDoc, PostRotatedKeysSignInputs, PostRotatedKeysErr = GenerateRotatedKeysDidDocWithSignInputs(PostDidDoc, PostSignInputs, resp.TxHash)
	Expect(PostRotatedKeysErr).To(BeNil())
	resp, err = cli.UpdateDid(PostRotatedKeysDidDoc, PostSignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(resp.Code).To(BeEquivalentTo(0))

	fmt.Printf("%s Post() successful.", cli.GREEN)

	return nil
}
