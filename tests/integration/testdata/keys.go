package testdata

import (
	"math/rand"
)

const ED25519_PRIVATE_KEY_LENGTH = 32

func GenerateByteEntropy() []byte {
	entropy := make([]byte, rand.Intn(ED25519_PRIVATE_KEY_LENGTH))
	rand.Read(entropy)
	return entropy
}