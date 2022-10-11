package tests

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

type KeyPair struct {
	Private ed25519.PrivateKey
	Public  ed25519.PublicKey
}

type SignInput struct {
	VerificationMethodId string
	Key                  ed25519.PrivateKey
}

type DidInfo struct {
	Msg       *types.MsgCreateDidPayload
	Did       string
	KeyPair   KeyPair
	KeyId     string
	SignInput SignInput
}

type CreatedDidInfo struct {
	DidInfo
	VersionId string
}
