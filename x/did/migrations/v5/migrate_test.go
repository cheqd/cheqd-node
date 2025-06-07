package v5_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"

	did "github.com/cheqd/cheqd-node/x/did"
	"github.com/cheqd/cheqd-node/x/did/exported"
	v5 "github.com/cheqd/cheqd-node/x/did/migrations/v5"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
)

type mockSubspace struct {
	ps types.LegacyFeeParams
}

func newMockSubspace(ps types.LegacyFeeParams) mockSubspace {
	return mockSubspace{ps: ps}
}

func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps exported.ParamSet) {
	*ps.(*types.LegacyFeeParams) = ms.ps
}

func (ms mockSubspace) Get(ctx sdk.Context, key []byte, ps interface{}) {
	*ps.(*types.LegacyFeeParams) = ms.ps
}

func TestMigrate(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(did.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := sdktestutil.DefaultContext(storeKey, tKey)
	kvStoreService := runtime.NewKVStoreService(storeKey)
	store := kvStoreService.OpenKVStore(ctx)
	sb := collections.NewSchemaBuilder(kvStoreService)
	countCollection := collections.NewItem(sb, collections.Prefix(types.DidDocCountKey),
		"did_count", collections.Uint64Value)
	docCollection := collections.NewMap(sb, types.DidDocVersionKeyPrefix, "did_version",
		collections.PairKeyCodec(collections.StringKey, collections.StringKey),
		codec.CollValue[types.DidDocWithMetadata](cdc))

	// set count key in old store
	var countValue uint64 = 5
	require.NoError(t, store.Set([]byte(types.DidDocCountKey), []byte(strconv.FormatUint(countValue, 10))))

	// set document in old store
	testId := "test-id"
	testVersion := "test-version"
	doc := types.DidDocWithMetadata{
		DidDoc: &types.DidDoc{
			Id: testId,
		},
		Metadata: &types.Metadata{
			VersionId: testVersion,
		},
	}
	require.NoError(t, store.Set([]byte(types.DidDocVersionKey+testId+":"+testVersion), cdc.MustMarshal(&doc)))

	legacySubspace := newMockSubspace(*types.DefaultLegacyFeeParams())
	require.NoError(t, v5.MigrateStore(ctx, kvStoreService, legacySubspace, cdc, countCollection, docCollection))

	var res types.LegacyFeeParams
	bz, err := store.Get(types.ParamStoreKey)
	require.NoError(t, err)
	require.NoError(t, cdc.Unmarshal(bz, &res))
	require.Equal(t, legacySubspace.ps, res)

	// check set count value
	actualCount, err := countCollection.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, countValue, actualCount)

	// check document
	docRes, err := docCollection.Get(ctx, collections.Join(testId, testVersion))
	require.NoError(t, err)
	require.Equal(t, doc, docRes)
}
