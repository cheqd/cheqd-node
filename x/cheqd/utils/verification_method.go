package utils

import "fmt"

var MultibaseStartSymbol = "z"

func ValidatePublicKeyMultibase(key string) error {
	// key cannot be empty
	if key == "" {
		return fmt.Errorf("PublicKeyMultibase cannot be empty string")
	}
	// check that starts from z
	if key[0:1] != MultibaseStartSymbol {
		return fmt.Errorf("PublicKeyMultibase should be started with #{MultibaseStartSymbol}")
	}
	// check that it's a base58 string
	if !ValidBase58Regexp.MatchString(key) {
		return fmt.Errorf("PublicKeyMultibase should be a string in base58 format")
	}
	return nil
}
