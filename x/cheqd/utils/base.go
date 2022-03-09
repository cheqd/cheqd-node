package utils

import (
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	multibase "github.com/multiformats/go-multibase"
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
	return ValidateMultibaseEncoding(string(multibase.Base58BTC) + data, multibase.Base58BTC)
}


func ValidateJWKEncoding(jwk_string string) error {
	_, err := jwk.ParseString(jwk_string)
	if err != nil {
		return fmt.Errorf("invalid format for JWK key, error from validation: %s", err.Error())
	}

	return nil
}