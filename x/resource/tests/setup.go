package tests

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	// "github.com/btcsuite/btcutil/base58"
	cheqdtests "github.com/cheqd/cheqd-node/x/cheqd/tests"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/resource"
	"github.com/cheqd/cheqd-node/x/resource/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"

	// "github.com/multiformats/go-multibase"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cheqd/cheqd-node/x/resource/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TestSetup struct {
	Cdc     codec.Codec
	Ctx     sdk.Context
	Keeper  keeper.Keeper
	Handler sdk.Handler
}


func Setup() TestSetup {

	cheqdtests.Setup()

	// Init Codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	cdc := codec.NewProtoCodec(ir)

	// Init KVSore
	db := dbm.NewMemDB()

	dbStore := store.NewCommitMultiStore(db)
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	dbStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, nil)

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

	handler := resource.NewHandler(*newKeeper)

	setup := TestSetup{
		Cdc:     cdc,
		Ctx:     ctx,
		Keeper:  *newKeeper,
		Handler: handler,
	}
	return setup
}

func GenerateCreateResourcePayload(resource types.Resource) *types.MsgCreateResourcePayload {
	return &types.MsgCreateResourcePayload{
		CollectionId: resource.CollectionId,
		Id:           resource.Id,
		Name:         resource.Name,
		ResourceType: resource.ResourceType,
		MimeType:     resource.MimeType,
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
	_, err := s.Handler(s.Ctx, s.WrapCreateRequest(msg, keys))
	if err != nil {
		return nil, err
	}

	created, _ := s.Keeper.GetResource(&s.Ctx, msg.CollectionId, msg.Id)
	return &created, nil
}


func InitEnv(t *testing.T, publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) (TestSetup, cheqdtests.TestSetup) {
	resourceSetup := Setup()
	didSetup := cheqdtests.Setup()

	didDoc := didSetup.CreateDid(publicKey, ExistingDID)
	_, err := didSetup.SendCreateDid(didDoc, map[string]ed25519.PrivateKey{ExistingDIDKey: privateKey})
	require.NoError(t, err)

	resourcePayload := GenerateCreateResourcePayload(ExistingResource())
	_, err = resourceSetup.SendCreateResource(resourcePayload, map[string]ed25519.PrivateKey{ExistingDIDKey: privateKey})
	require.NoError(t, err)

	return resourceSetup, didSetup
}

func GenerateTestKeys() map[string]cheqdtests.KeyPair {
	return map[string]cheqdtests.KeyPair{
		ExistingDIDKey:    cheqdtests.GenerateKeyPair(),
	}
}
