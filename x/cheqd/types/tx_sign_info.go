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

func GetSignInfoIds(infos []*SignInfo) []string {
	res := make([]string, len(infos))

	for i := range infos {
		res[i] = infos[i].VerificationMethodId
	}

	return res
}

func IsUniqueSignInfoList(infos []*SignInfo) bool {
	tmp_ := map[SignInfo]bool{}
	for _, si := range infos {
		_, found := tmp_[*si]
		if found {
			return false
		}
		tmp_[*si] = true
	}
	return true
}

// FindSignInfosBySigner returns the sign infos that corresponds to the provided signer's did
func FindSignInfosBySigner(infos []*SignInfo, signer string) []SignInfo {
	var result []SignInfo

	for _, info := range infos {
		did, _, _, _ := utils.MustSplitDIDUrl(info.VerificationMethodId)

		if did == signer {
			result = append(result, *info)
		}
	}

	return result
}

// FindSignInfoBySigner returns the first sign info that corresponds to the provided signer's did
func FindSignInfoBySigner(infos []*SignInfo, signer string) (info SignInfo, found bool) {
	infos_ := FindSignInfosBySigner(infos, signer)

	if len(infos_) == 0 {
		return SignInfo{}, false
	}

	return infos_[0], true
}

// Validate

func (si SignInfo) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&si,
		validation.Field(&si.VerificationMethodId, validation.Required, IsDIDUrl(allowedNamespaces, Empty, Empty, Required)),
		validation.Field(&si.Signature, validation.Required, is.Base64),
	)
}

func ValidSignInfoRule(allowedNamespaces []string) *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(SignInfo)
		if !ok {
			panic("ValidSignInfoRule must be only applied on sign infos")
		}

		return casted.Validate(allowedNamespaces)
	})
}

func IsUniqueSignInfoListByIdRule() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.([]*SignInfo)
		if !ok {
			panic("IsUniqueVerificationMethodListByIdRule must be only applied on VM lists")
		}

		ids := GetSignInfoIds(casted)
		if !utils.IsUnique(ids) {
			return errors.New("there are sign info records with the same ID")
		}

		return nil
	})
}

func IsUniqueSignInfoListRule() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.([]*SignInfo)
		if !ok {
			panic("IsUniqueVerificationMethodListByIdRule must be only applied on VM lists")
		}

		if !IsUniqueSignInfoList(casted) {
			return errors.New("there are full sign info duplicates")
		}

		return nil
	})
}
