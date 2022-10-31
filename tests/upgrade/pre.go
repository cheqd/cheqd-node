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

// Pre() is a function that runs before the upgrade test suite.
// Idiomatically, it is called from the upgrade_suite_test.go file, in the BeforeSuite() function.
// We will keep both BeforeSuite() and Pre() callback here for easiness of conceptual understanding.
var _ = BeforeSuite(func() {
	err := migration.AssertHandlers()
	Expect(err).To(BeNil())

	err = Pre()
	Expect(err).To(BeNil())
})

func Pre() error {
	By("Ensuring the Err in memory is nil")
	Expect(Err).To(BeNil())

	By("Ensuring CreateDid Tx is successful")
	res, err := cli.CreateDid(DidDoc, SignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res.Code).To(BeEquivalentTo(0))

	ResourceFile, ResourceFileErr := integrationtestdata.CreateTestJson(GinkgoT().TempDir())
	By("Ensuring the ResourceFileErr in memory is nil")
	Expect(ResourceFileErr).To(BeNil())

	By("Ensuring the ResourceErr in memory is nil")
	Expect(ResourceErr).To(BeNil())

	By("Ensuring CreateResource Tx is successful")
	res, err = cli.CreateResource(ResourcePayload.CollectionId, ResourcePayload.Id, ResourcePayload.Name, ResourcePayload.ResourceType, ResourceFile, SignInputs, cli.VALIDATOR1)

	By("Ensuring the RotatedKeysErr in memory is nil")
	RotatedKeysErr = GenerateRotatedKeysDidDocWithSignInputs(&DidDoc, &RotatedKeysDidDoc, &SignInputs, &RotatedKeysSignInputs, res.TxHash)
	Expect(RotatedKeysErr).To(BeNil())

	By("Ensuring the UpdateDid Tx is successful")
	res, err = cli.UpdateDid(RotatedKeysDidDoc, RotatedKeysSignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res.Code).To(BeEquivalentTo(0))

	By("Ensuring the QueryDid query is successful")
	res_, err := cli.QueryDid(DidDoc.Id, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res_.Did.Id).To(BeEquivalentTo(RotatedKeysDidDoc.Id))
	Expect(res_.Did.Controller).To(BeEquivalentTo(RotatedKeysDidDoc.Controller))
	Expect(res_.Did.VerificationMethod).To(HaveLen(1))
	Expect(res_.Did.VerificationMethod[0].Id).To(BeEquivalentTo(RotatedKeysDidDoc.VerificationMethod[0].Id))
	Expect(res_.Did.VerificationMethod[0].Type).To(BeEquivalentTo(RotatedKeysDidDoc.VerificationMethod[0].Type))
	Expect(res_.Did.VerificationMethod[0].Controller).To(BeEquivalentTo(RotatedKeysDidDoc.VerificationMethod[0].Controller))
	Expect(res_.Did.VerificationMethod[0].PublicKeyMultibase).To(BeEquivalentTo(RotatedKeysDidDoc.VerificationMethod[0].PublicKeyMultibase))
	Expect(res_.Did.Authentication).To(HaveLen(1))
	Expect(res_.Did.Authentication[0]).To(BeEquivalentTo(RotatedKeysDidDoc.Authentication[0]))

	QueriedDidDoc = *res_.Did

	By("Ensuring the QueryResource query is successful")
	res__, err := cli.QueryResource(ResourcePayload.CollectionId, ResourcePayload.Id, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res__.Resource.Header.CollectionId).To(BeEquivalentTo(ResourcePayload.CollectionId))
	Expect(res__.Resource.Header.Id).To(BeEquivalentTo(ResourcePayload.Id))
	Expect(res__.Resource.Header.Name).To(BeEquivalentTo(ResourcePayload.Name))
	Expect(res__.Resource.Header.ResourceType).To(BeEquivalentTo(ResourcePayload.ResourceType))
	Expect(res__.Resource.Header.MediaType).To(BeEquivalentTo("application/json"))
	Expect(res__.Resource.Data).To(BeEquivalentTo(ResourcePayload.Data))

	QueriedResource = *res__.Resource

	fmt.Printf("%s Pre() successful.", cli.GREEN)

	return nil
}
