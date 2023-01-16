package helpers

import (
	"time"

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

func MustParseFromStringTimeToGoTime(timeString string) time.Time {
	// If timeString is empty return default nullable time value (0001-01-01 00:00:00 +0000 UTC)
	if timeString == "" {
		return time.Time{}
	}

	t, err := time.Parse(time.RFC3339, timeString)
	if err == nil {
		return t
	}
	t, err = time.Parse(time.RFC3339Nano, timeString)
	if err == nil {
		return t
	}
	t, err = time.Parse(OldTimeFormat, timeString)
	if err != nil {
		panic(err)
	}
	return t
}
