package setup

import (
	"crypto/rand"
	"time"

	"github.com/cheqd/cheqd-node/app"
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"

	"cosmossdk.io/store/metrics"
	"github.com/cheqd/cheqd-node/x/resource"
	"github.com/cheqd/cheqd-node/x/resource/keeper"

	// sdk "github.com/cosmos/cosmos-sdk/types"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	portkeeper "github.com/cosmos/ibc-go/v8/modules/core/05-port/keeper"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

type TestSetup struct {
	didsetup.TestSetup

	ResourceKeeper      keeper.Keeper
	ResourceMsgServer   types.MsgServer
	ResourceQueryServer types.QueryServer
	IBCModule           resource.IBCModule
	// storeService        store.KVStoreService
}

func Setup() TestSetup {
	// Init Codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	authtypes.RegisterInterfaces(ir)
	banktypes.RegisterInterfaces(ir)
	didtypes.RegisterInterfaces(ir)
	banktypes.RegisterInterfaces(ir)
	authtypes.RegisterInterfaces(ir)
	cdc := codec.NewProtoCodec(ir)
	aminoCdc := codec.NewLegacyAmino()

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Init KVStore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

	keys := storetypes.NewKVStoreKeys(
		capabilitytypes.StoreKey,
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
	)

	maccPerms := map[string][]string{
		minttypes.ModuleName:           {authtypes.Minter},
		types.ModuleName:               {authtypes.Minter, authtypes.Burner},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	}
	// Mount did store
	didStoreKey := storetypes.NewKVStoreKey(didtypes.StoreKey)
	dbStore.MountStoreWithDB(didStoreKey, storetypes.StoreTypeIAVL, nil)

	dbStore.MountStoreWithDB(keys[authtypes.StoreKey], storetypes.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(keys[banktypes.StoreKey], storetypes.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(keys[stakingtypes.StoreKey], storetypes.StoreTypeIAVL, nil)

	// Mount resource store
	resourceStoreKey := storetypes.NewKVStoreKey(types.StoreKey)
	dbStore.MountStoreWithDB(resourceStoreKey, storetypes.StoreTypeIAVL, nil)

	// Mount capability store - required for ibc port tests
	capabilityStoreKey := storetypes.NewKVStoreKey(capabilitytypes.StoreKey)
	dbStore.MountStoreWithDB(capabilityStoreKey, storetypes.StoreTypeIAVL, nil)
	memStoreKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	dbStore.MountStoreWithDB(memStoreKeys[capabilitytypes.MemStoreKey], storetypes.StoreTypeMemory, nil)

	// Mount param store - required for ibc port tests with default genesis
	paramsStoreKey := storetypes.NewKVStoreKey(paramstypes.StoreKey)
	dbStore.MountStoreWithDB(paramsStoreKey, storetypes.StoreTypeIAVL, nil)
	paramsTStoreKey := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)
	dbStore.MountStoreWithDB(paramsTStoreKey, storetypes.StoreTypeTransient, nil)

	_ = dbStore.LoadLatestVersion()

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(app.AccountAddressPrefix),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authority,
	)

	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		nil,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		log.NewNopLogger(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(cdc, runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(app.AccountAddressPrefix),
		authcodec.NewBech32Codec(app.ConsNodeAddressPrefix))

	paramsKeeper := initParamsKeeper(cdc, aminoCdc, paramsStoreKey, paramsTStoreKey)

	didKeeper := didkeeper.NewKeeper(cdc, didStoreKey, getSubspace(didtypes.ModuleName, paramsKeeper), accountKeeper, bankKeeper, stakingKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	capabilityKeeper := capabilitykeeper.NewKeeper(cdc, capabilityStoreKey, memStoreKeys[capabilitytypes.MemStoreKey])

	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	portKeeper := portkeeper.NewKeeper(scopedIBCKeeper)

	scopedResourceKeeper := capabilityKeeper.ScopeToModule(types.ModuleName)
	resourceKeeper := keeper.NewKeeper(cdc, runtime.NewKVStoreService(keys[types.StoreKey]),
		&portKeeper,
		scopedResourceKeeper, authority)

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

	params := stakingtypes.DefaultParams()
	params.BondDenom = didtypes.BaseMinimalDenom
	err := stakingKeeper.SetParams(ctx, params)
	if err != nil {
		panic("error while setting up the params")
	}
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
