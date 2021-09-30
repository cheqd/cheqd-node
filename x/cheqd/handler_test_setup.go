package cheqd

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"time"

	"github.com/cheqd/cheqd-node/app/params"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	ptypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	newKeeper := keeper.NewKeeper(cdc, storeKey, memStoreKey)

	// Create context
	blockTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00.000Z")
	ctx := sdk.NewContext(dbStore,
		tmproto.Header{ChainID: "cheqd-node", Time: blockTime},
		false, log.NewNopLogger())

	handler := NewHandler(*newKeeper)

	setup := TestSetup{
		Cdc:     cdc,
		Ctx:     ctx,
		Keeper:  *newKeeper,
		Handler: handler,
	}

	return setup
}

func (s *TestSetup) CreateDid(pubKey ed25519.PublicKey) *types.MsgCreateDid {
	PublicKeyMultibase := "z" + base58.Encode(pubKey)

	VerificationMethod := types.VerificationMethod{
		Id:                 "did:cheqd:test:alice#key-1",
		Type:               "Ed25519VerificationKey2020",
		Controller:         "Controller",
		PublicKeyMultibase: PublicKeyMultibase,
	}

	Service := types.DidService{
		Id:              "1",
		Type:            "type",
		ServiceEndpoint: "endpoint",
	}

	return &types.MsgCreateDid{
		Id:                   "did:cheqd:test:alice",
		Controller:           []string{"controller"},
		VerificationMethod:   []*types.VerificationMethod{&VerificationMethod},
		Authentication:       []string{"did:cheqd:test:alice#key-1"},
		AssertionMethod:      []string{"AssertionMethod"},
		CapabilityInvocation: []string{"CapabilityInvocation"},
		CapabilityDelegation: []string{"CapabilityDelegation"},
		KeyAgreement:         []string{"KeyAgreement"},
		AlsoKnownAs:          []string{"AlsoKnownAs"},
		Service:              []*types.DidService{&Service},
	}
}

func (s *TestSetup) WrapRequest(privKey ed25519.PrivateKey, data *ptypes.Any, metadata map[string]string) *types.MsgWriteRequest {
	metadataBytes, _ := json.Marshal(&metadata)
	dataBytes := data.Value

	signingInput := base64.StdEncoding.EncodeToString(metadataBytes) + "." + base64.StdEncoding.EncodeToString(dataBytes)
	signingInputBytes := []byte(base64.StdEncoding.EncodeToString([]byte(signingInput)))
	signature := base64.StdEncoding.EncodeToString(ed25519.Sign(privKey, signingInputBytes))

	return &types.MsgWriteRequest{
		Data:       data,
		Metadata:   metadata,
		Signatures: map[string]string{"did:cheqd:test:alice#key-1": signature},
	}
}
