package setup

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cheqd/cheqd-node/x/cheqd/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	cdc := codec.NewProtoCodec(ir)

	// Init KVSore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	dbStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init Keepers
	newKeeper := keeper.NewKeeper(cdc, storeKey)

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
		Cdc: cdc,

		SdkCtx: ctx,
		StdCtx: sdk.WrapSDKContext(ctx),

		Keeper:      *newKeeper,
		MsgServer:   msgServer,
		QueryServer: queryServer,
	}

	setup.Keeper.SetDidNamespace(&ctx, DID_NAMESPACE)
	return setup
}
