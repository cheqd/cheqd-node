package types

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type VerificationMaterial interface {
	Type() string
	Validate() error
}

// Ed25519VerificationKey2020

type Ed25519VerificationKey2020 struct {
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

var _ VerificationMaterial = (*Ed25519VerificationKey2020)(nil)

func (vm Ed25519VerificationKey2020) Type() string {
	return "Ed25519VerificationKey2020"
}

func (vm Ed25519VerificationKey2020) Validate() error {
	return validation.ValidateStruct(&vm,
		validation.Field(&vm.PublicKeyMultibase, validation.Required, IsMultibase(), IsMultibaseEncodedEd25519PubKey()),
	)
}

// JsonWebKey2020

type JsonWebKey2020 struct {
	PublicKeyJwk json.RawMessage `json:"publicKeyJwk"`
}

var _ VerificationMaterial = (*JsonWebKey2020)(nil)

func (vm JsonWebKey2020) Type() string {
	return "JsonWebKey2020"
}

func (vm JsonWebKey2020) Validate() error {
	return validation.Validate(string(vm.PublicKeyJwk), validation.Required, IsJWK())
}

// Validation

func ValidEd25519VerificationKey2020Rule() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("ValidVerificationMethodRule must be only applied on verification methods")
		}

		var vm Ed25519VerificationKey2020
		err := json.Unmarshal([]byte(casted), &vm)
		if err != nil {
			return err
		}

		return vm.Validate()
	})
}

func ValidJsonWebKey2020Rule() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("ValidVerificationMethodRule must be only applied on verification methods")
		}

		var vm JsonWebKey2020
		err := json.Unmarshal([]byte(casted), &vm)
		if err != nil {
			return err
		}

		return vm.Validate()
	})
}
