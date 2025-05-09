package setup

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/cheqd/cheqd-node/x/did/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"

	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	appparams "github.com/cheqd/cheqd-node/app"
	"github.com/cheqd/cheqd-node/x/did/keeper"
	"github.com/cosmos/cosmos-sdk/runtime"
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
)

type TestSetup struct {
	Cdc codec.Codec

	SdkCtx sdk.Context
	StdCtx context.Context

	Keeper        keeper.Keeper
	MsgServer     types.MsgServer
	QueryServer   types.QueryServer
	BankKeeper    bankkeeper.Keeper
	AccountKeeper authkeeper.AccountKeeper
}

func Setup() TestSetup {
	// Init Codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	authtypes.RegisterInterfaces(ir)
	banktypes.RegisterInterfaces(ir)
	stakingtypes.RegisterInterfaces(ir)

	Cdc := codec.NewProtoCodec(ir)
	aminoCdc := codec.NewLegacyAmino()
	// Init KVStore
	db := dbm.NewMemDB()

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		types.StoreKey,
		stakingtypes.StoreKey,
	)
	dbStore := store.NewCommitMultiStore(db, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	dbStore.MountStoreWithDB(keys[types.StoreKey], storetypes.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(keys[authtypes.StoreKey], storetypes.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(keys[banktypes.StoreKey], storetypes.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(keys[stakingtypes.StoreKey], storetypes.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	maccPerms := map[string][]string{
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		authtypes.FeeCollectorName:     nil,
		types.ModuleName:               {authtypes.Minter, authtypes.Burner},
		govtypes.ModuleName:            {authtypes.Burner, authtypes.Minter},
	}

	// Init ParamsKeeper KVStore
	paramsStoreKey := storetypes.NewKVStoreKey(paramstypes.StoreKey)
	paramsTStoreKey := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)

	paramsKeeper := initParamsKeeper(Cdc, aminoCdc, paramsStoreKey, paramsTStoreKey)
	accountKeeper := authkeeper.NewAccountKeeper(Cdc, runtime.NewKVStoreService(keys[authtypes.StoreKey]), authtypes.ProtoBaseAccount, maccPerms, authcodec.NewBech32Codec("cheqd"), "cheqd", authtypes.NewModuleAddress(govtypes.ModuleName).String())
	bankKeeper := bankkeeper.NewBaseKeeper(
		Cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		nil,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		log.NewNopLogger(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(Cdc, runtime.NewKVStoreService(keys[stakingtypes.StoreKey]), accountKeeper, bankKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(), authcodec.NewBech32Codec(appparams.ValidatorAddressPrefix), authcodec.NewBech32Codec(appparams.ConsNodeAddressPrefix))
	newKeeper := keeper.NewKeeper(Cdc, runtime.NewKVStoreService(keys[types.StoreKey]), getSubspace(types.ModuleName, paramsKeeper), accountKeeper, bankKeeper, stakingKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	// Create Tx
	txBytes := make([]byte, 28)
	_, _ = rand.Read(txBytes)

	// Create context
	blockTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00.000Z")
	ctx := sdk.NewContext(dbStore,
		tmproto.Header{ChainID: "test", Time: blockTime},
		false, log.NewNopLogger()).WithTxBytes(txBytes)

	msgServer := keeper.NewMsgServer(*newKeeper)
	queryServer := keeper.NewQueryServer(*newKeeper)
	goCtx := ctx

	params := stakingtypes.DefaultParams()
	params.BondDenom = "ncheq"
	err := stakingKeeper.SetParams(goCtx, params)
	if err != nil {
		panic("error while setting up the params")
	}
	setup := TestSetup{
		Cdc: Cdc,

		SdkCtx: ctx,
		StdCtx: goCtx,

		Keeper:        *newKeeper,
		MsgServer:     msgServer,
		QueryServer:   queryServer,
		BankKeeper:    bankKeeper,
		AccountKeeper: accountKeeper,
	}
	err = setup.Keeper.SetDidNamespace(goCtx, DidNamespace)
	if err != nil {
		panic(err)
	}
	return setup
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key storetypes.StoreKey, tkey storetypes.StoreKey) paramskeeper.Keeper {
	// create keeper
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// set params subspaces
	paramsKeeper.Subspace(types.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)

	return paramsKeeper
}

func getSubspace(moduleName string, paramsKeeper paramskeeper.Keeper) paramstypes.Subspace {
	subspace, _ := paramsKeeper.GetSubspace(moduleName)
	return subspace
}
