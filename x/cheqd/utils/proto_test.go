package utils_test

import (
	"testing"

	bank_types "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"

	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func Test_MsgTypeUrl(t *testing.T) {
	assert.Equal(t, "/cosmos.bank.v1beta1.DenomUnit", utils.MsgTypeURL(&bank_types.DenomUnit{}))
}
