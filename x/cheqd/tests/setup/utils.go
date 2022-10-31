package setup

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multibase"
	. "github.com/onsi/gomega"
)

func randBase58Seq(bytes int) string {
	b := make([]byte, bytes)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base58.Encode(b)
}

type IDType int

const (
	Base58_16bytes IDType = iota
	UUID           IDType = iota
)

func GenerateDID(idtype IDType) string {
	prefix := "did:cheqd:" + DID_NAMESPACE + ":"

	switch idtype {
	case Base58_16bytes:
		return prefix + randBase58Seq(16)
	case UUID:
		return prefix + uuid.NewString()
	default:
		panic("Unknown ID type")
	}
}

func GenerateKeyPair() KeyPair {
	PublicKey, PrivateKey, _ := ed25519.GenerateKey(rand.Reader)
	return KeyPair{PrivateKey, PublicKey}
}

func MustEncodeBase58(data []byte) string {
	encoded, err := multibase.Encode(multibase.Base58BTC, data)
	Expect(err).To(BeNil())
	return encoded
}
