//go:build integration

package integration

import (
	"crypto/ed25519"
	"fmt"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	//"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	clitypes "github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive diddoc pricing", func() {
	var tmpDir string
	var payload types.MsgCreateDidDocPayload
	var signInputs []clitypes.SignInput

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		// Create a new DID Doc
		did := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
		Expect(err).To(BeNil())

		payload = types.MsgCreateDidDocPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 "Ed25519VerificationKey2020",
					Controller:           did,
					VerificationMaterial: "{\"publicKeyMultibase\": \"" + string(pubKeyMultibase58) + "\"}",
				},
			},
			Authentication: []string{keyId},
			VersionId:      uuid.NewString(),
		}

		signInputs = []clitypes.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey,
			},
		}
	})

	/* It("should tax create diddoc message - case: fixed fee", func() {
		feeParams := helpers.GenerateFees("5000000000ncheq")

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, feeParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		fmt.Println("res", res)

		// WIP:
		// 	 1. Check the balance of fee payer account
		//   2. Check supply deflation
		//   3. Check events
		//      - fee: {amount: 5000000000ncheq, payer: <fee payer>, granter: <granter>}
	}) */

	It("should tax create diddoc message - case: simulation", func() {
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_SIMULATION_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		fmt.Println("res", res)

		// WIP:
		// 	 1. Check the balance of fee payer account
		//   2. Check supply deflation
		//   3. Check events
		//      - fee: {amount: 5000000000ncheq, payer: <fee payer>, granter: <granter>}
	})
})
