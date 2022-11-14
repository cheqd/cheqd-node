package setup

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/x/did/types"
)

type KeyPair struct {
	Private ed25519.PrivateKey
	Public  ed25519.PublicKey
}

type SignInput struct {
	VerificationMethodId string
	Key                  ed25519.PrivateKey
}

type MinimalDidDocInfo struct {
	Msg       *types.MsgCreateDidDocPayload
	SignInput SignInput
}

type DidDocInfo struct {
	Msg          *types.MsgCreateDidDocPayload
	Did          string
	CollectionId string
	KeyPair      KeyPair
	KeyId        string
	SignInput    SignInput
}

type CreatedDidDocInfo struct {
	DidDocInfo
	VersionId string
}
