package cheqd

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

	_, did, _ := setup.InitDid()

	// query Did
	receivedDid, _, _ := setup.Keeper.GetDid(setup.Ctx, did.Id)

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
}

func TestHandler_UpdateDid(t *testing.T) {
	setup := Setup()

	//Init did
	privKey, did, _ := setup.InitDid()

	// query Did
	receivedDid, didMetadata, _ := setup.Keeper.GetDid(setup.Ctx, did.Id)

	//Init priv key
	newPubKey, _, _ := ed25519.GenerateKey(rand.Reader)

	// add new Did
	metadata := map[string]string{
		"versionId": didMetadata.VersionId,
	}

	didMsgUpdate := setup.UpdateDid(receivedDid, newPubKey)
	dataUpdate, _ := ptypes.NewAnyWithValue(didMsgUpdate)
	resultUpdate, _ := setup.Handler(setup.Ctx, setup.WrapRequest(privKey, dataUpdate, metadata))
	didUpdated := types.MsgUpdateDidResponse{}
	errUpdate := didUpdated.Unmarshal(resultUpdate.Data)

	if errUpdate != nil {
		log.Fatal(errUpdate)
	}

	// query Did
	receivedUpdatedDid, _, _ := setup.Keeper.GetDid(setup.Ctx, did.Id)

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
	require.NotEqual(t, receivedDid.VerificationMethod, receivedUpdatedDid.VerificationMethod)
}

func TestHandler_UpdateDidInvalidSignature(t *testing.T) {
	setup := Setup()

	_, did, _ := setup.InitDid()

	// query Did
	receivedDid, _, _ := setup.Keeper.GetDid(setup.Ctx, did.Id)

	//Init priv key
	newPubKey, newPrivKey, _ := ed25519.GenerateKey(rand.Reader)

	// add new Did
	didMsgUpdate := setup.UpdateDid(receivedDid, newPubKey)
	dataUpdate, _ := ptypes.NewAnyWithValue(didMsgUpdate)
	_, err := setup.Handler(setup.Ctx, setup.WrapRequest(newPrivKey, dataUpdate, make(map[string]string)))
	require.Error(t, err)
	require.Equal(t, "Invalid signature: invalid signature detected", err.Error())
}

func TestHandler_CreateSchema(t *testing.T) {
	setup := Setup()

	privKey, _, _ := setup.InitDid()
	msg := setup.CreateSchema()

	data, _ := ptypes.NewAnyWithValue(msg)
	result, _ := setup.Handler(setup.Ctx, setup.WrapRequest(privKey, data, make(map[string]string)))

	schema := types.MsgCreateSchemaResponse{}
	err := schema.Unmarshal(result.Data)

	if err != nil {
		log.Fatal(err)
	}

	// query Did
	receivedSchema, _ := setup.Keeper.GetSchema(setup.Ctx, schema.Id)

	require.Equal(t, schema.Id, receivedSchema.Id)
	require.Equal(t, msg.Name, receivedSchema.Name)
	require.Equal(t, msg.Version, receivedSchema.Version)
	require.Equal(t, msg.AttrNames, receivedSchema.AttrNames)
}

func TestHandler_CreateCredDef(t *testing.T) {
	setup := Setup()

	privKey, _, _ := setup.InitDid()
	msg := setup.CreateCredDef()

	data, _ := ptypes.NewAnyWithValue(msg)
	result, _ := setup.Handler(setup.Ctx, setup.WrapRequest(privKey, data, make(map[string]string)))

	credDef := types.MsgCreateCredDefResponse{}
	err := credDef.Unmarshal(result.Data)

	if err != nil {
		log.Fatal(err)
	}

	// query Cred Def
	receivedCredDef, _ := setup.Keeper.GetCredDef(setup.Ctx, credDef.Id)

	expectedValue := msg.Value.(*types.MsgCreateCredDef_ClType)
	actualValue := receivedCredDef.Value.(*types.CredDef_ClType)

	require.Equal(t, credDef.Id, receivedCredDef.Id)
	require.Equal(t, expectedValue.ClType, actualValue.ClType)
	require.Equal(t, msg.SchemaId, receivedCredDef.SchemaId)
	require.Equal(t, msg.Tag, receivedCredDef.Tag)
	require.Equal(t, msg.SignatureType, receivedCredDef.SignatureType)
}

func TestHandler_CreateSchemaInvalidSignature(t *testing.T) {
	setup := Setup()

	_, privKey, _ := ed25519.GenerateKey(rand.Reader)
	msg := setup.CreateSchema()

	data, _ := ptypes.NewAnyWithValue(msg)
	_, err := setup.Handler(setup.Ctx, setup.WrapRequest(privKey, data, make(map[string]string)))

	require.Error(t, err)
	require.Equal(t, "Invalid signature: invalid signature detected", err.Error())
}

func TestHandler_DidDocAlreadyExists(t *testing.T) {
	setup := Setup()

	privKey, _, _ := setup.InitDid()
	_, _, err := setup.InitDid()

	require.Error(t, err)
	require.Equal(t, "DID DOC already exists for DID did:cheqd:test:alice: did doc exists", err.Error())

	credDefMsg := setup.CreateCredDef()
	data, _ := ptypes.NewAnyWithValue(credDefMsg)
	_, _ = setup.Handler(setup.Ctx, setup.WrapRequest(privKey, data, make(map[string]string)))
	_, err = setup.Handler(setup.Ctx, setup.WrapRequest(privKey, data, make(map[string]string)))

	require.Error(t, err)
	require.Equal(t, "DID DOC already exists for CredDef did:cheqd:test:cred-def-1: did doc exists", err.Error())

	schemaMsg := setup.CreateSchema()
	data, _ = ptypes.NewAnyWithValue(schemaMsg)
	_, _ = setup.Handler(setup.Ctx, setup.WrapRequest(privKey, data, make(map[string]string)))
	_, err = setup.Handler(setup.Ctx, setup.WrapRequest(privKey, data, make(map[string]string)))

	require.Error(t, err)
	require.Equal(t, "DID DOC already exists for Schema did:cheqd:test:schema-1: did doc exists", err.Error())
}
