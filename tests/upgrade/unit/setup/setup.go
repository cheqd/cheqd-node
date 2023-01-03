package setup

import (
	"context"
	"crypto/rand"
	"time"

	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didkeeperv1 "github.com/cheqd/cheqd-node/x/did/keeper/v1"
	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"

	// didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	resourcekeeperv1 "github.com/cheqd/cheqd-node/x/resource/keeper/v1"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

type TestSetup struct {
	Cdc codec.Codec

	SdkCtx sdk.Context
	StdCtx context.Context

	DidKeeperV1      didkeeperv1.Keeper
	ResourceKeeperV1 resourcekeeperv1.Keeper

	DidKeeper      didkeeper.Keeper
	DidMsgServer   didtypes.MsgServer
	DidQueryServer didtypes.QueryServer

	ResourceKeeper      resourcekeeper.Keeper
	ResourceMsgServer   resourcetypes.MsgServer
	ResourceQueryServer resourcetypes.QueryServer

	DidStoreKey      *storetypes.KVStoreKey
	ResourceStoreKey *storetypes.KVStoreKey

	ParamsKeeper paramskeeper.Keeper
}

func Setup() TestSetup {
	// Init Codec
	ir := codectypes.NewInterfaceRegistry()
	didtypes.RegisterInterfaces(ir)
	// didtypesv1.RegisterInterfaces(ir) // TODO: Is v1 needed?
	Cdc := codec.NewProtoCodec(ir)
	aminoCdc := codec.NewLegacyAmino()

	// Init KVSore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)

	didStoreKey := sdk.NewKVStoreKey(didtypes.StoreKey)
	resourceStoreKey := sdk.NewKVStoreKey(resourcetypes.StoreKey)

	dbStore.MountStoreWithDB(didStoreKey, storetypes.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(resourceStoreKey, storetypes.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init ParamsKeeper KVStore
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	paramsTStoreKey := sdk.NewTransientStoreKey(paramstypes.TStoreKey)

	// Init Keepers
	paramsKeeper := initParamsKeeper(Cdc, aminoCdc, paramsStoreKey, paramsTStoreKey)

	// Init previous keepers
	didKeeperPrevious := didkeeperv1.NewKeeper(Cdc, didStoreKey)
	resourceKeeperPrevious := resourcekeeperv1.NewKeeper(Cdc, resourceStoreKey)

	// Init Keepers
	didKeeper := didkeeper.NewKeeper(Cdc, didStoreKey, getSubspace(didtypes.ModuleName, paramsKeeper))
	resourceKeeper := resourcekeeper.NewKeeper(Cdc, resourceStoreKey, getSubspace(resourcetypes.ModuleName, paramsKeeper))

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

	resourceMsgServer := resourcekeeper.NewMsgServer(*resourceKeeper, *didKeeper)
	resourceQueryServer := resourcekeeper.NewQueryServer(*resourceKeeper, *didKeeper)

	setup := TestSetup{
		Cdc: Cdc,

		SdkCtx: ctx,
		StdCtx: sdk.WrapSDKContext(ctx),

		DidKeeperV1:      *didKeeperPrevious,
		ResourceKeeperV1: *resourceKeeperPrevious,

		DidKeeper:      *didKeeper,
		DidMsgServer:   didMsgServer,
		DidQueryServer: didQueryServer,

		ResourceKeeper:      *resourceKeeper,
		ResourceMsgServer:   resourceMsgServer,
		ResourceQueryServer: resourceQueryServer,

		DidStoreKey:      didStoreKey,
		ResourceStoreKey: resourceStoreKey,

		ParamsKeeper: paramsKeeper,
	}

	setup.DidKeeper.SetDidNamespace(&ctx, didsetup.DidNamespace) // TODO: Think about it
	return setup
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key storetypes.StoreKey, tkey storetypes.StoreKey) paramskeeper.Keeper {
	// create keeper
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// set params subspaces
	paramsKeeper.Subspace(didtypes.ModuleName)

	return paramsKeeper
}

func getSubspace(moduleName string, paramsKeeper paramskeeper.Keeper) paramstypes.Subspace {
	subspace, _ := paramsKeeper.GetSubspace(moduleName)
	return subspace
}
