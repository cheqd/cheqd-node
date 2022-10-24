package setup

import (
	"crypto/rand"
	"time"

	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	cheqdsetup "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cheqd/cheqd-node/x/resource/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TestSetup struct {
	cheqdsetup.TestSetup

	ResourceKeeper      keeper.Keeper
	ResourceMsgServer   types.MsgServer
	ResourceQueryServer types.QueryServer
}

func Setup() TestSetup {
	// Init Codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	cheqdtypes.RegisterInterfaces(ir)
	cdc := codec.NewProtoCodec(ir)

	// Init KVSore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)

	cheqdStoreKey := sdk.NewKVStoreKey(cheqdtypes.StoreKey)
	resourceStoreKey := sdk.NewKVStoreKey(types.StoreKey)

	dbStore.MountStoreWithDB(cheqdStoreKey, sdk.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(resourceStoreKey, sdk.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init Keepers
	cheqdKeeper := cheqdkeeper.NewKeeper(cdc, cheqdStoreKey)
	resourceKeeper := keeper.NewKeeper(cdc, resourceStoreKey)

	// Create Tx
	txBytes := make([]byte, 28)
	_, _ = rand.Read(txBytes)

	// Create context
	blockTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00.000Z")
	ctx := sdk.NewContext(dbStore,
		tmproto.Header{ChainID: "test", Time: blockTime},
		false, log.NewNopLogger()).WithTxBytes(txBytes)

	// Init servers
	cheqdMsgServer := cheqdkeeper.NewMsgServer(*cheqdKeeper)
	cheqdQueryServer := cheqdkeeper.NewQueryServer(*cheqdKeeper)

	msgServer := keeper.NewMsgServer(*resourceKeeper, *cheqdKeeper)
	queryServer := keeper.NewQueryServer(*resourceKeeper, *cheqdKeeper)

	setup := TestSetup{
		TestSetup: cheqdsetup.TestSetup{
			Cdc: cdc,

			SdkCtx: ctx,
			StdCtx: sdk.WrapSDKContext(ctx),

			Keeper:      *cheqdKeeper,
			MsgServer:   cheqdMsgServer,
			QueryServer: cheqdQueryServer,
		},

		ResourceKeeper:      *resourceKeeper,
		ResourceMsgServer:   msgServer,
		ResourceQueryServer: queryServer,
	}

	setup.Keeper.SetDidNamespace(&ctx, cheqdsetup.DID_NAMESPACE)
	return setup
}
