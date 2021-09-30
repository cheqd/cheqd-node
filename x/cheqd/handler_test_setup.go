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
	// ptypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type TestSetup struct {
	Cdc     codec.Marshaler
	Ctx     sdk.Context
	Keeper  keeper.Keeper
	Handler sdk.Handler
}

func Setup() TestSetup {
	// Init Codec
	encodingConfig := params.MakeEncodingConfig()
	cdc := encodingConfig.Marshaler

	// Init KVSore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := sdk.NewKVStoreKey(types.MemStoreKey)
	dbStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init Keepers
	keeper := keeper.NewKeeper(cdc, storeKey, memStoreKey)

	// Create context
	blockTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00.000Z")
	ctx := sdk.NewContext(dbStore,
		tmproto.Header{ChainID: "cheqd-node", Time: blockTime},
		false, log.NewNopLogger())

	handler := NewHandler(*keeper)

	setup := TestSetup{
		Cdc:     cdc,
		Ctx:     ctx,
		Keeper:  *keeper,
		Handler: handler,
	}

	return setup
}

func (s *TestSetup) CreateDid() *types.MsgCreateDid {
	VerificationMethod := types.VerificationMethod{
		Id:                 "12",
		Type:               "Ed25519VerificationKey2020",
		Controller:         "Controller",
		PublicKeyMultibase: "21312",
	}

	Service := types.DidService{
		Id:              "1",
		Type:            "type",
		ServiceEndpoint: "endpoint",
	}

	return &types.MsgCreateDid{
		Id:                   "1",
		Controller:           []string{"controller"},
		VerificationMethod:   []*types.VerificationMethod{&VerificationMethod},
		Authentication:       []string{"Authentication"},
		AssertionMethod:      []string{"AssertionMethod"},
		CapabilityInvocation: []string{"CapabilityInvocation"},
		CapabilityDelegation: []string{"CapabilityDelegation"},
		KeyAgreement:         []string{"KeyAgreement"},
		AlsoKnownAs:          []string{"AlsoKnownAs"},
		Service:              []*types.DidService{&Service},
	}
}
