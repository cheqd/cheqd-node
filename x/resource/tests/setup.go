package tests

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"testing"
	"time"

	"github.com/cheqd/cheqd-node/x/cheqd"
	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"

	cheqdtests "github.com/cheqd/cheqd-node/x/cheqd/tests"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/resource"
	"github.com/cheqd/cheqd-node/x/resource/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cheqd/cheqd-node/x/resource/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TestSetup struct {
	cheqdtests.TestSetup

	ResourceKeeper  keeper.Keeper
	ResourceHandler sdk.Handler
	QueryServer     types.QueryServer
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
	queryServer := keeper.NewQueryServer(*resourceKeeper, *cheqdKeeper)

	// Create Tx
	txBytes := make([]byte, 28)
	_, _ = rand.Read(txBytes)

	// Create context
	blockTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00.000Z")
	ctx := sdk.NewContext(dbStore,
		tmproto.Header{ChainID: "test", Time: blockTime},
		false, log.NewNopLogger()).WithTxBytes(txBytes)

	cheqdHandler := cheqd.NewHandler(*cheqdKeeper)
	resourceHandler := resource.NewHandler(*resourceKeeper, *cheqdKeeper)

	setup := TestSetup{
		TestSetup: cheqdtests.TestSetup{
			Cdc:     cdc,
			Ctx:     ctx,
			Keeper:  *cheqdKeeper,
			Handler: cheqdHandler,
		},
		ResourceKeeper:  *resourceKeeper,
		ResourceHandler: resourceHandler,
		QueryServer:     queryServer,
	}

	setup.Keeper.SetDidNamespace(&ctx, "test")
	return setup
}

func GenerateCreateResourcePayload(resource types.Resource) *types.MsgCreateResourcePayload {
	return &types.MsgCreateResourcePayload{
		CollectionId: resource.Header.CollectionId,
		Id:           resource.Header.Id,
		Name:         resource.Header.Name,
		ResourceType: resource.Header.ResourceType,
		Data:         resource.Data,
	}
}

func (s *TestSetup) WrapCreateRequest(payload *types.MsgCreateResourcePayload, keys map[string]ed25519.PrivateKey) *types.MsgCreateResource {
	var signatures []*cheqdtypes.SignInfo
	signingInput := payload.GetSignBytes()

	for privKeyId, privKey := range keys {
		signature := base64.StdEncoding.EncodeToString(ed25519.Sign(privKey, signingInput))
		signatures = append(signatures, &cheqdtypes.SignInfo{
			VerificationMethodId: privKeyId,
			Signature:            signature,
		})
	}

	return &types.MsgCreateResource{
		Payload:    payload,
		Signatures: signatures,
	}
}

func (s *TestSetup) SendCreateResource(msg *types.MsgCreateResourcePayload, keys map[string]ed25519.PrivateKey) (*types.Resource, error) {
	_, err := s.ResourceHandler(s.Ctx, s.WrapCreateRequest(msg, keys))
	if err != nil {
		return nil, err
	}

	created, _ := s.ResourceKeeper.GetResource(&s.Ctx, msg.CollectionId, msg.Id)
	return &created, nil
}

func InitEnv(t *testing.T, publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) TestSetup {
	resourceSetup := Setup()

	didDoc := resourceSetup.CreateDid(publicKey, ExistingDID)
	_, err := resourceSetup.SendCreateDid(didDoc, map[string]ed25519.PrivateKey{ExistingDIDKey: privateKey})
	require.NoError(t, err)

	resourcePayload := GenerateCreateResourcePayload(ExistingResource())
	_, err = resourceSetup.SendCreateResource(resourcePayload, map[string]ed25519.PrivateKey{ExistingDIDKey: privateKey})
	require.NoError(t, err)

	return resourceSetup
}

func GenerateTestKeys() map[string]cheqdtests.KeyPair {
	return map[string]cheqdtests.KeyPair{
		ExistingDIDKey: cheqdtests.GenerateKeyPair(),
	}
}

func CreateChecksum(data []byte) []byte {
	h := sha256.New()
	h.Write(data) 
	return h.Sum(nil)
}
