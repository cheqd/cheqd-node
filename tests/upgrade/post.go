//go:build upgrade

package upgrade

import (
	"fmt"
	"os"
	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"
	migration "github.com/cheqd/cheqd-node/tests/upgrade/migration"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Post() is a function that runs after the upgrade test suite.
// Idiomatically, it is called from the upgrade_suite_test.go file, in the BeforeSuite() function.
// We will keep both AfterSuite() and Post() callback here for easiness of conceptual understanding.
var _ = AfterSuite(func() {
	DeferCleanup(func() error {
		return os.RemoveAll(GinkgoT().TempDir())
	})
	err := Post()
	Expect(err).To(BeNil())

	err = migration.AssertMigration(&QueriedDidDoc, &QueriedResource)
	Expect(err).To(BeNil())
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
	PostErr = GenerateDidDocWithSignInputs(&PostDidDoc, &PostSignInputs)
	Expect(PostErr).To(BeNil())
	resp, err := cli.CreateDid(PostDidDoc, PostSignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(resp.Code).To(BeEquivalentTo(0))

	By("Ensuring the CreateResource Tx is successful for a new Resource")
	PostResourceErr = GenerateResource(&PostResourcePayload)
	Expect(PostResourceErr).To(BeNil())
	resp, err = cli.CreateResource(PostResourcePayload.CollectionId, PostResourcePayload.Id, PostResourcePayload.Name, PostResourcePayload.ResourceType, PostResourceFile, PostSignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(resp.Code).To(BeEquivalentTo(0))

	By("Ensuring the UpdateDid Tx is successful for a new DID")
	PostRotatedKeysErr = GenerateRotatedKeysDidDocWithSignInputs(&PostDidDoc, &PostRotatedKeysDidDoc, &PostSignInputs, &PostRotatedKeysSignInputs, resp.TxHash)
	Expect(PostRotatedKeysErr).To(BeNil())
	resp, err = cli.UpdateDid(PostRotatedKeysDidDoc, PostSignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(resp.Code).To(BeEquivalentTo(0))

	fmt.Printf("%s Post() successful.", cli.GREEN)

	return nil
}
