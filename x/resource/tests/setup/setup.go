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
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cheqd/cheqd-node/x/resource"
	"github.com/cheqd/cheqd-node/x/resource/keeper"

	storetypes "cosmossdk.io/store/types"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
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

	oraclekeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
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
	authtypes.RegisterInterfaces(ir)
	banktypes.RegisterInterfaces(ir)
	didtypes.RegisterInterfaces(ir)
	banktypes.RegisterInterfaces(ir)
	authtypes.RegisterInterfaces(ir)
	distrtypes.RegisterInterfaces(ir)
	oracletypes.RegisterInterfaces(ir)
	cdc := codec.NewProtoCodec(ir)
	aminoCdc := codec.NewLegacyAmino()

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Init KVStore
	keys := storetypes.NewKVStoreKeys(
		capabilitytypes.StoreKey,
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		paramstypes.StoreKey,
		types.StoreKey,
		didtypes.StoreKey,
		distrtypes.StoreKey,
		oracletypes.StoreKey,
	)

	transientKeys := storetypes.NewTransientStoreKeys(
		paramstypes.TStoreKey,
	)

	memKeys := storetypes.NewMemoryStoreKeys(
		capabilitytypes.MemStoreKey,
	)

	ctx := sdktestutil.DefaultContextWithKeys(keys, transientKeys, memKeys)

	maccPerms := map[string][]string{
		minttypes.ModuleName:           {authtypes.Minter},
		types.ModuleName:               {authtypes.Minter, authtypes.Burner},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		distrtypes.ModuleName:          nil,
		oracletypes.ModuleName:         {authtypes.Minter},
	}

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
	distrKeeper := distrkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	paramsKeeper := initParamsKeeper(cdc, aminoCdc, keys[paramstypes.StoreKey], transientKeys[paramstypes.StoreKey])

	OracleKeeper := oraclekeeper.NewKeeper(
		cdc,
		keys[oracletypes.ModuleName],
		getSubspace(oracletypes.ModuleName, paramsKeeper),
		accountKeeper,
		bankKeeper,
		distrKeeper,
		stakingKeeper,
		distrtypes.ModuleName,
		true, // cast.ToBool(appOpts.Get("telemetry.enabled")),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	didKeeper := didkeeper.NewKeeper(cdc, runtime.NewKVStoreService(keys[didtypes.StoreKey]), getSubspace(didtypes.ModuleName, paramsKeeper), accountKeeper, bankKeeper, stakingKeeper, OracleKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	capabilityKeeper := capabilitykeeper.NewKeeper(cdc, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	portKeeper := portkeeper.NewKeeper(scopedIBCKeeper)

	scopedResourceKeeper := capabilityKeeper.ScopeToModule(types.ModuleName)
	resourceKeeper := keeper.NewKeeper(cdc, runtime.NewKVStoreService(keys[types.StoreKey]),
		getSubspace(types.ModuleName, paramsKeeper),
		&portKeeper,
		scopedResourceKeeper, authority)

	ibcModule := resource.NewIBCModule(*resourceKeeper)

	// Create Tx
	txBytes := make([]byte, 28)
	_, _ = rand.Read(txBytes)

	// Create context
	blockTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00.000Z")
	ctx = ctx.WithBlockTime(blockTime).WithTxBytes(txBytes)

	// Init servers
	didMsgServer := didkeeper.NewMsgServer(*didKeeper)
	didQueryServer := didkeeper.NewQueryServer(*didKeeper)

	msgServer := keeper.NewMsgServer(*resourceKeeper, *didKeeper, OracleKeeper)
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
			StdCtx: ctx,

			Keeper:      *didKeeper,
			MsgServer:   didMsgServer,
			QueryServer: didQueryServer,
		},

		ResourceKeeper:      *resourceKeeper,
		ResourceMsgServer:   msgServer,
		ResourceQueryServer: queryServer,
		IBCModule:           ibcModule,
	}
	err = setup.Keeper.SetDidNamespace(ctx, didsetup.DidNamespace)
	if err != nil {
		panic(err)
	}

	return setup
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key storetypes.StoreKey, tkey storetypes.StoreKey) paramskeeper.Keeper {
	// create keeper
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// set params subspaces
	paramsKeeper.Subspace(didtypes.ModuleName)
	paramsKeeper.Subspace(types.ModuleName).WithKeyTable(types.ParamKeyTable())
	paramsKeeper.Subspace(oracletypes.ModuleName).WithKeyTable(oracletypes.ParamKeyTable())

	return paramsKeeper
}

func getSubspace(moduleName string, paramsKeeper paramskeeper.Keeper) paramstypes.Subspace {
	subspace, _ := paramsKeeper.GetSubspace(moduleName)
	return subspace
}
