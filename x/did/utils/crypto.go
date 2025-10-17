package utils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"filippo.io/edwards25519"

	"github.com/btcsuite/btcd/btcec/v2"
	btcececdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
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

	// Detect curve type
	if pubKey.Curve.Params().Name == "secp256k1" {
		// Use btcec for secp256k1
		secpKey := btcec.NewPublicKey(bigIntToFieldVal(pubKey.X), bigIntToFieldVal(pubKey.Y))
		return VerifySecp256k1Signature(secpKey, digest, signature)
	}

	ok := ecdsa.VerifyASN1(&pubKey, digest, signature)
	if !ok {
		return errors.New("invalid ecdsa signature")
	}
	return nil
}

func VerifySecp256k1Signature(pubKey *btcec.PublicKey, digest []byte, signature []byte) error {
	sig, err := btcececdsa.ParseDERSignature(signature)
	if err != nil {
		return err
	}

	if !sig.Verify(digest, pubKey) {
		return errors.New("invalid secp256k1 signature")
	}

	return nil
}

func GetEd25519VerificationKey2020(keyBytes []byte) []byte {
	return keyBytes[2:]
}

// bigIntToFieldVal converts a *big.Int to a *btcec.FieldVal.
func bigIntToFieldVal(x *big.Int) *btcec.FieldVal {
	var fv btcec.FieldVal

	// The byte slice must be 32 bytes (256 bits) long for secp256k1,
	// padded with leading zeros if necessary.
	// We use a fixed 32-byte array to ensure correct padding.
	paddedBytes := make([]byte, 32)
	xBytes := x.Bytes()
	copy(paddedBytes[32-len(xBytes):], xBytes)

	fv.SetByteSlice(paddedBytes)
	return &fv
}
