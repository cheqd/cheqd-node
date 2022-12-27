package types

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// validation
func (au AlternativeUri) Validate() error {
	return validation.ValidateStruct(&au,
		validation.Field(&au.Uri, validation.Required, validation.Length(1, 256)),
		validation.Field(&au.Description, validation.Length(1, 128)),
	)
}

func ValidAlternativeURI() *didtypes.CustomErrorRule {
	return didtypes.NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(AlternativeUri)
		if !ok {
			panic("ValidAlternativeUri must be only applied on AlternativeUri properties")
		}

		return casted.Validate()
	})
}
