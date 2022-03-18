package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_PackUnpackAny(t *testing.T) {
	original := &Did{
		Id:                   "test",
	}

	// Construct codec
	registry := types.NewInterfaceRegistry()
	RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	// Marshal
	bz, err := cdc.MarshalInterface(original)
	require.NoError(t, err)

	// Assert type url
	var any types.Any
	err = any.Unmarshal(bz)
	assert.NoError(t, err)
	assert.Equal(t, any.TypeUrl, utils.MsgTypeURL(&Did{}))

	// Unmarshal
	var decoded StateValueData
	err = cdc.UnmarshalInterface(bz, &decoded)
	require.NoError(t, err)
	require.IsType(t, &Did{}, decoded)
	require.Equal(t, original, decoded)
}
