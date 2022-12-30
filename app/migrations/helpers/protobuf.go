package helpers

import (
	"github.com/multiformats/go-multibase"
)

func BuildEd25519VerificationKey2020VerificationMaterial(publicKey string) (string, error) {
	encoding, publicKeyBytes, err := multibase.Decode(publicKey)
	if encoding != multibase.Base58BTC {
		panic("Only Base58BTC encoding is supported")
	}
	if err != nil {
		return "", err
	}
	multicodecAndKey := []byte{0xed, 0x01}
	multicodecAndKey = append(multicodecAndKey, publicKeyBytes...)
	keyStr, err := multibase.Encode(multibase.Base58BTC, multicodecAndKey)
	return keyStr, err
}
