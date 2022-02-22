package tests

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MsgTypeUrl(t *testing.T) {
	assert.Equal(t, "/cheqdid.cheqdnode.cheqd.v1.Did", types.MsgTypeURL(&types.Did{}))
}
