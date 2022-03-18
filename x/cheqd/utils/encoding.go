package utils

import (
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
