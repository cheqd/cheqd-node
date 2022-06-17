package types

import (
	cheqdTypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/resource/utils"
)

// Validation helpers

func IsUUID() *cheqdTypes.CustomErrorRule {
	return cheqdTypes.NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("IsDID must be only applied on string properties")
		}

		return utils.ValidateUUID(casted)
	})
}

func isResourceTypeAllowed()*cheqdTypes.CustomErrorRule {
	return cheqdTypes.NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("isResourceTypeAllowed must be only applied on string properties")
		}

		return utils.ValidateResourceType(casted)
	})
}

func isMimeTypeAllowed()*cheqdTypes.CustomErrorRule {
	return cheqdTypes.NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("isMimeTypeAllowed must be only applied on string properties")
		}

		return utils.ValidateMimeType(casted)
	})
}
