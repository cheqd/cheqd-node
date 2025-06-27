package v4_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"

	"github.com/cheqd/cheqd-node/x/resource"
	"github.com/cheqd/cheqd-node/x/resource/exported"
	v4 "github.com/cheqd/cheqd-node/x/resource/migration/v4"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
)

type mockSubspace struct {
	ps types.FeeParams
}

func newMockSubspace(ps types.FeeParams) mockSubspace {
	return mockSubspace{ps: ps}
}

func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps exported.ParamSet) {
	*ps.(*types.FeeParams) = ms.ps
}

func (ms mockSubspace) Get(ctx sdk.Context, key []byte, ps interface{}) {
	*ps.(*types.FeeParams) = ms.ps
}

func TestMigrate(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(resource.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := sdktestutil.DefaultContext(storeKey, tKey)
	kvStoreService := runtime.NewKVStoreService(storeKey)
	store := kvStoreService.OpenKVStore(ctx)
	sb := collections.NewSchemaBuilder(kvStoreService)
	countCollection := collections.NewItem(sb, collections.Prefix(types.ResourceCountKey),
		"resource_count", collections.Uint64Value)
	metadataCollection := collections.NewMap(
		sb,
		collections.NewPrefix(types.ResourceMetadataKey),
		"resource_metadata",
		collections.PairKeyCodec(collections.StringKey, collections.StringKey),
		codec.CollValue[types.Metadata](cdc),
	)
	dataCollection := collections.NewMap(
		sb,
		collections.NewPrefix(types.ResourceDataKey),
		"resource_data",
		collections.PairKeyCodec(collections.StringKey, collections.StringKey),
		collections.BytesValue,
	)
	latestResourceVersionCollection := collections.NewMap(
		sb,
		collections.NewPrefix(types.ResourceLatestVersionKey),
		"resource_latest_version",
		collections.TripleKeyCodec(collections.StringKey, collections.StringKey, collections.StringKey),
		collections.StringValue,
	)

	// set count key in old store
	var countValue uint64 = 5
	require.NoError(t, store.Set([]byte(types.ResourceCountKey), []byte(strconv.FormatUint(countValue, 10))))

	// set resource in old store
	testCollectionId := "collection-id"
	testId := "test-id"
	testData := []byte("testdata")
	metadata := types.Metadata{
		CollectionId: testCollectionId,
		Id:           testId,
		Name:         "test-resource",
	}
	require.NoError(t, store.Set([]byte(types.ResourceMetadataKey+testCollectionId+":"+testId), cdc.MustMarshal(&metadata)))
	require.NoError(t, store.Set(v4.LegacyResourceDataKey(testCollectionId, testId), testData))

	legacySubspace := newMockSubspace(*types.DefaultFeeParams())
	require.NoError(t, v4.MigrateStore(ctx, runtime.NewKVStoreService(storeKey), legacySubspace, cdc,
		countCollection, metadataCollection, dataCollection, latestResourceVersionCollection))

	var res types.FeeParams
	bz, err := store.Get(types.ParamStoreKeyFeeParams)
	require.NoError(t, err)
	require.NoError(t, cdc.Unmarshal(bz, &res))
	require.Equal(t, legacySubspace.ps, res)

	// check set count value
	actualCount, err := countCollection.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, countValue, actualCount)

	// check resource records
	metadataRes, err := metadataCollection.Get(ctx, collections.Join(testCollectionId, testId))
	require.NoError(t, err)
	require.Equal(t, metadata, metadataRes)

	dataRes, err := dataCollection.Get(ctx, collections.Join(testCollectionId, testId))
	require.NoError(t, err)
	require.Equal(t, testData, dataRes)
}
