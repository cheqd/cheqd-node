package tests

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"
	"testing"
)

var letters = []rune("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


func GenerateDID() string {
	return "did:cheqd:test:" + randSeq(16)
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

func InitEnv(t *testing.T, keys map[string]KeyPair) (TestSetup){
	setup := Setup()
	err := setup.CreateTestDIDs(keys)
	require.NoError(t, err)
	return setup
}