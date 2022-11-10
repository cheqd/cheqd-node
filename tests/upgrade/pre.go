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

	// TODO: Add localnet volume mount cleanup here.
	// This allows for a clean start of the localnet containers.

	err = Pre()
	Expect(err).To(BeNil())
})

func Pre() error {
	By("Ensuring the Err in memory is nil")
	DidDoc, SignInputs, Err = GenerateDidDocWithSignInputs()
	Expect(Err).To(BeNil())

	By("Ensuring CreateDid Tx is successful")
	res, err := cli.CreateDid(DidDoc, SignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res.Code).To(BeEquivalentTo(0))

	By("Ensuring the ResourceFileErr in memory is nil")
	ResourceFile, ResourceFileErr := integrationtestdata.CreateTestJson(GinkgoT().TempDir())
	Expect(ResourceFileErr).To(BeNil())

	By("Ensuring the ResourceErr in memory is nil")
	ResourcePayload, ResourceErr = GenerateResource(DidDoc)
	Expect(ResourceErr).To(BeNil())

	By("Ensuring the ResourceFile is copied to the localnet container")
	_, err = cli.LocalnetExecCopyAbsoluteWithPermissions(ResourceFile, cli.DOCKER_HOME, cli.VALIDATOR1)
	Expect(err).To(BeNil())

	By("Ensuring CreateResource Tx is successful")
	res, err = cli.CreateResource(ResourcePayload.CollectionId, ResourcePayload.Id, ResourcePayload.Name, ResourcePayload.ResourceType, integrationtestdata.JSON_FILE_NAME, SignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res.Code).To(BeEquivalentTo(0))

	By("Ensuring the QueryResource query is successful")
	res__, err := cli.QueryResource(ResourcePayload.CollectionId, ResourcePayload.Id, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res__.Resource.Metadata.CollectionId).To(BeEquivalentTo(ResourcePayload.CollectionId))
	Expect(res__.Resource.Metadata.Id).To(BeEquivalentTo(ResourcePayload.Id))
	Expect(res__.Resource.Metadata.Name).To(BeEquivalentTo(ResourcePayload.Name))
	Expect(res__.Resource.Metadata.ResourceType).To(BeEquivalentTo(ResourcePayload.ResourceType))
	Expect(res__.Resource.Metadata.MediaType).To(BeEquivalentTo("application/json"))
	Expect(res__.Resource.Resource.Data).To(BeEquivalentTo(ResourcePayload.Data))

	QueriedResource = *res__.Resource.Resource

	By("Ensuring the QueryDid query is successful")
	res_, err := cli.QueryDid(DidDoc.Id, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res_.Value.DidDoc.Id).To(BeEquivalentTo(DidDoc.Id))
	Expect(res_.Value.DidDoc.Controller).To(BeEquivalentTo(DidDoc.Controller))
	Expect(res_.Value.DidDoc.VerificationMethod).To(HaveLen(1))
	Expect(res_.Value.DidDoc.VerificationMethod[0].Id).To(BeEquivalentTo(DidDoc.VerificationMethod[0].Id))
	Expect(res_.Value.DidDoc.VerificationMethod[0].Type).To(BeEquivalentTo(DidDoc.VerificationMethod[0].Type))
	Expect(res_.Value.DidDoc.VerificationMethod[0].Controller).To(BeEquivalentTo(DidDoc.VerificationMethod[0].Controller))
	Expect(res_.Value.DidDoc.VerificationMethod[0].VerificationMaterial).To(BeEquivalentTo(DidDoc.VerificationMethod[0].VerificationMaterial))
	Expect(res_.Value.DidDoc.Authentication).To(HaveLen(1))
	Expect(res_.Value.DidDoc.Authentication[0]).To(BeEquivalentTo(DidDoc.Authentication[0]))

	By("Ensuring the RotatedKeysErr in memory is nil")
	RotatedKeysDidDoc, RotatedKeysSignInputs, RotatedKeysErr = GenerateRotatedKeysDidDocWithSignInputs(DidDoc, SignInputs, res_.Value.Metadata.VersionId)
	Expect(RotatedKeysErr).To(BeNil())

	By("Ensuring the UpdateDid Tx is successful")
	res, err = cli.UpdateDid(RotatedKeysDidDoc, RotatedKeysSignInputs, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res.Code).To(BeEquivalentTo(0))

	By("Ensuring the QueryDid query is successful for the updated DIDDoc")
	res_, err = cli.QueryDid(DidDoc.Id, cli.VALIDATOR1)
	Expect(err).To(BeNil())
	Expect(res_.Value.DidDoc.Id).To(BeEquivalentTo(RotatedKeysDidDoc.Id))
	Expect(res_.Value.DidDoc.Controller).To(BeEquivalentTo(RotatedKeysDidDoc.Controller))
	Expect(res_.Value.DidDoc.VerificationMethod).To(HaveLen(1))
	Expect(res_.Value.DidDoc.VerificationMethod[0].Id).To(BeEquivalentTo(RotatedKeysDidDoc.VerificationMethod[0].Id))
	Expect(res_.Value.DidDoc.VerificationMethod[0].Type).To(BeEquivalentTo(RotatedKeysDidDoc.VerificationMethod[0].Type))
	Expect(res_.Value.DidDoc.VerificationMethod[0].Controller).To(BeEquivalentTo(RotatedKeysDidDoc.VerificationMethod[0].Controller))
	Expect(res_.Value.DidDoc.VerificationMethod[0].VerificationMaterial).To(BeEquivalentTo(RotatedKeysDidDoc.VerificationMethod[0].VerificationMaterial))
	Expect(res_.Value.DidDoc.Authentication).To(HaveLen(1))
	Expect(res_.Value.DidDoc.Authentication[0]).To(BeEquivalentTo(RotatedKeysDidDoc.Authentication[0]))

	QueriedDidDoc = *res_.Value.DidDoc

	fmt.Printf("%sPre-Upgrade successful.\n", cli.GREEN)

	return nil
}
