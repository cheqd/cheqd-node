package types

import (
	"errors"

	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func NewService(id string, type_ string, serviceEndpoint string) *Service {
	return &Service{
		Id:              id,
		Type:            type_,
		ServiceEndpoint: serviceEndpoint,
	}
}

// Helpers

func GetServiceIds(vms []*Service) []string {
	res := make([]string, len(vms))

	for i := range vms {
		res[i] = vms[i].Id
	}

	return res
}

// Validation

func (s Service) Validate(baseDid string, allowedNamespaces []string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Id, validation.Required, IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(baseDid)),
		validation.Field(&s.Type, validation.Required, validation.Length(1, 255)),
		validation.Field(&s.ServiceEndpoint, validation.Required),
	)
}

func ValidServiceRule(baseDid string, allowedNamespaces []string) *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(Service)
		if !ok {
			panic("ValidVerificationMethodRule must be only applied on verification methods")
		}

		return casted.Validate(baseDid, allowedNamespaces)
	})
}

func IsUniqueServiceListByIdRule() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.([]*Service)
		if !ok {
			panic("IsUniqueServiceListByIdRule must be only applied on service lists")
		}

		ids := GetServiceIds(casted)
		if !utils.IsUnique(ids) {
			return errors.New("there are service duplicates")
		}

		return nil
	})
}
