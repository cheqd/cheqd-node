package types

import (
	utils "github.com/cheqd/cheqd-node/x/resource/utils"
)

// Helper enums

type ValidationType int

const (
	Optional ValidationType = iota
	Required ValidationType = iota
	Empty    ValidationType = iota
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

// Validate URL

func IsUUID() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(string)
		if !ok {
			panic("IsDID must be only applied on string properties")
		}

		return utils.ValidateUUID(casted)
	})
}
