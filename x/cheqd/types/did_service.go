package types

import (
	"errors"
	"strings"

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

func NormalizeIdentifier(did string) string {
	_, _, sUniqueId, err := utils.TrySplitDID(did)
	if err != nil {
		return did
	}
	if utils.IsValidUUID(sUniqueId) {
		return strings.ToLower(did)
	}
	return did
}

func NormalizeIdForFragmentUrl(didUrl string) string {
	id, _, _, fragmentId, err := utils.TrySplitDIDUrl(didUrl)
	if err != nil {
		return didUrl
	}
	if fragmentId == "" {
		return NormalizeIdentifier(id)
	}
	return NormalizeIdentifier(id) + "#" + fragmentId
}

func NormalizeDID(didDoc *Did) *Did {
	didDoc.Id = NormalizeIdentifier(didDoc.Id)
	for _, vm := range didDoc.VerificationMethod {
		vm.Controller = NormalizeIdentifier(vm.Controller)
		vm.Id = NormalizeIdForFragmentUrl(vm.Id)
	}
	for _, s := range didDoc.Service {
		s.Id = NormalizeIdForFragmentUrl(s.Id)
	}
	didDoc.Controller = NormalizeIdentifiersList(didDoc.Controller)
	didDoc.Authentication = NormalizeIdentifiersList(didDoc.Authentication)
	didDoc.AssertionMethod = NormalizeIdentifiersList(didDoc.AssertionMethod)
	didDoc.CapabilityInvocation = NormalizeIdentifiersList(didDoc.CapabilityInvocation)
	didDoc.CapabilityDelegation = NormalizeIdentifiersList(didDoc.CapabilityDelegation)
	didDoc.KeyAgreement = NormalizeIdentifiersList(didDoc.KeyAgreement)
	didDoc.AlsoKnownAs = NormalizeIdentifiersList(didDoc.AlsoKnownAs)
	return didDoc
}

func NormalizeIdentifiersList(keys []string) []string {
	if keys == nil {
		return nil
	}
	newKeys := []string{}
	for _, id := range keys {
		newKeys = append(newKeys, NormalizeIdForFragmentUrl(id))
	}
	return newKeys
}

func NormalizeSignatureUUIDIdentifiers(signatures []*SignInfo) {
	for _, s := range signatures {
		s.VerificationMethodId = NormalizeIdForFragmentUrl(s.VerificationMethodId)
	}
}
