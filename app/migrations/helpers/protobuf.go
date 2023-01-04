package helpers

import (
	"github.com/multiformats/go-multibase"
)

func GenerateEd25519VerificationKey2020VerificationMaterial(publicKey string) (string, error) {
	encoding, publicKeyBytes, err := multibase.Decode(publicKey)
	if encoding != multibase.Base58BTC {
		panic("Only Base58BTC encoding is supported")
	}
	if err != nil {
		return "", err
	}
	publicKeyMultibaseBytes := []byte{0xed, 0x01}
	publicKeyMultibaseBytes = append(publicKeyMultibaseBytes, publicKeyBytes...)

	return multibase.Encode(multibase.Base58BTC, publicKeyMultibaseBytes)
}
