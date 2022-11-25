//go:build integration

package integration

import (
	"crypto/ed25519"
	"fmt"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	helpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	cli_types "github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - negative did", func() {
	var tmpDir string

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()
	})

	It("cannot create diddoc with missing cli arguments, sign inputs mismatch, non-supported VM type, already existing did", func() {
		// Define a valid new DID Doc
		did := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
		Expect(err).To(BeNil())

		payload := types.MsgCreateDidDocPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 "Ed25519VerificationKey2020",
					Controller:           did,
					VerificationMaterial: "{\"publicKeyMultibase\":\"" + string(pubKeyMultibase58) + "\"}",
				},
			},
			Authentication: []string{keyId},
			VersionId:      uuid.NewString(),
		}

		signInputs := []cli_types.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey,
			},
		}

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Second new valid DID Doc
		did2 := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
		keyId2 := did2 + "#key1"

		pubKey2, privKey2, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase582, err := multibase.Encode(multibase.Base58BTC, pubKey2)
		Expect(err).To(BeNil())

		payload2 := types.MsgCreateDidDocPayload{
			Id: did2,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId2,
					Type:                 "Ed25519VerificationKey2020",
					Controller:           did2,
					VerificationMaterial: "{\"publicKeyMultibase\":\"" + string(pubKeyMultibase582) + "\"}",
				},
			},
			Authentication: []string{keyId2},
			VersionId:      uuid.NewString(),
		}

		signInputs2 := []cli_types.SignInput{
			{
				VerificationMethodId: keyId2,
				PrivKey:              privKey2,
			},
		}

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot create diddoc with missing cli arguments"))
		// Fail to create a new DID Doc with missing cli arguments
		//   a. missing payload, sign inputs and account
		_, err = cli.CreateDidDoc(tmpDir, types.MsgCreateDidDocPayload{}, []cli_types.SignInput{}, "", cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   b. missing payload, sign inputs
		_, err = cli.CreateDidDoc(tmpDir, types.MsgCreateDidDocPayload{}, []cli_types.SignInput{}, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   c. missing payload, account
		_, err = cli.CreateDidDoc(tmpDir, types.MsgCreateDidDocPayload{}, signInputs2, "", cli.CLI_GAS_PARAMS)

		//   d. missing sign inputs, account
		_, err = cli.CreateDidDoc(tmpDir, payload2, []cli_types.SignInput{}, "", cli.CLI_GAS_PARAMS)

		//   e. missing payload
		_, err = cli.CreateDidDoc(tmpDir, types.MsgCreateDidDocPayload{}, signInputs2, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   f. missing sign inputs
		_, err = cli.CreateDidDoc(tmpDir, payload2, []cli_types.SignInput{}, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   g. missing account
		_, err = cli.CreateDidDoc(tmpDir, payload2, signInputs2, "", cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot create diddoc with sign inputs mismatch"))
		// Fail to create a new DID Doc with sign inputs mismatch
		//   a. sign inputs mismatch
		_, err = cli.CreateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   b. non-existing key id
		_, err = cli.CreateDidDoc(tmpDir, payload2, []cli_types.SignInput{
			{
				VerificationMethodId: "non-existing-key-id",
				PrivKey:              privKey2,
			},
		}, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   c. non-matching private key
		_, err = cli.CreateDidDoc(tmpDir, payload2, []cli_types.SignInput{
			{
				VerificationMethodId: keyId2,
				PrivKey:              privKey,
			},
		}, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot create diddoc with non-supported VM type"))
		// Fail to create a new DID Doc with non-supported VM type
		payload3 := payload2
		payload3.VerificationMethod[0].Type = "NonSupportedVMType"
		_, err = cli.CreateDidDoc(tmpDir, payload3, signInputs2, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot create diddoc with already existing DID"))
		// Fail to create a new DID Doc with already existing DID
		_, err = cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())
	})

	deepCopierUpdateDid := helpers.DeepCopyUpdateDid{}

	It("cannot update a DID Doc with missing cli arguments, sign inputs mismatch, non-supported VM type, non-existing did, unchanged payload", func() {
		// Define a valid DID Doc to be updated
		did := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
		Expect(err).To(BeNil())

		payload := types.MsgCreateDidDocPayload{
			Id:         did,
			Controller: []string{did},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 "Ed25519VerificationKey2020",
					Controller:           did,
					VerificationMaterial: "{\"publicKeyMultibase\":\"" + string(pubKeyMultibase58) + "\"}",
				},
			},
			Authentication: []string{keyId},
			VersionId:      uuid.NewString(),
		}

		signInputs := []cli_types.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey,
			},
		}

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Update the DID Doc
		updatedPayload := types.MsgUpdateDidDocPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 "Ed25519VerificationKey2020",
					Controller:           did,
					VerificationMaterial: "{\"publicKeyMultibase\":\"" + string(pubKeyMultibase58) + "\"}",
				},
			},
			Authentication:  []string{keyId},
			AssertionMethod: []string{keyId},
			VersionId:       uuid.NewString(),
		}

		res, err = cli.UpdateDidDoc(tmpDir, updatedPayload, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Generate second controller
		did2 := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
		keyId2 := did2 + "#key1"
		keyId2AsExtraController := did + "#key2"

		pubKey2, privKey2, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase582, err := multibase.Encode(multibase.Base58BTC, pubKey2)
		Expect(err).To(BeNil())

		payload2 := types.MsgCreateDidDocPayload{
			Id:         did2,
			Controller: []string{did2},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId2,
					Type:                 "Ed25519VerificationKey2020",
					Controller:           did2,
					VerificationMaterial: "{\"publicKeyMultibase\":\"" + string(pubKeyMultibase582) + "\"}",
				},
			},
			Authentication: []string{keyId2},
			VersionId:      uuid.NewString(),
		}

		signInputs2 := []cli_types.SignInput{
			{
				VerificationMethodId: keyId2,
				PrivKey:              privKey2,
			},
		}

		res_, err := cli.CreateDidDoc(tmpDir, payload2, signInputs2, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
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

		signInputsFuzzed := []cli_types.SignInput{
			{
				VerificationMethodId: keyIdFuzzed,
				PrivKey:              privKeyFuzzed,
			},
			{
				VerificationMethodId: keyIdFuzzed2,
				PrivKey:              privKeyFuzzed2,
			},
		}

		// Following valid DID Doc to be updated
		followingUpdatedPayload := deepCopierUpdateDid.DeepCopy(updatedPayload)
		followingUpdatedPayload.Controller = []string{did, did2}
		followingUpdatedPayload.VerificationMethod = append(followingUpdatedPayload.VerificationMethod, &types.VerificationMethod{
			Id:                   keyId2AsExtraController,
			Type:                 "Ed25519VerificationKey2020",
			Controller:           did2,
			VerificationMaterial: "{\"publicKeyMultibase\":\"" + string(pubKeyMultibase582) + "\"}",
		})
		followingUpdatedPayload.Authentication = append(followingUpdatedPayload.Authentication, keyId2AsExtraController)
		followingUpdatedPayload.CapabilityDelegation = []string{keyId}
		followingUpdatedPayload.CapabilityInvocation = []string{keyId}
		followingUpdatedPayload.VersionId = uuid.NewString()

		signInputsAugmented := append(signInputs, signInputs2...)

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot update diddoc with missing cli arguments"))
		// Fail to update the DID Doc with missing cli arguments
		//   a. missing payload, sign inputs and account
		_, err = cli.UpdateDidDoc(tmpDir, types.MsgUpdateDidDocPayload{}, []cli_types.SignInput{}, "", cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   b. missing payload, sign inputs
		_, err = cli.UpdateDidDoc(tmpDir, types.MsgUpdateDidDocPayload{}, []cli_types.SignInput{}, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   c. missing payload, account
		_, err = cli.UpdateDidDoc(tmpDir, types.MsgUpdateDidDocPayload{}, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   d. missing sign inputs, account
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, []cli_types.SignInput{}, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   e. missing payload
		_, err = cli.UpdateDidDoc(tmpDir, types.MsgUpdateDidDocPayload{}, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   f. missing sign inputs
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, []cli_types.SignInput{}, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   g. missing account
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, signInputs, "", cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot update diddoc with sign inputs mismatch"))
		// Fail to update the DID Doc with sign inputs mismatch
		//   a. sign inputs total mismatch
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, signInputsFuzzed, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   b. sign inputs invalid length
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//   c. non-existing key id
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, []cli_types.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey,
			},
			{
				VerificationMethodId: "non-existing-key-id",
				PrivKey:              privKey2,
			},
		}, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//  d. non-matching private key
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, []cli_types.SignInput{
			{
				VerificationMethodId: keyId2AsExtraController,
				PrivKey:              privKey,
			},
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey2,
			},
		}, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		//  e. invalid private key
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, []cli_types.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKeyFuzzedExtra,
			},
			{
				VerificationMethodId: keyId2AsExtraController,
				PrivKey:              privKey2,
			},
		}, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot update diddoc with a non-supported VM type"))
		// Fail to update the DID Doc with a non-supported VM type
		invalidVmTypePayload := deepCopierUpdateDid.DeepCopy(followingUpdatedPayload)
		invalidVmTypePayload.VerificationMethod = []*types.VerificationMethod{
			followingUpdatedPayload.VerificationMethod[0],
			{
				Id:                   followingUpdatedPayload.VerificationMethod[1].Id,
				Type:                 "NonSupportedVmType",
				Controller:           followingUpdatedPayload.VerificationMethod[1].Controller,
				VerificationMaterial: "{\"publicKeyMultibase\":\"pretty-long-public-key-multibase\"}",
			},
		}
		invalidVmTypePayload.VersionId = uuid.NewString()
		_, err = cli.UpdateDidDoc(tmpDir, invalidVmTypePayload, signInputsAugmented, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot update diddoc with a non-existing DID"))
		// Fail to update a non-existing DID Doc
		nonExistingDid := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
		nonExistingDidPayload := deepCopierUpdateDid.DeepCopy(followingUpdatedPayload)
		nonExistingDidPayload.Id = nonExistingDid
		_, err = cli.UpdateDidDoc(tmpDir, nonExistingDidPayload, signInputsAugmented, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())

		// Finally, update the DID Doc
		res, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, signInputsAugmented, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot update diddoc with an unchanged payload"))
		// Fail to update the DID Doc with an unchanged payload
		followingUpdatedPayload.VersionId = uuid.NewString()
		_, err = cli.UpdateDidDoc(tmpDir, followingUpdatedPayload, signInputsAugmented, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil()) // TODO: Decide if this should be an error, if the DID Doc is unchanged
	})

	It("cannot query a diddoc with missing cli arguments, non-existing diddoc", func() {
		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot query diddoc with missing cli arguments"))
		// Fail to query the DID Doc with missing cli arguments
		//   a. missing did
		_, err := cli.QueryDidDoc("")
		Expect(err).ToNot(BeNil())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot query diddoc with a non-existing DID"))
		// Fail to query a non-existing DID Doc
		nonExistingDid := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
		_, err = cli.QueryDidDoc(nonExistingDid)
		Expect(err).ToNot(BeNil())
	})
})
