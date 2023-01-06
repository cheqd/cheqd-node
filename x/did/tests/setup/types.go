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
	VerificationMethodID string
	Key                  ed25519.PrivateKey
}

type DidDocInfo struct {
	Msg          *types.MsgCreateDidDocPayload
	Did          string
	CollectionID string
	KeyPair      KeyPair
	KeyID        string
	SignInput    SignInput
}

type CreatedDidDocInfo struct {
	DidDocInfo
	VersionID string
}
