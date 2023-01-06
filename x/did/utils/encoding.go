package utils

import (
	"encoding/json"
	"fmt"

	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multibase"
)

func ValidateMultibase(data string) error {
	_, _, err := multibase.Decode(data)
	return err
}

func ValidateMultibaseEncoding(data string, expectedEncoding multibase.Encoding) error {
	actualEncoding, _, err := multibase.Decode(data)
	if err != nil {
		return err
	}

	if actualEncoding != expectedEncoding {
		return fmt.Errorf("invalid actualEncoding. expected: %s actual: %s",
			multibase.EncodingToStr[expectedEncoding], multibase.EncodingToStr[actualEncoding])
	}

	return nil
}

func ValidateBase58(data string) error {
	return ValidateMultibaseEncoding(string(multibase.Base58BTC)+data, multibase.Base58BTC)
}

func IsValidBase58(data string) bool {
	return ValidateBase58(data) == nil
}

func MustEncodeMultibaseBase58(data []byte) string {
	encoded, err := multibase.Encode(multibase.Base58BTC, data)
	if err != nil {
		panic(err)
	}

	return encoded
}

func MustEncodeJSON(data interface{}) string {
	encoded, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(encoded)
}

func ValidateMulticodecEd25519VerificationKey2020(keyBytes []byte) error {
	if keyBytes[0] != 0xed && keyBytes[1] != 0x01 {
		return fmt.Errorf("invalid two-byte prefix for Ed25519VerificationKey2020. expected: %s actual: %s",
			"0xed01", fmt.Sprintf("0x%02x%02x", keyBytes[0], keyBytes[1]))
	}
	return nil
}

func ValidateMultibaseEd25519VerificationKey2020(data string) error {
	encoding, keyBytes, err := multibase.Decode(data)
	if err != nil {
		return err
	}

	if encoding != multibase.Base58BTC {
		return fmt.Errorf("invalid encoding for Ed25519VerificationKey2020. expected: %s actual: %s",
			multibase.EncodingToStr[multibase.Base58BTC], multibase.EncodingToStr[encoding])
	}

	err = ValidateMulticodecEd25519VerificationKey2020(keyBytes)
	if err != nil {
		return err
	}

	pubKey := GetEd25519VerificationKey2020(keyBytes)
	return ValidateEd25519PubKey(pubKey)
}

func ValidateBase58Ed25519VerificationKey2018(data string) error {
	pubKey, err := base58.Decode(data)
	if err != nil {
		return err
	}
	return ValidateEd25519PubKey(pubKey)
}
