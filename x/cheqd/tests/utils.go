package tests

import (
	"crypto/ed25519"
	"crypto/rand"
	mathrand "math/rand"

	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
	. "github.com/onsi/gomega"
)

var base58Runes = []rune("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func randBase58Seq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = base58Runes[mathrand.Intn(len(base58Runes))]
	}
	return string(b)
}

type IDType int

const (
	Base58_16chars IDType = iota
	Base58_32chars IDType = iota
	UUID           IDType = iota
)

func GenerateDID(idtype IDType) string {
	prefix := "did:cheqd:" + DID_NAMESPACE + ":"

	switch idtype {
	case Base58_16chars:
		return prefix + randBase58Seq(16)
	case Base58_32chars:
		return prefix + randBase58Seq(32)
	case UUID:
		return prefix + uuid.NewString()
	default:
		panic("Unknown ID type")
	}
}

func GenerateTestKeys() map[string]KeyPair {
	return map[string]KeyPair{
		AliceKey1:    GenerateKeyPair(),
		AliceKey2:    GenerateKeyPair(),
		BobKey1:      GenerateKeyPair(),
		BobKey2:      GenerateKeyPair(),
		BobKey3:      GenerateKeyPair(),
		BobKey4:      GenerateKeyPair(),
		CharlieKey1:  GenerateKeyPair(),
		CharlieKey2:  GenerateKeyPair(),
		CharlieKey3:  GenerateKeyPair(),
		ImposterKey1: GenerateKeyPair(),
	}
}

func InitEnv(keys map[string]KeyPair) TestSetup {
	setup := Setup()
	err := setup.CreateTestDIDs(keys)
	Expect(err).To(BeNil())
	return setup
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
