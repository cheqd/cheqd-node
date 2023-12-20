package setup

import (
	"context"
	"crypto/rand"
	"time"

	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didkeeperv1 "github.com/cheqd/cheqd-node/x/did/keeper/v1"
	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	"github.com/cheqd/cheqd-node/x/resource"
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
	portkeeper "github.com/cosmos/ibc-go/v6/modules/core/05-port/keeper"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
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

	ResourceKeeper       resourcekeeper.Keeper
	ResourceMsgServer    resourcetypes.MsgServer
	ResourceQueryServer  resourcetypes.QueryServer
	IBCModule            resource.IBCModule
	PortKeeper           portkeeper.Keeper
	ScopedResourceKeeper capabilitykeeper.ScopedKeeper

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

	// mount did store
	didStoreKey := sdk.NewKVStoreKey(didtypes.StoreKey)
	dbStore.MountStoreWithDB(didStoreKey, storetypes.StoreTypeIAVL, nil)

	// mount resource store
	resourceStoreKey := sdk.NewKVStoreKey(resourcetypes.StoreKey)
	dbStore.MountStoreWithDB(resourceStoreKey, storetypes.StoreTypeIAVL, nil)

	// mount capability store - required for ibc port tests
	capabilityStoreKey := sdk.NewKVStoreKey(capabilitytypes.StoreKey)
	dbStore.MountStoreWithDB(capabilityStoreKey, storetypes.StoreTypeIAVL, nil)
	memStoreKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	dbStore.MountStoreWithDB(memStoreKeys[capabilitytypes.MemStoreKey], storetypes.StoreTypeMemory, nil)

	// mount param store - required for ibc port tests with default genesis
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	dbStore.MountStoreWithDB(paramsStoreKey, storetypes.StoreTypeIAVL, nil)
	paramsTStoreKey := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	dbStore.MountStoreWithDB(paramsTStoreKey, storetypes.StoreTypeTransient, nil)

	_ = dbStore.LoadLatestVersion()

	// init previous keepers
	didKeeperPrevious := didkeeperv1.NewKeeper(Cdc, didStoreKey)
	resourceKeeperPrevious := resourcekeeperv1.NewKeeper(Cdc, resourceStoreKey)

	// init Keepers
	paramsKeeper := initParamsKeeper(Cdc, aminoCdc, paramsStoreKey, paramsTStoreKey)
	didKeeper := didkeeper.NewKeeper(Cdc, didStoreKey, getSubspace(didtypes.ModuleName, paramsKeeper))
	capabilityKeeper := capabilitykeeper.NewKeeper(Cdc, capabilityStoreKey, memStoreKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibchost.ModuleName)
	portKeeper := portkeeper.NewKeeper(scopedIBCKeeper)
	scopedResourceKeeper := capabilityKeeper.ScopeToModule(resourcetypes.ModuleName)
	resourceKeeper := resourcekeeper.NewKeeper(Cdc, resourceStoreKey, getSubspace(resourcetypes.ModuleName, paramsKeeper), &portKeeper, scopedResourceKeeper)

	// init IBC Module
	ibcModule := resource.NewIBCModule(*resourceKeeper)

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

		ResourceKeeper:       *resourceKeeper,
		ResourceMsgServer:    resourceMsgServer,
		ResourceQueryServer:  resourceQueryServer,
		IBCModule:            ibcModule,
		PortKeeper:           portKeeper,
		ScopedResourceKeeper: scopedResourceKeeper,

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
