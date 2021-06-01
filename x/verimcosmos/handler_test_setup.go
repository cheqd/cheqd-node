package verimcosmos

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/verim-id/verim-cosmos/app/params"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/require"
	_ "github.com/tendermint/tendermint/abci/types"
	"github.com/verim-id/verim-cosmos/x/verimcosmos/keeper"
	"github.com/verim-id/verim-cosmos/x/verimcosmos/types"
)

type TestSetup struct {
	Cdc       codec.Marshaler
	Ctx       sdk.Context
	NymKeeper keeper.Keeper
	Handler   sdk.Handler
	Querier   sdk.Querier
	Vendor    sdk.AccAddress
}

func Setup() TestSetup {
	// Init Codec
	//encodingConfig := app.MakeEncodingConfig()
	encodingConfig := params.MakeEncodingConfig()
	cdc := encodingConfig.Marshaler
	//sdk.RegisterCodec(cdc)
	//codec.RegisterCrypto(cdc)

	// Init KVSore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)
	//
	nymKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := sdk.NewKVStoreKey(types.MemStoreKey)
	dbStore.MountStoreWithDB(nymKey, sdk.StoreTypeIAVL, nil)
	//
	//authKey := sdk.NewKVStoreKey(auth.StoreKey)
	//dbStore.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init Keepers
	//nymKeeper := keeper.Keeper{}
	nymKeeper := keeper.NewKeeper(cdc, nymKey, memStoreKey)

	// Create context
	blockTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00.000Z")
	ctx := sdk.NewContext(dbStore, tmproto.Header{ChainID: "verim-cosmos", Time: blockTime}, false, log.NewNopLogger())

	// Create Handler and Querier
	//querier := keeper.NewQuerier(nymKeeper, leg)
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

func TestMsgUpdateNym(id uint64) types.MsgUpdateNym {
	return types.MsgUpdateNym{
		Creator: "creator_address",
		Id:      id,
		Alias:   "test_alias",
		Verkey:  "test_verkey",
		Did:     "test_did",
		Role:    "test_role",
	}
}

func TestMsgDeleteNym(id uint64) types.MsgDeleteNym {
	return types.MsgDeleteNym{
		Creator: "creator_address",
		Id:      id,
	}
}

func TestQueryGetNym(id uint64) types.QueryGetNymRequest {
	return types.QueryGetNymRequest{
		Id: id,
	}
}
