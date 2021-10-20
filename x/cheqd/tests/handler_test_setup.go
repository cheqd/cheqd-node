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
	memStoreKey := sdk.NewKVStoreKey(types.MemStoreKey)
	dbStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, nil)

	_ = dbStore.LoadLatestVersion()

	// Init Keepers
	newKeeper := keeper.NewKeeper(cdc, storeKey, memStoreKey)

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

	Service := types.DidService{
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
		Service:              []*types.DidService{&Service},
	}
}

func (s *TestSetup) UpdateDid(did *types.Did, pubKey ed25519.PublicKey, versionId string) *types.MsgUpdateDid {
	PublicKeyMultibase := "z" + base58.Encode(pubKey)

	VerificationMethod := types.VerificationMethod{
		Id:                 "did:cheqd:test:alice#key-2",
		Type:               "Ed25519VerificationKey2020",
		Controller:         "Controller",
		PublicKeyMultibase: PublicKeyMultibase,
	}

	return &types.MsgUpdateDid{
		Id:                   did.Id,
		Controller:           did.Controller,
		VerificationMethod:   []*types.VerificationMethod{did.VerificationMethod[0], &VerificationMethod},
		Authentication:       append(did.Authentication, "#key-2"),
		AssertionMethod:      did.AssertionMethod,
		CapabilityInvocation: did.CapabilityInvocation,
		CapabilityDelegation: did.CapabilityDelegation,
		KeyAgreement:         did.KeyAgreement,
		AlsoKnownAs:          did.AlsoKnownAs,
		Service:              did.Service,
		VersionId:            versionId,
		Context:              did.Context,
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

func (s *TestSetup) CreateSchema() *types.MsgCreateSchema {
	return &types.MsgCreateSchema{
		Id:         "did:cheqd:test:schema-1?service=CL-Schema",
		Type:       "CL-Schema",
		Name:       "name",
		Version:    "version",
		AttrNames:  []string{"attr1", "attr2"},
		Controller: []string{"did:cheqd:test:alice"},
	}
}

func (s *TestSetup) CreateCredDef() *types.MsgCreateCredDef {
	Value := types.MsgCreateCredDef_ClType{
		ClType: &types.CredDefValue{
			Primary:    map[string]*ptypes.Any{"first": nil},
			Revocation: map[string]*ptypes.Any{"second": nil},
		},
	}

	return &types.MsgCreateCredDef{
		Id:         "did:cheqd:test:cred-def-1?service=CL-CredDef",
		SchemaId:   "schema-1",
		Tag:        "tag",
		Type:       "CL-CredDef",
		Value:      &Value,
		Controller: []string{"did:cheqd:test:alice"},
	}
}

func (s *TestSetup) WrapRequest(data *ptypes.Any, keys map[string]ed25519.PrivateKey, metadata map[string]string) *types.MsgWriteRequest {
	result := types.MsgWriteRequest{
		Data:     data,
		Metadata: metadata,
	}

	signatures := make(map[string]string)
	signingInput, _ := keeper.BuildSigningInput(s.Cdc, &result)

	for privKeyId, privKey := range keys {
		signature := base64.StdEncoding.EncodeToString(ed25519.Sign(privKey, signingInput))
		signatures[privKeyId] = signature
	}

	return &types.MsgWriteRequest{
		Data:       data,
		Metadata:   metadata,
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

	result, err := s.Handler(s.Ctx, s.WrapRequest(data, keys, make(map[string]string)))
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
	data, err := ptypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	// query Did
	state, _ := s.Keeper.GetDid(&s.Ctx, msg.Id)
	msg.VersionId = state.Metadata.VersionId

	_, err = s.Handler(s.Ctx, s.WrapRequest(data, keys, map[string]string{}))
	if err != nil {
		return nil, err
	}

	did, _ := s.Keeper.GetDid(&s.Ctx, msg.Id)
	return did.GetDid(), nil
}

func (s *TestSetup) SendCreateDid(msg *types.MsgCreateDid, keys map[string]ed25519.PrivateKey) (*types.Did, error) {
	data, err := ptypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	_, err = s.Handler(s.Ctx, s.WrapRequest(data, keys, map[string]string{}))
	if err != nil {
		return nil, err
	}

	state, _ := s.Keeper.GetDid(&s.Ctx, msg.Id)
	return state.GetDid(), nil
}

func ConcatKeys(dst map[string]ed25519.PrivateKey, src map[string]ed25519.PrivateKey) map[string]ed25519.PrivateKey {
	for k, v := range src {
		dst[k] = v
	}

	return dst
}
