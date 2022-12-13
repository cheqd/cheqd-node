package setup

import (
	"crypto/rand"
	"time"

	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cheqd/cheqd-node/x/resource/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type TestSetup struct {
	didsetup.TestSetup

	ResourceKeeper      keeper.Keeper
	ResourceMsgServer   types.MsgServer
	ResourceQueryServer types.QueryServer
}

func Setup() TestSetup {
	// Init Codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	didtypes.RegisterInterfaces(ir)
	cdc := codec.NewProtoCodec(ir)
	aminoCdc := codec.NewLegacyAmino()

	// Init KVSore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)

	didStoreKey := sdk.NewKVStoreKey(didtypes.StoreKey)
	resourceStoreKey := sdk.NewKVStoreKey(types.StoreKey)

	dbStore.MountStoreWithDB(didStoreKey, storetypes.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(resourceStoreKey, storetypes.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init ParamsKeeper KVStore
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	paramsTStoreKey := sdk.NewTransientStoreKey(paramstypes.TStoreKey)

	// Init Keepers
	paramsKeeper := initParamsKeeper(cdc, aminoCdc, paramsStoreKey, paramsTStoreKey)
	didKeeper := didkeeper.NewKeeper(cdc, didStoreKey, getSubspace(didtypes.ModuleName, paramsKeeper))
	resourceKeeper := keeper.NewKeeper(cdc, resourceStoreKey, getSubspace(types.ModuleName, paramsKeeper))

	// Create Tx
	txBytes := make([]byte, 28)
	_, _ = rand.Read(txBytes)

	// Create context
	blockTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00.000Z")
	ctx := sdk.NewContext(dbStore,
		tmproto.Header{ChainID: "test", Time: blockTime},
		false, log.NewNopLogger()).WithTxBytes(txBytes)

	// Init servers
	didMsgServer := didkeeper.NewMsgServer(*didKeeper)
	didQueryServer := didkeeper.NewQueryServer(*didKeeper)

	msgServer := keeper.NewMsgServer(*resourceKeeper, *didKeeper)
	queryServer := keeper.NewQueryServer(*resourceKeeper, *didKeeper)

	setup := TestSetup{
		TestSetup: didsetup.TestSetup{
			Cdc: cdc,

			SdkCtx: ctx,
			StdCtx: sdk.WrapSDKContext(ctx),

			Keeper:      *didKeeper,
			MsgServer:   didMsgServer,
			QueryServer: didQueryServer,
		},

		ResourceKeeper:      *resourceKeeper,
		ResourceMsgServer:   msgServer,
		ResourceQueryServer: queryServer,
	}

	setup.Keeper.SetDidNamespace(&ctx, didsetup.DID_NAMESPACE)
	return setup
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key storetypes.StoreKey, tkey storetypes.StoreKey) paramskeeper.Keeper {
	// create keeper
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// set params subspaces
	paramsKeeper.Subspace(didtypes.ModuleName)
	paramsKeeper.Subspace(types.ModuleName)

	return paramsKeeper
}

func getSubspace(moduleName string, paramsKeeper paramskeeper.Keeper) paramstypes.Subspace {
	subspace, _ := paramsKeeper.GetSubspace(moduleName)
	return subspace
}
