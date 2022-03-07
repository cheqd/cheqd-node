package types

import (
	"errors"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

// Custom error rule

type CustomErrorRule struct {
	fn func(value interface{}) error
}

func NewCustomErrorRule(fn func(value interface{}) error) *CustomErrorRule {
	return &CustomErrorRule{fn: fn}
}

func (c CustomErrorRule) Validate(value interface{}) error {
	return c.fn(value)
}


// Validation helpers

func IsDID() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("IsDID must be only applied on string properties")
		}

		return utils.ValidateDID(casted, "", nil)
	})
}

func IsDIDUrl() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("IsDIDUrl must be only applied on string properties")
		}

		return utils.ValidateDIDUrl(casted, "", nil)
	})
}

func IsMultibase() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("IsMultibase must be only applied on string properties")
		}

		return utils.ValidateMultibase(casted)
	})
}

func IsUnique() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.([]string)
		if !ok {
			panic("IsSet must be only applied on string array properties")
		}

		if !utils.IsUnique(casted) {
			return errors.New("there should be no duplicates")
		}

		return nil
	})
}
