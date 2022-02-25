package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/go-playground/validator/v10"
)

func BuildValidator(allowedDIDNamespaces []string) (*validator.Validate, error) {
	validate := validator.New()

	err := validate.RegisterValidation("did", func(fl validator.FieldLevel) bool {
		return utils.IsValidDID(fl.Field().String(), allowedDIDNamespaces)
	})
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation("did-url", func(fl validator.FieldLevel) bool {
		return utils.IsValidDIDUrl(fl.Field().String(), allowedDIDNamespaces)
	})
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation("did-url-no-path", func(fl validator.FieldLevel) bool {
		_, _, _, path, _, _ := utils.SplitDIDUrl(fl.Field().String())
		return path == ""
	})
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation("did-url-no-query", func(fl validator.FieldLevel) bool {
		_, _, _, _, query, _ := utils.SplitDIDUrl(fl.Field().String())
		return query == ""
	})
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation("did-url-fragment", func(fl validator.FieldLevel) bool {
		_, _, _, _, _, fragment := utils.SplitDIDUrl(fl.Field().String())
		return fragment != ""
	})
	if err != nil {
		return nil, err
	}

	return validate, nil
}
