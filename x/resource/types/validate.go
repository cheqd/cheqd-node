package types

import (
	cheqdTypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/resource/utils"
)

// Validation helpers

func IsAllowedResourceType() *cheqdTypes.CustomErrorRule {
	return cheqdTypes.NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("IsAllowedResourceType must be only applied on string properties")
		}

		return utils.ValidateResourceType(casted)
	})
}

func IsAllowedMediaType() *cheqdTypes.CustomErrorRule {
	return cheqdTypes.NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("IsAllowedMediaType must be only applied on string properties")
		}

		return utils.ValidateMediaType(casted)
	})
}
