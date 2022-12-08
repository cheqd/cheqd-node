package utils

import (
	"encoding/json"
	"fmt"

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

func MustEncodeJson(data interface{}) string {
	encoded, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(encoded)
}
