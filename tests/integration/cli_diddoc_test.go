//go:build integration

package integration

import (
	"crypto/ed25519"
	"fmt"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive did", func() {
	var tmpDir string
	var feeParams types.FeeParams

	BeforeAll(func() {
		// Query fee params
		res, err := cli.QueryParams(types.ModuleName, string(types.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = helpers.Codec.UnmarshalJSON([]byte(res.Value), &feeParams)
		Expect(err).To(BeNil())
	})

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()
	})

	It("can create diddoc, update it and query the result (Ed25519VerificationKey2020)", func() {
		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can create diddoc (Ed25519VerificationKey2020)"))
		// Create a new DID Doc
		did := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyID := did + "#key1"

		publicKey, privateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(publicKey)

		payload := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyID,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": publicKeyMultibase,
				},
			},
			Authentication: []string{keyID},
		}

		signInputs := []didcli.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              privateKey,
			},
		}

		versionID := uuid.NewString()

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can update diddoc (Ed25519VerificationKey2020)"))
		// Update the DID Doc
		newPublicKey, newPrivateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		newPublicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(newPublicKey)

		payload2 := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyID,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": newPublicKeyMultibase,
				},
			},
			Authentication: []string{keyID},
		}

		signInputs2 := []didcli.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              privateKey,
			},
			{
				VerificationMethodID: keyID,
				PrivKey:              newPrivateKey,
			},
		}

		versionID = uuid.NewString()

		res2, err := cli.UpdateDidDoc(tmpDir, payload2, signInputs2, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).To(BeNil())
		Expect(res2.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query diddoc (Ed25519VerificationKey2020)"))
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
		Expect(didDoc.VerificationMethod[0].VerificationMaterial).To(BeEquivalentTo(newPublicKeyMultibase))

		// Check that DIDDoc is not deactivated
		Expect(resp.Value.Metadata.Deactivated).To(BeFalse())

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can deactivate diddoc (Ed25519VerificationKey2020)"))
		// Deactivate the DID Doc
		payload3 := types.MsgDeactivateDidDocPayload{
			Id: did,
		}

		signInputs3 := []didcli.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              newPrivateKey,
			},
		}

		versionID = uuid.NewString()

		res3, err := cli.DeactivateDidDoc(tmpDir, payload3, signInputs3, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.DeactivateDid.String()))
		Expect(err).To(BeNil())
		Expect(res3.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query deactivated diddoc (Ed25519VerificationKey2020)"))
		// Query the DID Doc

		resp2, err := cli.QueryDidDoc(did)
		Expect(err).To(BeNil())

		didDoc2 := resp2.Value.DidDoc
		Expect(didDoc2).To(BeEquivalentTo(didDoc))

		// Check that the DID Doc is deactivated
		Expect(resp2.Value.Metadata.Deactivated).To(BeTrue())
	})

	It("can create diddoc, update it and query the result (JsonWebKey2020)", func() {
		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can create diddoc (JsonWebKey2020)"))
		// Create a new DID Doc
		did := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyID := did + "#key1"

		publicKey, privateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyJwkJSON := testsetup.GenerateJSONWebKey2020VerificationMaterial(publicKey)
		publicKeyJwk, err := testsetup.ParseJSONToMap(publicKeyJwkJSON)
		Expect(err).To(BeNil())

		payload := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":           keyID,
					"type":         "JsonWebKey2020",
					"controller":   did,
					"publicKeyJwk": publicKeyJwk,
				},
			},
			Authentication: []string{keyID},
		}

		signInputs := []didcli.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              privateKey,
			},
		}

		versionID := uuid.NewString()

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can update diddoc (JsonWebKey2020)"))
		// Update the DID Doc
		newPublicKey, newPrivateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		newPublicKeyJwkJSON := testsetup.GenerateJSONWebKey2020VerificationMaterial(newPublicKey)
		newPublicKeyJwk, err := testsetup.ParseJSONToMap(newPublicKeyJwkJSON)
		Expect(err).To(BeNil())

		payload2 := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":           keyID,
					"type":         "JsonWebKey2020",
					"controller":   did,
					"publicKeyJwk": newPublicKeyJwk,
				},
			},
			Authentication: []string{keyID},
		}

		signInputs2 := []didcli.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              privateKey,
			},
			{
				VerificationMethodID: keyID,
				PrivKey:              newPrivateKey,
			},
		}

		versionID = uuid.NewString()

		res2, err := cli.UpdateDidDoc(tmpDir, payload2, signInputs2, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).To(BeNil())
		Expect(res2.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query diddoc (JsonWebKey2020)"))
		// Query the DID Doc
		resp, err := cli.QueryDidDoc(did)
		Expect(err).To(BeNil())

		didDoc := resp.Value.DidDoc
		Expect(didDoc.Id).To(BeEquivalentTo(did))
		Expect(didDoc.Authentication).To(HaveLen(1))
		Expect(didDoc.Authentication[0]).To(BeEquivalentTo(keyID))
		Expect(didDoc.VerificationMethod).To(HaveLen(1))
		Expect(didDoc.VerificationMethod[0].Id).To(BeEquivalentTo(keyID))
		Expect(didDoc.VerificationMethod[0].VerificationMethodType).To(BeEquivalentTo("JsonWebKey2020"))
		Expect(didDoc.VerificationMethod[0].Controller).To(BeEquivalentTo(did))
		Expect(didDoc.VerificationMethod[0].VerificationMaterial).To(BeEquivalentTo(newPublicKeyJwkJSON))

		// Check that DIDDoc is not deactivated
		Expect(resp.Value.Metadata.Deactivated).To(BeFalse())

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can deactivate diddoc (JsonWebKey2020)"))
		// Deactivate the DID Doc
		payload3 := types.MsgDeactivateDidDocPayload{
			Id: did,
		}

		signInputs3 := []didcli.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              newPrivateKey,
			},
		}

		versionID = uuid.NewString()

		res3, err := cli.DeactivateDidDoc(tmpDir, payload3, signInputs3, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.DeactivateDid.String()))
		Expect(err).To(BeNil())
		Expect(res3.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query deactivated diddoc (JsonWebKey2020)"))
		// Query the DID Doc

		resp2, err := cli.QueryDidDoc(did)
		Expect(err).To(BeNil())

		didDoc2 := resp2.Value.DidDoc
		Expect(didDoc2).To(BeEquivalentTo(didDoc))

		// Check that the DID Doc is deactivated
		Expect(resp2.Value.Metadata.Deactivated).To(BeTrue())
	})

	It("can create diddoc, update it and query the result (Ed25519VerificationKey2018)", func() {
		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can create diddoc (Ed25519VerificationKey2018)"))
		// Create a new DID Doc
		did := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyID := did + "#key1"

		publicKey, privateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyBase58 := testsetup.GenerateEd25519VerificationKey2018VerificationMaterial(publicKey)

		payload := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":              keyID,
					"type":            "Ed25519VerificationKey2018",
					"controller":      did,
					"publicKeyBase58": publicKeyBase58,
				},
			},
			Authentication: []string{keyID},
		}

		signInputs := []didcli.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              privateKey,
			},
		}

		versionID := uuid.NewString()

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can update diddoc (Ed25519VerificationKey2018)"))
		// Update the DID Doc
		newPublicKey, newPrivateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		newPublicKeyBase58 := testsetup.GenerateEd25519VerificationKey2018VerificationMaterial(newPublicKey)

		payload2 := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":              keyID,
					"type":            "Ed25519VerificationKey2018",
					"controller":      did,
					"publicKeyBase58": newPublicKeyBase58,
				},
			},
			Authentication: []string{keyID},
		}

		signInputs2 := []didcli.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              privateKey,
			},
			{
				VerificationMethodID: keyID,
				PrivKey:              newPrivateKey,
			},
		}

		versionID = uuid.NewString()

		res2, err := cli.UpdateDidDoc(tmpDir, payload2, signInputs2, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).To(BeNil())
		Expect(res2.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query diddoc (Ed25519VerificationKey2018)"))
		// Query the DID Doc
		resp, err := cli.QueryDidDoc(did)
		Expect(err).To(BeNil())

		didDoc := resp.Value.DidDoc
		Expect(didDoc.Id).To(BeEquivalentTo(did))
		Expect(didDoc.Authentication).To(HaveLen(1))
		Expect(didDoc.Authentication[0]).To(BeEquivalentTo(keyID))
		Expect(didDoc.VerificationMethod).To(HaveLen(1))
		Expect(didDoc.VerificationMethod[0].Id).To(BeEquivalentTo(keyID))
		Expect(didDoc.VerificationMethod[0].VerificationMethodType).To(BeEquivalentTo("Ed25519VerificationKey2018"))
		Expect(didDoc.VerificationMethod[0].Controller).To(BeEquivalentTo(did))
		Expect(didDoc.VerificationMethod[0].VerificationMaterial).To(BeEquivalentTo(newPublicKeyBase58))

		// Check that DIDDoc is not deactivated
		Expect(resp.Value.Metadata.Deactivated).To(BeFalse())

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can deactivate diddoc (Ed25519VerificationKey2018)"))
		// Deactivate the DID Doc
		payload3 := types.MsgDeactivateDidDocPayload{
			Id: did,
		}

		signInputs3 := []didcli.SignInput{
			{
				VerificationMethodID: keyID,
				PrivKey:              newPrivateKey,
			},
		}

		versionID = uuid.NewString()

		res3, err := cli.DeactivateDidDoc(tmpDir, payload3, signInputs3, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.DeactivateDid.String()))
		Expect(err).To(BeNil())
		Expect(res3.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query deactivated diddoc (Ed25519VerificationKey2018)"))
		// Query the DID Doc

		resp2, err := cli.QueryDidDoc(did)
		Expect(err).To(BeNil())

		didDoc2 := resp2.Value.DidDoc
		Expect(didDoc2).To(BeEquivalentTo(didDoc))

		// Check that the DID Doc is deactivated
		Expect(resp2.Value.Metadata.Deactivated).To(BeTrue())
	})
})
