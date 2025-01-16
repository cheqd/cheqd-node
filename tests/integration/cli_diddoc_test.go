//go:build integration

package integration

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"strconv"

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

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		// Query fee params
		res, err := cli.QueryParams(types.ModuleName, string(types.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = helpers.Codec.UnmarshalJSON([]byte(res.Value), &feeParams)
		Expect(err).To(BeNil())
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

	It("can create diddoc with augmented assertionMethod, update it and query the result (Ed25519VerificationKey2020)", func() {
		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can create diddoc with augmented assertionMethod (Ed25519VerificationKey2020)"))
		// Create a new DID Doc
		did := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyID := did + "#key1"

		publicKey, privateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(publicKey)
		publicKeyBase58 := testsetup.GenerateEd25519VerificationKey2018VerificationMaterial(publicKey)

		participantId := 123
		paramsRef := "https://resolver.cheqd.net/1.0/identifiers/did:cheqd:testnet:09b20561-7339-40ea-a377-05ea35a0e82a/resources/08f35fe3-bc2a-4666-90da-972a5b05645f"
		curveType := "Bls12381BBSVerificationKeyDock2023"

		assertionMethodJSONEscaped := func() string {
			b, _ := json.Marshal(types.AssertionMethodJSONUnescaped{
				Id:              keyID,
				Type:            "Ed25519VerificationKey2018",
				Controller:      did,
				PublicKeyBase58: &publicKeyBase58, // arbitrarily chosen, loosely validated
				Metadata: 	  &types.AssertionMethodJSONUnescapedMetadata{
					ParticipantId: &participantId,
					ParamsRef:     &paramsRef,
					CurveType:     &curveType,
				},
			})
			return strconv.Quote(string(b))
		}()

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
			Authentication:  []string{keyID},
			AssertionMethod: []string{keyID, assertionMethodJSONEscaped},
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

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can update diddoc with augmented assertionMethod (Ed25519VerificationKey2020)"))
		// Update the DID Doc

		assertionMethodJSONEscaped2 := func() string {
			b, _ := json.Marshal(types.AssertionMethodJSONUnescaped{
				Id:                 did + "#key2",
				Type:               "Ed25519VerificationKey2020",
				Controller:         did,
				PublicKeyMultibase: &publicKeyMultibase, // arbitrarily chosen, loosely validated
			})
			return strconv.Quote(string(b))
		}()

		payload2 := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyID,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": publicKeyMultibase,
				},
			},
			Authentication:  []string{keyID},
			AssertionMethod: []string{keyID, assertionMethodJSONEscaped, assertionMethodJSONEscaped2},
		}

		versionID = uuid.NewString()

		res2, err := cli.UpdateDidDoc(tmpDir, payload2, signInputs, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).To(BeNil())
		Expect(res2.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query diddoc with augmented assertionMethod (Ed25519VerificationKey2020)"))
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
		Expect(didDoc.VerificationMethod[0].VerificationMaterial).To(BeEquivalentTo(publicKeyMultibase))
		Expect(didDoc.AssertionMethod).To(HaveLen(3))
		Expect(didDoc.AssertionMethod[0]).To(BeEquivalentTo(keyID))
		Expect(didDoc.AssertionMethod[1]).To(BeEquivalentTo(assertionMethodJSONEscaped))
		Expect(didDoc.AssertionMethod[2]).To(BeEquivalentTo(assertionMethodJSONEscaped2))

		// Check that DIDDoc is not deactivated
		Expect(resp.Value.Metadata.Deactivated).To(BeFalse())

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can deactivate diddoc with augmented assertionMethod (Ed25519VerificationKey2020)"))
		// Deactivate the DID Doc
		payload3 := types.MsgDeactivateDidDocPayload{
			Id: did,
		}

		versionID = uuid.NewString()

		res3, err := cli.DeactivateDidDoc(tmpDir, payload3, signInputs, versionID, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.DeactivateDid.String()))
		Expect(err).To(BeNil())
		Expect(res3.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can query deactivated diddoc with augmented assertionMethod (Ed25519VerificationKey2020)"))
		// Query the DID Doc

		resp2, err := cli.QueryDidDoc(did)
		Expect(err).To(BeNil())

		didDoc2 := resp2.Value.DidDoc
		Expect(didDoc2).To(BeEquivalentTo(didDoc))

		// Check that the DID Doc is deactivated
		Expect(resp2.Value.Metadata.Deactivated).To(BeTrue())
	})

	It("can create diddoc with empty controller, update it using the authentication key and query the result (Ed25519VerificationKey2020)", func() {
		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can create diddoc with empty controller (Ed25519VerificationKey2020)"))
		// Create a new DID Doc
		did := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyID := did + "#key1"

		publicKey, privateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(publicKey)

		keyID2 := did + "#key2"

		publicKey2, _, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase2 := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(publicKey2)

		payload := didcli.DIDDocument{
			ID: did,
			Controller: []string{},
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyID,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": publicKeyMultibase,
				},
				map[string]any{
					"id":                 keyID2,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": publicKeyMultibase2,
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

		AddReportEntry("Integration", fmt.Sprintf("%sPositive: %s", cli.Green, "can update diddoc with empty controller (Ed25519VerificationKey2020)"))
		// Update the DID Doc, removing the second verification method
		payload2 := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				payload.VerificationMethod[0],
			},
			Authentication: []string{keyID},
		}

		versionID2 := uuid.NewString()

		res2, err := cli.UpdateDidDoc(tmpDir, payload2, signInputs, versionID2, testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))

		Expect(err).To(BeNil())
		Expect(res2.Code).To(BeEquivalentTo(0))
	})
})
