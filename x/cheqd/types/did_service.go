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

func UpdateUUIDForDID(did string) string {
	_, _, sUniqueId, err := utils.TrySplitDID(did)
	if err != nil {
		return did
	}
	if utils.IsValidUUID(sUniqueId) {
		return strings.ToLower(did)
	}
	return did
}

func UpdateUUIDForFragmentUrl(didUrl string) string {
	id, _, _, fragmentId, err := utils.TrySplitDIDUrl(didUrl)
	if err != nil {
		return didUrl
	}
	return UpdateUUIDForDID(id) + "#" + fragmentId
}

func UpdateUUIDIdentifiers(didDoc *Did) {
	didDoc.Id = UpdateUUIDForDID(didDoc.Id)
	for _, vm := range didDoc.VerificationMethod {
		vm.Controller = UpdateUUIDForDID(vm.Controller)
		vm.Id = UpdateUUIDForFragmentUrl(vm.Id)
	}
	for _, s := range didDoc.Service {
		s.Id = UpdateUUIDForFragmentUrl(s.Id)
	}
	didDoc.Authentication = UpdateDidKeyIdentifiersList(didDoc.Authentication)
	didDoc.AssertionMethod = UpdateDidKeyIdentifiersList(didDoc.AssertionMethod)
	didDoc.CapabilityInvocation = UpdateDidKeyIdentifiersList(didDoc.CapabilityInvocation)
	didDoc.CapabilityDelegation = UpdateDidKeyIdentifiersList(didDoc.CapabilityDelegation)
	didDoc.KeyAgreement = UpdateDidKeyIdentifiersList(didDoc.KeyAgreement)
	didDoc.AlsoKnownAs = UpdateDidKeyIdentifiersList(didDoc.AlsoKnownAs)
}

func UpdateDidKeyIdentifiersList(keys []string) []string {
	newKeys := []string{}
	for _, id := range keys {
		newKeys = append(newKeys, UpdateUUIDForFragmentUrl(id))
	}
	return newKeys
}

func UpdateSignatureUUIDIdentifiers(signatures []*SignInfo) {
	for _, s := range signatures {
		s.VerificationMethodId = UpdateUUIDForFragmentUrl(s.VerificationMethodId)
	}
}
