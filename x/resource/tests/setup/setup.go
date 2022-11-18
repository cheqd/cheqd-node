package setup

import (
	"crypto/rand"
	"time"

	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cheqd/cheqd-node/x/resource/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TestSetup struct {
	didsetup.TestSetup

	ResourceKeeper      keeper.Keeper
	ResourceMsgServer   types.MsgServer
	ResourceQueryServer types.QueryServer
	ResourceStoreKey    *storetypes.KVStoreKey
}

func Setup() TestSetup {
	// Init Codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	didtypes.RegisterInterfaces(ir)
	cdc := codec.NewProtoCodec(ir)

	// Init KVSore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)

	didStoreKey := sdk.NewKVStoreKey(didtypes.StoreKey)
	resourceStoreKey := sdk.NewKVStoreKey(types.StoreKey)

	dbStore.MountStoreWithDB(didStoreKey, storetypes.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(resourceStoreKey, storetypes.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init Keepers
	didKeeper := didkeeper.NewKeeper(cdc, didStoreKey)
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
			DidStoreKey: didStoreKey,
		},

		ResourceKeeper:      *resourceKeeper,
		ResourceMsgServer:   msgServer,
		ResourceQueryServer: queryServer,
		ResourceStoreKey:    resourceStoreKey,
	}

	setup.Keeper.SetDidNamespace(&ctx, didsetup.DID_NAMESPACE)
	return setup
}
