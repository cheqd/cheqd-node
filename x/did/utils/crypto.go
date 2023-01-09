package utils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"errors"
	"fmt"
	"reflect"

	"filippo.io/edwards25519"

	"github.com/lestrrat-go/jwx/jwk"
)

func ValidateJWK(jwkString string) error {
	var raw interface{}
	err := jwk.ParseRawKey([]byte(jwkString), &raw)
	if err != nil {
		return fmt.Errorf("can't parse jwk: %s", err.Error())
	}

	switch key := raw.(type) {
	case *rsa.PublicKey:
		break
	case *ecdsa.PublicKey:
		break
	case ed25519.PublicKey:
		err := ValidateEd25519PubKey(key)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported jwk type: %s. supported types are: rsa/pub, ecdsa/pub, ed25519/pub", reflect.TypeOf(raw).Name())
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

func VerifyED25519Signature(pubKey ed25519.PublicKey, message []byte, signature []byte) error {
	valid := ed25519.Verify(pubKey, message, signature)
	if !valid {
		return errors.New("invalid ed25519 signature")
	}

	return nil
}

// VerifyRSASignature uses PSS padding and SHA256 digest
// A good explanation of different paddings: https://security.stackexchange.com/questions/183179/what-is-rsa-oaep-rsa-pss-in-simple-terms
func VerifyRSASignature(pubKey rsa.PublicKey, message []byte, signature []byte) error {
	hasher := crypto.SHA256.New()
	hasher.Write(message)
	digest := hasher.Sum(nil)

	err := rsa.VerifyPSS(&pubKey, crypto.SHA256, digest, signature, nil)
	if err != nil {
		return err
	}
	return nil
}

// VerifyECDSASignature uses ASN1 to decode r and s, SHA265 to calculate message digest
func VerifyECDSASignature(pubKey ecdsa.PublicKey, message []byte, signature []byte) error {
	hasher := crypto.SHA256.New()
	hasher.Write(message)
	digest := hasher.Sum(nil)

	ok := ecdsa.VerifyASN1(&pubKey, digest, signature)
	if !ok {
		return errors.New("invalid ecdsa signature")
	}
	return nil
}

func GetEd25519VerificationKey2020(keyBytes []byte) []byte {
	return keyBytes[2:]
}
