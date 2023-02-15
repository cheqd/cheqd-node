package types

import (
	"errors"

	"github.com/cheqd/cheqd-node/x/did/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mr-tron/base58"
)

func NewSignInfo(verificationMethodID string, signature []byte) *SignInfo {
	return &SignInfo{VerificationMethodId: verificationMethodID, Signature: signature}
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
	hash := func(si *SignInfo) string {
		return si.VerificationMethodId + ":" + base58.Encode(si.Signature)
	}

	tmp := map[string]bool{}
	for _, si := range infos {
		h := hash(si)

		_, found := tmp[h]
		if found {
			return false
		}

		tmp[h] = true
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
	infosBS := FindSignInfosBySigner(infos, signer)

	if len(infosBS) == 0 {
		return SignInfo{}, false
	}

	return infosBS[0], true
}

// Validate

func (si SignInfo) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&si,
		validation.Field(&si.VerificationMethodId, validation.Required, IsDIDUrl(allowedNamespaces, Empty, Empty, Required)),
		validation.Field(&si.Signature, validation.Required),
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

func IsUniqueSignInfoListByIDRule() *CustomErrorRule {
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

// Normalization

func (si *SignInfo) Normalize() {
	si.VerificationMethodId = utils.NormalizeDIDUrl(si.VerificationMethodId)
}

func NormalizeSignInfoList(signatures []*SignInfo) {
	for _, s := range signatures {
		s.Normalize()
	}
}
