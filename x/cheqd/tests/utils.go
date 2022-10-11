package tests

import (
	"math/rand"

	. "github.com/onsi/gomega"
)

var base58Runes = []rune("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func randBase58Seq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = base58Runes[rand.Intn(len(base58Runes))]
	}
	return string(b)
}

func GenerateDID() string {
	return "did:cheqd:test:" + randBase58Seq(16)
}

func GenerateFragment(did string) string {
	return did + "#key-1"
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
