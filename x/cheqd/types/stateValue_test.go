package types

import (
	"testing"
	"time"

	"github.com/canow-co/cheqd-node/x/cheqd/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func Test_PackUnpackAny(t *testing.T) {
	original := &Did{
		Id: "test",
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

func Test_NewMetadataFromContext(t *testing.T) {
	createdTime := time.Now()
	ctx := sdk.NewContext(nil, tmproto.Header{ChainID: "test_chain_id", Time: createdTime}, true, nil)
	ctx.WithTxBytes([]byte("test_tx"))
	expectedMetadata := Metadata{
		Created:     createdTime.UTC().Format(time.RFC3339),
		Updated:     "",
		Deactivated: false,
		VersionId:   utils.GetTxHash(ctx.TxBytes()),
	}

	metadata := NewMetadataFromContext(ctx)

	require.Equal(t, expectedMetadata, metadata)
}

func Test_UpdateMetadata(t *testing.T) {
	createdTime := time.Now()
	updatedTime := createdTime.Add(time.Hour)

	ctx1 := NewContext(createdTime, []byte("test1_tx"))
	ctx2 := NewContext(updatedTime, []byte("test1_tx"))

	expectedMetadata := Metadata{
		Created:     createdTime.UTC().Format(time.RFC3339),
		Updated:     updatedTime.UTC().Format(time.RFC3339),
		Deactivated: false,
		VersionId:   utils.GetTxHash(ctx2.TxBytes()),
	}

	metadata := NewMetadataFromContext(ctx1)
	metadata.Update(ctx2)

	require.Equal(t, expectedMetadata, metadata)
}

func NewContext(time time.Time, txBytes []byte) sdk.Context {
	ctx := sdk.NewContext(nil, tmproto.Header{ChainID: "test_chain_id", Time: time}, true, nil)
	ctx.WithTxBytes(txBytes)
	return ctx
}
