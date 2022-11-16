package setup

import (
	"context"
	"crypto/rand"
	"time"

	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didkeeperv1 "github.com/cheqd/cheqd-node/x/did/keeper/v1"
	didsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	resourcekeeperv1 "github.com/cheqd/cheqd-node/x/resource/keeper/v1"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

type TestSetup struct {
	Cdc codec.Codec

	SdkCtx sdk.Context
	StdCtx context.Context

	DidKeekerPrevious      didkeeperv1.Keeper
	ResourceKeeperprevious resourcekeeperv1.Keeper

	DidKeeper      didkeeper.Keeper
	DidMsgServer   didtypes.MsgServer
	DidQueryServer didtypes.QueryServer

	ResourceKeeper      resourcekeeper.Keeper
	ResourceMsgServer   resourcetypes.MsgServer
	ResourceQueryServer resourcetypes.QueryServer
}

func Setup() TestSetup {
	// Init Codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	didtypes.RegisterInterfaces(ir) // TODO: Is v1 needed?
	cdc := codec.NewProtoCodec(ir)

	// Init KVSore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)

	didStoreKey := sdk.NewKVStoreKey(didtypes.StoreKey)
	resourceStoreKey := sdk.NewKVStoreKey(types.StoreKey)

	dbStore.MountStoreWithDB(didStoreKey, storetypes.StoreTypeIAVL, nil)
	dbStore.MountStoreWithDB(resourceStoreKey, storetypes.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init previous keepers
	didKeeperPrevious := didkeeperv1.NewKeeper(cdc, didStoreKey)
	resourceKeeperPrevious := resourcekeeperv1.NewKeeper(cdc, resourceStoreKey)

	// Init Keepers
	didKeeper := didkeeper.NewKeeper(cdc, didStoreKey)
	resourceKeeper := resourcekeeper.NewKeeper(cdc, resourceStoreKey)

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
		Cdc: cdc,

		SdkCtx: ctx,
		StdCtx: sdk.WrapSDKContext(ctx),

		DidKeekerPrevious:      *didKeeperPrevious,
		ResourceKeeperprevious: *resourceKeeperPrevious,

		DidKeeper:      *didKeeper,
		DidMsgServer:   didMsgServer,
		DidQueryServer: didQueryServer,

		ResourceKeeper:      *resourceKeeper,
		ResourceMsgServer:   resourceMsgServer,
		ResourceQueryServer: resourceQueryServer,
	}

	setup.DidKeeper.SetDidNamespace(&ctx, didsetup.DID_NAMESPACE) // TODO: Think about it
	return setup
}
