package setup

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/cheqd/cheqd-node/x/did/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"cosmossdk.io/store"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"

	storetypes "cosmossdk.io/store/types"
	"github.com/cheqd/cheqd-node/x/did/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		types.StoreKey,
		stakingtypes.StoreKey,
	)
	dbStore := store.NewCommitMultiStore(db)
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
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	paramsTStoreKey := sdk.NewTransientStoreKey(paramstypes.TStoreKey)

	paramsKeeper := initParamsKeeper(Cdc, aminoCdc, paramsStoreKey, paramsTStoreKey)
	accountKeeper := authkeeper.NewAccountKeeper(Cdc, keys[authtypes.StoreKey], authtypes.ProtoBaseAccount, maccPerms, "cheqd", authtypes.NewModuleAddress(govtypes.ModuleName).String())
	bankKeeper := bankkeeper.NewBaseKeeper(
		Cdc,
		keys[banktypes.StoreKey],
		accountKeeper,
		nil,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(Cdc, keys[stakingtypes.StoreKey], accountKeeper, bankKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	newKeeper := keeper.NewKeeper(Cdc, keys[types.StoreKey], getSubspace(types.ModuleName, paramsKeeper), accountKeeper, bankKeeper, stakingKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String())

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

	params := stakingtypes.DefaultParams()
	params.BondDenom = "ncheq"
	err := stakingKeeper.SetParams(ctx, params)
	if err != nil {
		panic("error while setting up the params")
	}
	setup := TestSetup{
		Cdc: Cdc,

		SdkCtx: ctx,
		StdCtx: sdk.WrapSDKContext(ctx),

		Keeper:        *newKeeper,
		MsgServer:     msgServer,
		QueryServer:   queryServer,
		BankKeeper:    bankKeeper,
		AccountKeeper: accountKeeper,
	}

	setup.Keeper.SetDidNamespace(&ctx, DidNamespace)
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
