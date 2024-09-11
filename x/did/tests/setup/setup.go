package setup

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/cheqd/cheqd-node/app"
	"github.com/cheqd/cheqd-node/x/did/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"

	"github.com/cheqd/cheqd-node/x/did/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type TestSetup struct {
	Cdc codec.Codec

	SdkCtx sdk.Context
	StdCtx context.Context

	Keeper      keeper.Keeper
	MsgServer   types.MsgServer
	QueryServer types.QueryServer
}

func Setup() TestSetup {
	// Init Codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	Cdc := codec.NewProtoCodec(ir)
	banktypes.RegisterInterfaces(ir)
	authtypes.RegisterInterfaces(ir)
	aminoCdc := codec.NewLegacyAmino()

	// Init KVStore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	dbStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init ParamsKeeper KVStore
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	paramsTStoreKey := sdk.NewTransientStoreKey(paramstypes.TStoreKey)

	// Mount account and bank stores
	authStoreKey := sdk.NewKVStoreKey(authtypes.StoreKey)
	bankStoreKey := sdk.NewKVStoreKey(banktypes.StoreKey)
	dbStore.MountStoreWithDB(authStoreKey, storetypes.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(bankStoreKey, storetypes.StoreTypeIAVL, nil)

	paramsKeeper := initParamsKeeper(Cdc, aminoCdc, paramsStoreKey, paramsTStoreKey)

	// Initialize accountKeeper
	maccPerms := map[string][]string{
		authtypes.FeeCollectorName: nil,
		types.ModuleName:           {authtypes.Minter, authtypes.Burner},
	}

	accountKeeper := authkeeper.NewAccountKeeper(
		Cdc,
		authStoreKey,
		authtypes.ProtoBaseAccount,
		maccPerms,
		app.AccountAddressPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	bankKeeper := bankkeeper.NewBaseKeeper(
		Cdc,
		bankStoreKey,
		accountKeeper,
		nil,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Init Keepers
	newKeeper := keeper.NewKeeper(Cdc, storeKey, getSubspace(types.ModuleName, paramsKeeper), bankKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String())

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

	setup := TestSetup{
		Cdc: Cdc,

		SdkCtx: ctx,
		StdCtx: sdk.WrapSDKContext(ctx),

		Keeper:      *newKeeper,
		MsgServer:   msgServer,
		QueryServer: queryServer,
	}

	setup.Keeper.SetDidNamespace(&ctx, DidNamespace)
	return setup
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key storetypes.StoreKey, tkey storetypes.StoreKey) paramskeeper.Keeper {
	// create keeper
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// set params subspaces
	paramsKeeper.Subspace(types.ModuleName)

	return paramsKeeper
}

func getSubspace(moduleName string, paramsKeeper paramskeeper.Keeper) paramstypes.Subspace {
	subspace, _ := paramsKeeper.GetSubspace(moduleName)
	return subspace
}
