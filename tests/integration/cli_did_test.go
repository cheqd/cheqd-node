//nogo:build integration

package integration

import (
	"crypto/ed25519"

	test_cli "github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	"github.com/cheqd/cheqd-node/x/cheqd/client/cli"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli", func() {
	It("can create diddoc", func() {
		did := "did:cheqd:testnet:" + uuid.NewString()
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
		Expect(err).To(BeNil())

		payload := types.MsgCreateDidPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 keyId,
					Type:               "Ed25519VerificationKey2020",
					Controller:         did,
					PublicKeyMultibase: string(pubKeyMultibase58),
				},
			},
			Authentication: []string{keyId},
		}

		signInputs := []cli.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey,
			},
		}

		res, err := test_cli.CreateDid(payload, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})
})
