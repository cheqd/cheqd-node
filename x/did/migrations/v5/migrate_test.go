package v5_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"

	did "github.com/cheqd/cheqd-node/x/did"
	"github.com/cheqd/cheqd-node/x/did/exported"
	v5 "github.com/cheqd/cheqd-node/x/did/migrations/v5"
	"github.com/cheqd/cheqd-node/x/did/types"
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
	encCfg := moduletestutil.MakeTestEncodingConfig(did.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := sdktestutil.DefaultContext(storeKey, tKey)
	kvStoreService := runtime.NewKVStoreService(storeKey)
	store := kvStoreService.OpenKVStore(ctx)

	legacySubspace := newMockSubspace(*types.DefaultFeeParams())
	require.NoError(t, v5.MigrateStore(ctx, runtime.NewKVStoreService(storeKey), legacySubspace, cdc))

	var res types.FeeParams
	bz, err := store.Get(types.ParamStoreKey)
	require.NoError(t, err)
	require.NoError(t, cdc.Unmarshal(bz, &res))
	require.Equal(t, legacySubspace.ps, res)
}
