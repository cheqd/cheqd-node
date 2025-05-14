//go:build integration

package integration

import (
	"crypto/ed25519"
	"fmt"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	helpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - negative did", func() {
	var tmpDir string
	var feeParams didtypes.FeeParams

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		// Query fee params
		res, err := cli.QueryDidParams()
		Expect(err).To(BeNil())

		feeParams = res.Params
	})

	It("cannot create diddoc with missing cli arguments, sign inputs mismatch, non-supported VM type, already existing did", func() {
		// Define a valid new DID Doc
		did := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(pubKey)

		payload := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyId,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": publicKeyMultibase,
				},
			},
			Authentication: []string{keyId},
		}

		signInputs := []didcli.SignInput{
			{
				VerificationMethodID: keyId,
				PrivKey:              privKey,
			},
		}

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Second new valid DID Doc
		did2 := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyId2 := did2 + "#key1"

		publicKey2, privateKey2, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase2 := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(publicKey2)

		payload2 := didcli.DIDDocument{
			ID: did2,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyId2,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did2,
					"publicKeyMultibase": publicKeyMultibase2,
				},
			},
			Authentication: []string{keyId2},
		}

		signInputs2 := []didcli.SignInput{
			{
				VerificationMethodID: keyId2,
				PrivKey:              privateKey2,
			},
		}

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot create diddoc with missing cli arguments"))
		// Fail to create a new DID Doc with missing cli arguments
		//   a. missing payload, sign inputs and account
		_, err = cli.CreateDidDoc(tmpDir, didcli.DIDDocument{}, []didcli.SignInput{}, "", "", helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		//   b. missing payload, sign inputs
		_, err = cli.CreateDidDoc(tmpDir, didcli.DIDDocument{}, []didcli.SignInput{}, "", testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		//   c. missing payload, account
		_, err = cli.CreateDidDoc(tmpDir, didcli.DIDDocument{}, signInputs2, "", "", helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		//   d. missing sign inputs, account
		_, err = cli.CreateDidDoc(tmpDir, payload2, []didcli.SignInput{}, "", "", helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		//   e. missing payload
		_, err = cli.CreateDidDoc(tmpDir, didcli.DIDDocument{}, signInputs2, "", testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		//   f. missing sign inputs
		_, err = cli.CreateDidDoc(tmpDir, payload2, []didcli.SignInput{}, "", testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		//   g. missing account
		_, err = cli.CreateDidDoc(tmpDir, payload2, signInputs2, "", "", helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot create diddoc with sign inputs mismatch"))
		// Fail to create a new DID Doc with sign inputs mismatch
		//   a. sign inputs mismatch
		_, err = cli.CreateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		//   b. non-existing key id
		_, err = cli.CreateDidDoc(tmpDir, payload2, []didcli.SignInput{
			{
				VerificationMethodID: "non-existing-key-id",
				PrivKey:              privateKey2,
			},
		}, "", testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		//   c. non-matching private key
		_, err = cli.CreateDidDoc(tmpDir, payload2, []didcli.SignInput{
			{
				VerificationMethodID: keyId2,
				PrivKey:              privKey,
			},
		}, "", testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot create diddoc with non-supported VM type"))
		// Fail to create a new DID Doc with non-supported VM type
		payload3 := payload2
		payload3.VerificationMethod[0]["type"] = "NonSupportedVMType"
		_, err = cli.CreateDidDoc(tmpDir, payload3, signInputs2, "", testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot create diddoc with already existing DID"))
		// Fail to create a new DID Doc with already existing DID
		_, err = cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).ToNot(BeNil())
	})

	deepCopierUpdateDid := helpers.DeepCopyDIDDocument{}

	It("cannot update a DID Doc with missing cli arguments, sign inputs mismatch, non-supported VM type, non-existing did, unchanged payload", func() {
		// Define a valid DID Doc to be updated
		did := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(pubKey)

		payload := didcli.DIDDocument{
			ID:         did,
			Controller: []string{did},
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyId,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": publicKeyMultibase,
				},
			},
			Authentication: []string{keyId},
		}

		signInputs := []didcli.SignInput{
			{
				VerificationMethodID: keyId,
				PrivKey:              privKey,
			},
		}

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Update the DID Doc
		updatedPayload := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyId,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": publicKeyMultibase,
				},
			},
			Authentication:  []string{keyId},
			AssertionMethod: []string{keyId},
		}

		res, err = cli.UpdateDidDoc(tmpDir, updatedPayload, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Generate second controller
		did2 := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyId2 := did2 + "#key1"
		keyId2AsExtraController := did + "#key2"

		publicKey2, privateKey2, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase2 := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(publicKey2)

		payload2 := didcli.DIDDocument{
			ID:         did2,
			Controller: []string{did2},
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyId2,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did2,
					"publicKeyMultibase": publicKeyMultibase2,
				},
			},
			Authentication: []string{keyId2},
		}

		signInputs2 := []didcli.SignInput{
			{
				VerificationMethodID: keyId2,
				PrivKey:              privateKey2,
			},
		}

		res_, err := cli.CreateDidDoc(tmpDir, payload2, signInputs2, "", testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).To(BeNil())
		Expect(res_.Code).To(BeEquivalentTo(0))

		// Extra fuzzed sign inputs
		//   a. first sign input
		//	   i. fresh keys
		keyIdFuzzed := did + "#key1-fuzzed"
		_, privKeyFuzzed, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		//	   ii. fuzzed private key, invalid and non-matching
		privKeyFuzzedExtra := testdata.GenerateByteEntropy()
		Expect(len(privKeyFuzzedExtra)).NotTo(BeEquivalentTo(len(privKeyFuzzed)))

		//   b. second sign input
		//	   i. fresh keys
		keyIdFuzzed2 := did + "#key2-fuzzed"
		_, privKeyFuzzed2, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		//	   ii. fuzzed private key, invalid and non-matching
		privKeyFuzzedExtra2 := testdata.GenerateByteEntropy()
		Expect(len(privKeyFuzzedExtra2)).NotTo(BeEquivalentTo(len(privKeyFuzzed2)))

		signInputsFuzzed := []didcli.SignInput{
			{
				VerificationMethodID: keyIdFuzzed,
				PrivKey:              privKeyFuzzed,
			},
			{
				VerificationMethodID: keyIdFuzzed2,
				PrivKey:              privKeyFuzzed2,
			},
		}

		// Following valid DID Doc to be updated
		followingUpdatedPayload := deepCopierUpdateDid.DeepCopy(updatedPayload)
		followingUpdatedPayload.Controller = []string{did, did2}
		followingUpdatedPayload.VerificationMethod = append(followingUpdatedPayload.VerificationMethod, didcli.VerificationMethod{
			"id":                 keyId2AsExtraController,
			"type":               "Ed25519VerificationKey2020",
			"controller":         did2,
			"publicKeyMultibase": publicKeyMultibase2,
		})
		followingUpdatedPayload.Authentication = append(followingUpdatedPayload.Authentication, keyId2AsExtraController)
		followingUpdatedPayload.CapabilityDelegation = []string{keyId}
		followingUpdatedPayload.CapabilityInvocation = []string{keyId}

		signInputsAugmented := append(signInputs, signInputs2...)

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot update diddoc with missing cli arguments"))
		// Fail to update the DID Doc with missing cli arguments
		//   a. missing payload, sign inputs and account
		_, err = cli.UpdateDidDoc(tmpDir, didcli.DIDDocument{}, []didcli.SignInput{}, "", "", helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		//   b. missing payload, sign inputs
		_, err = cli.UpdateDidDoc(tmpDir, didcli.DIDDocument{}, []didcli.SignInput{}, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		//   c. missing payload, account
		_, err = cli.UpdateDidDoc(tmpDir, didcli.DIDDocument{}, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		//   d. missing sign inputs, account
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, []didcli.SignInput{}, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		//   e. missing payload
		_, err = cli.UpdateDidDoc(tmpDir, didcli.DIDDocument{}, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		//   f. missing sign inputs
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, []didcli.SignInput{}, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		//   g. missing account
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, signInputs, "", "", helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot update diddoc with sign inputs mismatch"))
		// Fail to update the DID Doc with sign inputs mismatch
		//   a. sign inputs total mismatch
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, signInputsFuzzed, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		//   b. sign inputs invalid length
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		//   c. non-existing key id
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, []didcli.SignInput{
			{
				VerificationMethodID: keyId,
				PrivKey:              privKey,
			},
			{
				VerificationMethodID: "non-existing-key-id",
				PrivKey:              privateKey2,
			},
		}, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		//  d. non-matching private key
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, []didcli.SignInput{
			{
				VerificationMethodID: keyId2AsExtraController,
				PrivKey:              privKey,
			},
			{
				VerificationMethodID: keyId,
				PrivKey:              privateKey2,
			},
		}, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		//  e. invalid private key
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, []didcli.SignInput{
			{
				VerificationMethodID: keyId,
				PrivKey:              privKeyFuzzedExtra,
			},
			{
				VerificationMethodID: keyId2AsExtraController,
				PrivKey:              privateKey2,
			},
		}, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot update diddoc with a non-supported VM type"))
		// Fail to update the DID Doc with a non-supported VM type
		invalidVmTypePayload := deepCopierUpdateDid.DeepCopy(followingUpdatedPayload)
		invalidVmTypePayload.VerificationMethod = []didcli.VerificationMethod{
			followingUpdatedPayload.VerificationMethod[0],
			map[string]any{
				"Id":                     followingUpdatedPayload.VerificationMethod[1]["id"],
				"VerificationMethodType": "NonSupportedVmType",
				"Controller":             followingUpdatedPayload.VerificationMethod[1]["controller"],
				"VerificationMaterial":   "pretty-long-public-key-multibase",
			},
		}
		_, err = cli.UpdateDidDoc(tmpDir, invalidVmTypePayload, signInputsAugmented, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot update diddoc with a non-existing DID"))
		// Fail to update a non-existing DID Doc
		nonExistingDid := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		nonExistingDidPayload := deepCopierUpdateDid.DeepCopy(followingUpdatedPayload)
		nonExistingDidPayload.ID = nonExistingDid
		_, err = cli.UpdateDidDoc(tmpDir, nonExistingDidPayload, signInputsAugmented, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).ToNot(BeNil())

		// Finally, update the DID Doc
		res, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, signInputsAugmented, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot update diddoc with an unchanged payload"))
		// Fail to update the DID Doc with an unchanged payload
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, signInputsAugmented, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(feeParams.UpdateDid.String()))
		Expect(err).To(BeNil()) // TODO: Decide if this should be an error, if the DID Doc is unchanged
	})

	It("cannot query a diddoc with missing cli arguments, non-existing diddoc", func() {
		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot query diddoc with missing cli arguments"))
		// Fail to query the DID Doc with missing cli arguments
		//   a. missing did
		_, err := cli.QueryDidDoc("")
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.Purple, "cannot query diddoc with a non-existing DID"))
		// Fail to query a non-existing DID Doc
		nonExistingDid := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		_, err = cli.QueryDidDoc(nonExistingDid)
		Expect(err).ToNot(BeNil())
	})
})
