package setup

import (
	"crypto/rand"
	"time"

	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"

	"github.com/cheqd/cheqd-node/x/resource"
	"github.com/cheqd/cheqd-node/x/resource/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	portkeeper "github.com/cosmos/ibc-go/v7/modules/core/05-port/keeper"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
)

type TestSetup struct {
	didsetup.TestSetup

	ResourceKeeper      keeper.Keeper
	ResourceMsgServer   types.MsgServer
	ResourceQueryServer types.QueryServer
	IBCModule           resource.IBCModule
}

func Setup() TestSetup {
	// Init Codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	didtypes.RegisterInterfaces(ir)
	cdc := codec.NewProtoCodec(ir)
	aminoCdc := codec.NewLegacyAmino()

	// Init KVStore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)

	keys := sdk.NewKVStoreKeys(
		capabilitytypes.StoreKey,
		authtypes.StoreKey,
		banktypes.StoreKey,
	)
	// Mount did store
	didStoreKey := sdk.NewKVStoreKey(didtypes.StoreKey)
	dbStore.MountStoreWithDB(didStoreKey, storetypes.StoreTypeIAVL, nil)

	// Mount resource store
	resourceStoreKey := sdk.NewKVStoreKey(types.StoreKey)
	dbStore.MountStoreWithDB(resourceStoreKey, storetypes.StoreTypeIAVL, nil)

	// Mount capability store - required for ibc port tests
	capabilityStoreKey := sdk.NewKVStoreKey(capabilitytypes.StoreKey)
	dbStore.MountStoreWithDB(capabilityStoreKey, storetypes.StoreTypeIAVL, nil)
	memStoreKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	dbStore.MountStoreWithDB(memStoreKeys[capabilitytypes.MemStoreKey], storetypes.StoreTypeMemory, nil)

	// Mount param store - required for ibc port tests with default genesis
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	dbStore.MountStoreWithDB(paramsStoreKey, storetypes.StoreTypeIAVL, nil)
	paramsTStoreKey := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	dbStore.MountStoreWithDB(paramsTStoreKey, storetypes.StoreTypeTransient, nil)

	_ = dbStore.LoadLatestVersion()

	// Init Keepers
	accountKeeper := authkeeper.NewAccountKeeper(cdc, keys[authtypes.StoreKey], authtypes.ProtoBaseAccount, map[string][]string{}, "cheqd", string(authtypes.NewModuleAddress(govtypes.ModuleName)))
	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		keys[banktypes.StoreKey],
		accountKeeper,
		nil,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	paramsKeeper := initParamsKeeper(cdc, aminoCdc, paramsStoreKey, paramsTStoreKey)
	didKeeper := didkeeper.NewKeeper(cdc, didStoreKey, getSubspace(didtypes.ModuleName, paramsKeeper), accountKeeper, bankKeeper)
	capabilityKeeper := capabilitykeeper.NewKeeper(cdc, capabilityStoreKey, memStoreKeys[capabilitytypes.MemStoreKey])

	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	portKeeper := portkeeper.NewKeeper(scopedIBCKeeper)

	scopedResourceKeeper := capabilityKeeper.ScopeToModule(types.ModuleName)
	resourceKeeper := keeper.NewKeeper(cdc, resourceStoreKey, getSubspace(types.ModuleName, paramsKeeper), &portKeeper, scopedResourceKeeper)

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
		IBCModule:           ibcModule,
	}

	setup.Keeper.SetDidNamespace(&ctx, didsetup.DidNamespace)

	return setup
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key storetypes.StoreKey, tkey storetypes.StoreKey) paramskeeper.Keeper {
	// create keeper
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// set params subspaces
	paramsKeeper.Subspace(didtypes.ModuleName)
	paramsKeeper.Subspace(types.ModuleName).WithKeyTable(types.ParamKeyTable())

	return paramsKeeper
}

func getSubspace(moduleName string, paramsKeeper paramskeeper.Keeper) paramstypes.Subspace {
	subspace, _ := paramsKeeper.GetSubspace(moduleName)
	return subspace
}
