package cheqd

import (
	"time"

	"github.com/cheqd/cheqd-node/app/params"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TestSetup struct {
	Cdc       codec.Marshaler
	Ctx       sdk.Context
	NymKeeper keeper.Keeper
	Handler   sdk.Handler
}

func Setup() TestSetup {
	// Init Codec
	encodingConfig := params.MakeEncodingConfig()
	cdc := encodingConfig.Marshaler

	// Init KVSore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)
	nymKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := sdk.NewKVStoreKey(types.MemStoreKey)
	dbStore.MountStoreWithDB(nymKey, sdk.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init Keepers
	nymKeeper := keeper.NewKeeper(cdc, nymKey, memStoreKey)

	// Create context
	blockTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00.000Z")
	ctx := sdk.NewContext(dbStore,
		tmproto.Header{ChainID: "cheqd-node", Time: blockTime},
		false, log.NewNopLogger())

	handler := NewHandler(*nymKeeper)

	setup := TestSetup{
		Cdc:       cdc,
		Ctx:       ctx,
		NymKeeper: *nymKeeper,
		Handler:   handler,
	}

	return setup
}

func TestMsgCreateNym() *types.MsgCreateNym {
	return types.NewMsgCreateNym(
		"creator_address",
		"test_alias",
		"test_verkey",
		"test_did",
		"test_role",
	)
}

func TestMsgUpdateNym(id uint64) *types.MsgUpdateNym {
	return types.NewMsgUpdateNym(
		"creator_address",
		id,
		"test_alias_new",
		"test_verkey_new",
		"test_did_new",
		"test_role_new",
	)
}

func TestMsgDeleteNym(id uint64) *types.MsgDeleteNym {
	return types.NewMsgDeleteNym(
		"creator_address",
		id,
	)
}

func TestQueryGetNym(id uint64) types.QueryGetNymRequest {
	return types.QueryGetNymRequest{
		Id: id,
	}
}
