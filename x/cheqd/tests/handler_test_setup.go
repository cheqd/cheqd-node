package tests

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd"
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

type KeyPair struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

type TestSetup struct {
	Cdc     codec.Codec
	Ctx     sdk.Context
	Keeper  keeper.Keeper
	Handler sdk.Handler
}

func Setup() TestSetup {
	// Init Codec
	encodingConfig := params.MakeEncodingConfig()
	cdc := encodingConfig.Codec

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

	handler := cheqd.NewHandler(*newKeeper)

	setup := TestSetup{
		Cdc:     cdc,
		Ctx:     ctx,
		Keeper:  *newKeeper,
		Handler: handler,
	}

	return setup
}

func (s *TestSetup) CreateDid(pubKey ed25519.PublicKey, did string) *types.MsgCreateDid {
	PublicKeyMultibase := "z" + base58.Encode(pubKey)

	VerificationMethod := types.VerificationMethod{
		Id:                 did + "#key-1",
		Type:               "Ed25519VerificationKey2020",
		Controller:         did,
		PublicKeyMultibase: PublicKeyMultibase,
	}

	Service := types.ServiceEndpoint{
		Id:              "#service-2",
		Type:            "DIDCommMessaging",
		ServiceEndpoint: "endpoint",
	}

	return &types.MsgCreateDid{
		Id:                   did,
		Controller:           nil,
		VerificationMethod:   []*types.VerificationMethod{&VerificationMethod},
		Authentication:       []string{did + "#key-1"},
		AssertionMethod:      []string{did + "#key-1"},
		CapabilityInvocation: []string{did + "#key-1"},
		CapabilityDelegation: []string{did + "#key-1"},
		KeyAgreement:         []string{did + "#key-1"},
		AlsoKnownAs:          []string{did + "#key-1"},
		Context:              []string{"Context"},
		Service:              []*types.ServiceEndpoint{&Service},
	}
}

func (s *TestSetup) CreateToUpdateDid(did *types.MsgCreateDid) *types.MsgUpdateDid {
	return &types.MsgUpdateDid{
		Id:                   did.Id,
		Controller:           did.Controller,
		VerificationMethod:   did.VerificationMethod,
		Authentication:       did.Authentication,
		AssertionMethod:      did.AssertionMethod,
		CapabilityInvocation: did.CapabilityInvocation,
		CapabilityDelegation: did.CapabilityDelegation,
		KeyAgreement:         did.KeyAgreement,
		AlsoKnownAs:          did.AlsoKnownAs,
		Service:              did.Service,
		Context:              did.Context,
	}
}

func (s *TestSetup) WrapRequest(data *ptypes.Any, keys map[string]ed25519.PrivateKey) *types.MsgWriteRequest {
	result := types.MsgWriteRequest{
		Data: data,
	}

	var signatures []*types.SignInfo
	signingInput := result.Data.Value

	for privKeyId, privKey := range keys {
		signature := base64.StdEncoding.EncodeToString(ed25519.Sign(privKey, signingInput))
		signatures = append(signatures, &types.SignInfo{
			VerificationMethodId: privKeyId,
			Signature:            signature,
		})
	}

	return &types.MsgWriteRequest{
		Data:       data,
		Signatures: signatures,
	}
}

func GenerateKeyPair() KeyPair {
	PublicKey, PrivateKey, _ := ed25519.GenerateKey(rand.Reader)
	return KeyPair{PrivateKey, PublicKey}
}

func (s *TestSetup) InitDid(did string) (map[string]ed25519.PrivateKey, *types.MsgCreateDid, error) {
	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)

	// add new Did
	didMsg := s.CreateDid(pubKey, did)
	data, err := ptypes.NewAnyWithValue(didMsg)
	if err != nil {
		return nil, nil, err
	}

	keyId := did + "#key-1"
	keys := map[string]ed25519.PrivateKey{keyId: privKey}

	result, err := s.Handler(s.Ctx, s.WrapRequest(data, keys))
	if err != nil {
		return nil, nil, err
	}

	didResponse := types.MsgCreateDidResponse{}
	if err := didResponse.Unmarshal(result.Data); err != nil {
		return nil, nil, err
	}

	return keys, didMsg, nil
}

func (s *TestSetup) SendUpdateDid(msg *types.MsgUpdateDid, keys map[string]ed25519.PrivateKey) (*types.Did, error) {
	// query Did
	state, _ := s.Keeper.GetDid(&s.Ctx, msg.Id)
	if len(msg.VersionId) == 0 {
		msg.VersionId = state.Metadata.VersionId
	}

	data, err := ptypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	_, err = s.Handler(s.Ctx, s.WrapRequest(data, keys))
	if err != nil {
		return nil, err
	}

	updated, _ := s.Keeper.GetDid(&s.Ctx, msg.Id)
	return updated.GetDid()
}

func (s *TestSetup) SendCreateDid(msg *types.MsgCreateDid, keys map[string]ed25519.PrivateKey) (*types.Did, error) {
	data, err := ptypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	_, err = s.Handler(s.Ctx, s.WrapRequest(data, keys))
	if err != nil {
		return nil, err
	}

	created, _ := s.Keeper.GetDid(&s.Ctx, msg.Id)
	return created.GetDid()
}

func ConcatKeys(dst map[string]ed25519.PrivateKey, src map[string]ed25519.PrivateKey) map[string]ed25519.PrivateKey {
	for k, v := range src {
		dst[k] = v
	}

	return dst
}

func (s TestSetup) CreatePreparedDID() map[string]KeyPair {
	prefilledDids := []struct {
		keys    map[string]KeyPair
		signers []string
		msg     *types.MsgCreateDid
	}{
		{
			keys: map[string]KeyPair{
				AliceKey1: GenerateKeyPair(),
				AliceKey2: GenerateKeyPair(),
			},
			signers: []string{AliceKey1},
			msg: &types.MsgCreateDid{
				Id:             AliceDID,
				Authentication: []string{AliceKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			keys: map[string]KeyPair{
				BobKey1: GenerateKeyPair(),
				BobKey2: GenerateKeyPair(),
				BobKey3: GenerateKeyPair(),
				BobKey4: GenerateKeyPair(),
			},
			signers: []string{BobKey2},
			msg: &types.MsgCreateDid{
				Id: BobDID,
				Authentication: []string{
					BobKey1,
					BobKey2,
					BobKey3,
				},
				CapabilityDelegation: []string{
					BobKey4,
				},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         BobKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey3,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey4,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
		{
			keys: map[string]KeyPair{
				CharlieKey1: GenerateKeyPair(),
				CharlieKey2: GenerateKeyPair(),
				CharlieKey3: GenerateKeyPair(),
			},
			signers: []string{CharlieKey2},
			msg: &types.MsgCreateDid{
				Id: CharlieDID,
				Authentication: []string{
					CharlieKey1,
					CharlieKey2,
					CharlieKey3,
				},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         CharlieKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         CharlieKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         CharlieKey3,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
	}

	keys := map[string]KeyPair{}

	for _, prefilled := range prefilledDids {
		msg := prefilled.msg

		for _, vm := range msg.VerificationMethod {
			vm.PublicKeyMultibase = "z" + base58.Encode(prefilled.keys[vm.Id].PublicKey)
		}

		signerKeys := map[string]ed25519.PrivateKey{}
		for _, signer := range prefilled.signers {
			signerKeys[signer] = prefilled.keys[signer].PrivateKey
		}

		for keyId, key := range prefilled.keys {
			keys[keyId] = key
		}

		_, _ = s.SendCreateDid(msg, signerKeys)
	}

	return keys
}
