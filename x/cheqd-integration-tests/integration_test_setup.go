package cheqd_integration_tests

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/client/cli"
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
	"github.com/spf13/cobra"
)

type KeyPair struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

type TestSetup struct {
	txCmd    *cobra.Command
	queryCmd *cobra.Command
	keys     map[string]KeyPair
}

func Setup() (*TestSetup, error) {
	setup := TestSetup{
		txCmd:    cli.GetTxCmd(),
		queryCmd: cli.GetQueryCmd(),
	}

	keys, err := setup.CreatePreparedDID()
	if err != nil {
		return nil, err
	}

	return &TestSetup{
		txCmd:    cli.GetTxCmd(),
		queryCmd: cli.GetQueryCmd(),
		keys:     keys,
	}, nil
}

func WrapRequestCreateDid(payload *v1.MsgCreateDidPayload, keys map[string]ed25519.PrivateKey) *v1.MsgCreateDid {
	result := v1.MsgCreateDid{
		Payload: payload,
	}

	var signatures []*v1.SignInfo
	signingInput := result.Payload.GetSignBytes()

	for privKeyId, privKey := range keys {
		signature := base64.StdEncoding.EncodeToString(ed25519.Sign(privKey, signingInput))
		signatures = append(signatures, &v1.SignInfo{
			VerificationMethodId: privKeyId,
			Signature:            signature,
		})
	}

	return &v1.MsgCreateDid{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (t TestSetup) CreatePreparedDID() (map[string]KeyPair, error) {
	prefilledDids := []struct {
		keys    map[string]KeyPair
		signers []string
		msg     *v1.MsgCreateDidPayload
	}{
		{
			keys: map[string]KeyPair{
				AliceKey1: GenerateKeyPair(),
				AliceKey2: GenerateKeyPair(),
			},
			signers: []string{AliceKey1},
			msg: &v1.MsgCreateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey1},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			keys: map[string]KeyPair{
				BobKey1: GenerateKeyPair(),
				BobKey2: GenerateKeyPair(),
				BobKey3: GenerateKeyPair(),
				BobKey4: GenerateKeyPair(),
			},
			signers: []string{BobKey2},
			msg: &v1.MsgCreateDidPayload{
				Id: BobDID,
				Authentication: []string{
					BobKey1,
					BobKey2,
					BobKey3,
				},
				CapabilityDelegation: []string{
					BobKey4,
				},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         BobKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey3,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey4,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
		{
			keys: map[string]KeyPair{
				CharlieKey1: GenerateKeyPair(),
				CharlieKey2: GenerateKeyPair(),
				CharlieKey3: GenerateKeyPair(),
			},
			signers: []string{CharlieKey2},
			msg: &v1.MsgCreateDidPayload{
				Id: CharlieDID,
				Authentication: []string{
					CharlieKey1,
					CharlieKey2,
					CharlieKey3,
				},
				VerificationMethod: []*v1.VerificationMethod{
					{
						Id:         CharlieKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         CharlieKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         CharlieKey3,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
	}

	keys := map[string]KeyPair{}

	for _, prefilled := range prefilledDids {
		msg := prefilled.msg

		for _, vm := range msg.VerificationMethod {
			vm.PublicKeyMultibase = "z" + base58.Encode(prefilled.keys[vm.Id].PublicKey)
		}

		signerKeys := map[string]ed25519.PrivateKey{}
		for _, signer := range prefilled.signers {
			signerKeys[signer] = prefilled.keys[signer].PrivateKey
		}

		for keyId, key := range prefilled.keys {
			keys[keyId] = key
		}

		_, err := t.SendCreateDid(msg, signerKeys)
		if err != nil {
			return nil, err
		}
	}

	return keys, nil
}

func (t TestSetup) SendCreateDid(msg *v1.MsgCreateDidPayload, keys map[string]ed25519.PrivateKey) (string, error) {
	msgWriteRequestBytes, _ := WrapRequestCreateDid(msg, keys).Marshal()
	argWriteRequest := base64.StdEncoding.EncodeToString(msgWriteRequestBytes)
	return t.ExecuteCommand(t.txCmd, "create-did", argWriteRequest, "--from=cheqd1rnr5jrt4exl0samwj0yegv99jeskl0hsxmcz96")
}

func (t TestSetup) ExecuteCommand(cmd *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err = cmd.Execute()
	return buf.String(), err
}

func GenerateKeyPair() KeyPair {
	PublicKey, PrivateKey, _ := ed25519.GenerateKey(rand.Reader)
	return KeyPair{PrivateKey, PublicKey}
}
