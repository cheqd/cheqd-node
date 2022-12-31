//go:build integration

package integration

import (
	"crypto/ed25519"
	"fmt"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	cli_types "github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive did", func() {
	var tmpDir string

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()
	})

	It("can create diddoc, update it and query the result", func() {
		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can create diddoc"))
		// Create a new DID Doc
		did := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyID := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase58 := testsetup.BuildEd25519VerificationKey2020VerificationMaterial(pubKey)

		payload := types.MsgCreateDidDocPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyID,
					VerificationMethodType: "Ed25519VerificationKey2020",
					Controller:             did,
					VerificationMaterial:   pubKeyMultibase58,
				},
			},
			Authentication: []string{keyID},
			VersionId:      uuid.NewString(),
		}

		signInputs := []cli_types.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              privKey,
			},
		}

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_1, cli.CLIGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can update diddoc"))
		// Update the DID Doc
		newPubKey, newPrivKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		newPubKeyMultibase58 := testsetup.BuildEd25519VerificationKey2020VerificationMaterial(newPubKey)

		payload2 := types.MsgUpdateDidDocPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                     keyID,
					VerificationMethodType: "Ed25519VerificationKey2020",
					Controller:             did,
					VerificationMaterial:   newPubKeyMultibase58,
				},
			},
			Authentication: []string{keyID},
			VersionId:      uuid.NewString(),
		}

		signInputs2 := []cli_types.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              privKey,
			},
			{
				VerificationMethodID: keyID,
				PrivKey:              newPrivKey,
			},
		}

		res2, err := cli.UpdateDidDoc(tmpDir, payload2, signInputs2, testdata.BASE_ACCOUNT_1, cli.CLIGasParams)
		Expect(err).To(BeNil())
		Expect(res2.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query diddoc"))
		// Query the DID Doc
		resp, err := cli.QueryDidDoc(did)
		Expect(err).To(BeNil())

		didDoc := resp.Value.DidDoc
		Expect(didDoc.Id).To(BeEquivalentTo(did))
		Expect(didDoc.Authentication).To(HaveLen(1))
		Expect(didDoc.Authentication[0]).To(BeEquivalentTo(keyID))
		Expect(didDoc.VerificationMethod).To(HaveLen(1))
		Expect(didDoc.VerificationMethod[0].Id).To(BeEquivalentTo(keyID))
		Expect(didDoc.VerificationMethod[0].VerificationMethodType).To(BeEquivalentTo("Ed25519VerificationKey2020"))
		Expect(didDoc.VerificationMethod[0].Controller).To(BeEquivalentTo(did))
		Expect(didDoc.VerificationMethod[0].VerificationMaterial).To(BeEquivalentTo("{\"publicKeyMultibase\": \"" + newPubKeyMultibase58 + "\"}"))

		// Check that DIDDoc is not deactivated
		Expect(resp.Value.Metadata.Deactivated).To(BeFalse())

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can deactivate diddoc"))
		// Deactivate the DID Doc
		payload3 := types.MsgDeactivateDidDocPayload{
			Id:        did,
			VersionId: uuid.NewString(),
		}

		signInputs3 := []cli_types.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              newPrivKey,
			},
		}

		res3, err := cli.DeactivateDidDoc(tmpDir, payload3, signInputs3, testdata.BASE_ACCOUNT_1, cli.CLIGasParams)
		Expect(err).To(BeNil())
		Expect(res3.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query deactivated diddoc"))
		// Query the DID Doc

		resp2, err := cli.QueryDidDoc(did)
		Expect(err).To(BeNil())

		didDoc2 := resp2.Value.DidDoc
		Expect(didDoc2).To(BeEquivalentTo(didDoc))

		// Check that the DID Doc is deactivated
		Expect(resp2.Value.Metadata.Deactivated).To(BeTrue())
	})
})
