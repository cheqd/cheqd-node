package tests

import (
	types2 "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_PackUnpackAny(t *testing.T) {
	original := &types2.Did{
		Id:                   "test",
	}

	// Construct codec
	registry := types.NewInterfaceRegistry()
	types2.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	// Marshal
	bz, err := cdc.MarshalInterface(original)
	require.NoError(t, err)

	// Assert type url
	var any types.Any
	err = any.Unmarshal(bz)
	assert.NoError(t, err)
	assert.Equal(t, any.TypeUrl, types2.MsgTypeURL(&types2.Did{}))

	// Unmarshal
	var decoded types2.StateValueData
	err = cdc.UnmarshalInterface(bz, &decoded)
	require.NoError(t, err)
	require.IsType(t, &types2.Did{}, decoded)
	require.Equal(t, original, decoded)
}
