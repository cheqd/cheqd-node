package tests

import (
	"crypto/ed25519"
	"log"
	"testing"

	"crypto/rand"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	ptypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateDid(t *testing.T) {
	setup := Setup()

	_, did, _ := setup.InitDid("did:cheqd:test:alice")

	// query Did
	state, _ := setup.Keeper.GetDid(&setup.Ctx, did.Id)
	receivedDid := state.GetDid()

	// check
	require.Equal(t, did.Id, receivedDid.Id)
	require.Equal(t, did.Controller, receivedDid.Controller)
	require.Equal(t, did.VerificationMethod, receivedDid.VerificationMethod)
	require.Equal(t, did.Authentication, receivedDid.Authentication)
	require.Equal(t, did.AssertionMethod, receivedDid.AssertionMethod)
	require.Equal(t, did.CapabilityInvocation, receivedDid.CapabilityInvocation)
	require.Equal(t, did.CapabilityDelegation, receivedDid.CapabilityDelegation)
	require.Equal(t, did.KeyAgreement, receivedDid.KeyAgreement)
	require.Equal(t, did.AlsoKnownAs, receivedDid.AlsoKnownAs)
	require.Equal(t, did.Service, receivedDid.Service)
	require.Equal(t, did.Context, receivedDid.Context)
}

func TestHandler_UpdateDid(t *testing.T) {
	setup := Setup()

	//Init did
	keys, did, _ := setup.InitDid("did:cheqd:test:alice")

	// query Did
	state, _ := setup.Keeper.GetDid(&setup.Ctx, did.Id)

	//Init priv key
	newPubKey, _, _ := ed25519.GenerateKey(rand.Reader)

	didMsgUpdate := setup.UpdateDid(state.GetDid(), newPubKey, state.Metadata.VersionId)
	dataUpdate, _ := ptypes.NewAnyWithValue(didMsgUpdate)
	resultUpdate, err := setup.Handler(setup.Ctx, setup.WrapRequest(dataUpdate, keys, map[string]string{}))
	if err != nil {

	}

	didUpdated := types.MsgUpdateDidResponse{}
	errUpdate := didUpdated.Unmarshal(resultUpdate.Data)

	if errUpdate != nil {
		log.Fatal(errUpdate)
	}

	// query Did
	updatedState, _ := setup.Keeper.GetDid(&setup.Ctx, did.Id)
	receivedUpdatedDid := updatedState.GetDid()

	// check
	require.Equal(t, didUpdated.Id, receivedUpdatedDid.Id)
	require.Equal(t, didMsgUpdate.Controller, receivedUpdatedDid.Controller)
	require.Equal(t, didMsgUpdate.VerificationMethod, receivedUpdatedDid.VerificationMethod)
	require.Equal(t, didMsgUpdate.Authentication, receivedUpdatedDid.Authentication)
	require.Equal(t, didMsgUpdate.AssertionMethod, receivedUpdatedDid.AssertionMethod)
	require.Equal(t, didMsgUpdate.CapabilityInvocation, receivedUpdatedDid.CapabilityInvocation)
	require.Equal(t, didMsgUpdate.CapabilityDelegation, receivedUpdatedDid.CapabilityDelegation)
	require.Equal(t, didMsgUpdate.KeyAgreement, receivedUpdatedDid.KeyAgreement)
	require.Equal(t, didMsgUpdate.AlsoKnownAs, receivedUpdatedDid.AlsoKnownAs)
	require.Equal(t, didMsgUpdate.Service, receivedUpdatedDid.Service)
	require.NotEqual(t, state.GetDid().VerificationMethod, receivedUpdatedDid.VerificationMethod)
}

func TestHandler_CreateSchema(t *testing.T) {
	setup := Setup()

	keys, _, _ := setup.InitDid("did:cheqd:test:alice")
	msg := setup.CreateSchema()

	data, _ := ptypes.NewAnyWithValue(msg)
	result, _ := setup.Handler(setup.Ctx, setup.WrapRequest(data, keys, make(map[string]string)))

	schema := types.MsgCreateSchemaResponse{}
	err := schema.Unmarshal(result.Data)

	if err != nil {
		log.Fatal(err)
	}

	// query Did
	state, _ := setup.Keeper.GetSchema(setup.Ctx, schema.Id)
	receivedSchema := state.GetSchema()

	require.Equal(t, schema.Id, receivedSchema.Id)
	require.Equal(t, msg.Type, receivedSchema.Type)
	require.Equal(t, msg.Name, receivedSchema.Name)
	require.Equal(t, msg.Version, receivedSchema.Version)
	require.Equal(t, msg.AttrNames, receivedSchema.AttrNames)
	require.Equal(t, msg.Controller, receivedSchema.Controller)
}

func TestHandler_CreateCredDef(t *testing.T) {
	setup := Setup()

	keys, _, _ := setup.InitDid("did:cheqd:test:alice")
	msg := setup.CreateCredDef()

	data, _ := ptypes.NewAnyWithValue(msg)
	result, _ := setup.Handler(setup.Ctx, setup.WrapRequest(data, keys, make(map[string]string)))

	credDef := types.MsgCreateCredDefResponse{}
	err := credDef.Unmarshal(result.Data)

	if err != nil {
		log.Fatal(err)
	}

	// query Did
	state, _ := setup.Keeper.GetCredDef(setup.Ctx, credDef.Id)
	receivedCredDef := state.GetCredDef()

	expectedValue := msg.Value.(*types.MsgCreateCredDef_ClType)
	actualValue := receivedCredDef.Value.(*types.CredDef_ClType)

	require.Equal(t, credDef.Id, receivedCredDef.Id)
	require.Equal(t, expectedValue.ClType, actualValue.ClType)
	require.Equal(t, msg.SchemaId, receivedCredDef.SchemaId)
	require.Equal(t, msg.Tag, receivedCredDef.Tag)
	require.Equal(t, msg.Type, receivedCredDef.Type)
	require.Equal(t, msg.Controller, receivedCredDef.Controller)
}

func TestHandler_DidDocAlreadyExists(t *testing.T) {
	setup := Setup()

	keys, _, _ := setup.InitDid("did:cheqd:test:alice")
	_, _, err := setup.InitDid("did:cheqd:test:alice")

	require.Error(t, err)
	require.Equal(t, "DID DOC already exists for DID did:cheqd:test:alice: DID Doc exists", err.Error())

	credDefMsg := setup.CreateCredDef()
	data, _ := ptypes.NewAnyWithValue(credDefMsg)
	_, _ = setup.Handler(setup.Ctx, setup.WrapRequest(data, keys, make(map[string]string)))
	_, err = setup.Handler(setup.Ctx, setup.WrapRequest(data, keys, make(map[string]string)))

	require.Error(t, err)
	require.Equal(t, "DID DOC already exists for CredDef did:cheqd:test:cred-def-1: DID Doc exists", err.Error())

	schemaMsg := setup.CreateSchema()
	data, _ = ptypes.NewAnyWithValue(schemaMsg)
	_, _ = setup.Handler(setup.Ctx, setup.WrapRequest(data, keys, make(map[string]string)))
	_, err = setup.Handler(setup.Ctx, setup.WrapRequest(data, keys, make(map[string]string)))

	require.Error(t, err)
	require.Equal(t, "DID DOC already exists for Schema did:cheqd:test:schema-1: DID Doc exists", err.Error())
}
