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

	// set count key in old store
	var countValue uint64 = 5
	store.Set([]byte(types.ResourceCountKey), []byte(strconv.FormatUint(countValue, 10)))

	legacySubspace := newMockSubspace(*types.DefaultFeeParams())
	require.NoError(t, v4.MigrateStore(ctx, runtime.NewKVStoreService(storeKey), legacySubspace, cdc, countCollection))

	var res types.FeeParams
	bz, err := store.Get(types.ParamStoreKeyFeeParams)
	require.NoError(t, err)
	require.NoError(t, cdc.Unmarshal(bz, &res))
	require.Equal(t, legacySubspace.ps, res)

	// check set count value
	actualCount, err := countCollection.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, countValue, actualCount)
}
