package types

import (
	"github.com/go-playground/validator/v10"
)

func BuildValidator(DIDMethod string, allowedDIDNamespaces []string) (*validator.Validate, error) {
	validate := validator.New()

	// Custom tags
	err := validate.RegisterValidation("did", func(fl validator.FieldLevel) bool {
		return IsValidDID(fl.Field().String(), DIDMethod, allowedDIDNamespaces)
	})
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation("did-url", func(fl validator.FieldLevel) bool {
		return IsValidDIDUrl(fl.Field().String(), DIDMethod, allowedDIDNamespaces)
	})
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation("did-url-no-path", func(fl validator.FieldLevel) bool {
		_, path, _, _ := TrySplitDIDUrl(fl.Field().String())
		return path == ""
	})
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation("did-url-no-query", func(fl validator.FieldLevel) bool {
		_, _, query, _ := TrySplitDIDUrl(fl.Field().String())
		return query == ""
	})
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation("did-url-with-fragment", func(fl validator.FieldLevel) bool {
		_, _, _, fragment := TrySplitDIDUrl(fl.Field().String())
		return fragment != ""
	})
	if err != nil {
		return nil, err
	}

	// Custom struct level rules

	validate.RegisterStructValidation(VerificationMethodStructLevelValidation, VerificationMethod{})

	return validate, nil
}

//const (
//	PublicKeyJwk       = "PublicKeyJwk"
//	PublicKeyMultibase = "PublicKeyMultibase"
//)
//
//var VerificationMethodType = map[string]string{
//	"JsonWebKey2020":             PublicKeyJwk,
//	"Ed25519VerificationKey2020": PublicKeyMultibase,
//}
//
//var ServiceType = []string{
//	"LinkedDomains",
//	"DIDCommMessaging",
//}
//
//func GetVerificationMethodType(vmType string) string {
//	return VerificationMethodType[vmType]
//}
//
//func IsValidDidServiceType(sType string) bool {
//	return strings.Contains(ServiceType, sType)
//}