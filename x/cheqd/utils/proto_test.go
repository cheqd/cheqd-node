package utils

import (
	"testing"

	bank_types "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
)

func Test_MsgTypeUrl(t *testing.T) {
	assert.Equal(t, "/cosmos.bank.v1beta1.DenomUnit", MsgTypeURL(&bank_types.DenomUnit{}))
}
