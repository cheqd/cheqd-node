package utils

import (
	"crypto/ed25519"
	"filippo.io/edwards25519"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
)

func ValidateJWK(jwk_string string) error {
	_, err := jwk.ParseString(jwk_string)
	if err != nil {
		return fmt.Errorf("invalid format for JWK key, error from validation: %s", err.Error())
	}

	return nil
}

func ValidateEd25519PubKey(keyBytes []byte) error {
	if l := len(keyBytes); l != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519: bad public key length: %d", l)
	}
	_, err := (&edwards25519.Point{}).SetBytes(keyBytes)
	if err != nil {
		return err
	}
	return nil
}