package types

import (
	"errors"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func NewSignInfo(verificationMethodId string, signature string) *SignInfo {
	return &SignInfo{VerificationMethodId: verificationMethodId, Signature: signature}
}

// Helpers

func GetSignInfoIds(vms []*SignInfo) []string {
	res := make([]string, len(vms))

	for i := range vms {
		res[i] = vms[i].VerificationMethodId
	}

	return res
}


// Validate

func (si SignInfo) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&si,
		validation.Field(&si.VerificationMethodId, validation.Required, IsDIDUrl(allowedNamespaces, Empty, Empty, Required)),
		validation.Field(&si.Signature, validation.Required, is.Base64),
	)
}

func ValidSignInfo(allowedNamespaces []string) *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(SignInfo)
		if !ok {
			panic("ValidSignInfo must be only applied on sign infos")
		}

		return casted.Validate(allowedNamespaces)
	})
}

func IsUniqueSignInfoList() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.([]*SignInfo)
		if !ok {
			panic("IsUniqueVerificationMethodList must be only applied on VM lists")
		}

		ids := GetSignInfoIds(casted)
		if !utils.IsUnique(ids) {
			return errors.New("there are sign info duplicates")
		}

		return nil
	})
}
