package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MsgTypeUrl(t *testing.T) {
	assert.Equal(t, "/cheqdid.cheqdnode.cheqd.v1.Did", MsgTypeURL(&Did{}))
}
