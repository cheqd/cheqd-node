package tests

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestDIDDocControllerChanged(t *testing.T) {
	setup := Setup()

	//Init did
	aliceKeys, aliceDid, _ := setup.InitDid(AliceDID)
	bobKeys, _, _ := setup.InitDid(BobDID)

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.Controller = append(updatedDidDoc.Controller, BobDID)
	receivedDid, _ := setup.SendUpdateDid(updatedDidDoc, ConcatKeys(aliceKeys, bobKeys))

	// check
	require.NotEqual(t, aliceDid.Controller, receivedDid.Controller)
	require.NotEqual(t, []string{AliceDID, BobDID}, receivedDid.Controller)
	require.Equal(t, []string{BobDID}, receivedDid.Controller)
}

func TestDIDDocVerificationMethodChangedWithoutOldSignature(t *testing.T) {
	setup := Setup()

	//Init did
	_, aliceDid, _ := setup.InitDid(AliceDID)
	bobKeys, _, _ := setup.InitDid(BobDID)

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.VerificationMethod[0].Type = "Ed25519VerificationKey2020"
	_, err := setup.SendUpdateDid(updatedDidDoc, bobKeys)

	// check
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf("signature %v not found: invalid signature detected", AliceDID), err.Error())
}

func TestDIDDocVerificationMethodControllerChangedWithoutOldSignature(t *testing.T) {
	setup := Setup()

	//Init did
	_, aliceDid, _ := setup.InitDid(AliceDID)
	bobKeys, _, _ := setup.InitDid(BobDID)

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.VerificationMethod[0].Controller = BobDID
	_, err := setup.SendUpdateDid(updatedDidDoc, bobKeys)

	// check
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf("signature %v not found: invalid signature detected", AliceDID), err.Error())
}

func TestDIDDocControllerChangedWithoutOldSignature(t *testing.T) {
	setup := Setup()

	//Init did
	_, aliceDid, _ := setup.InitDid(AliceDID)
	bobKeys, _, _ := setup.InitDid(BobDID)

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.Controller = append(updatedDidDoc.Controller, BobDID)
	_, err := setup.SendUpdateDid(updatedDidDoc, bobKeys)

	// check
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf("signature %v not found: invalid signature detected", AliceDID), err.Error())
}

func TestDIDDocVerificationMethodDeletedWithoutOldSignature(t *testing.T) {
	setup := Setup()

	//Init did
	_, bodDidDoc, _ := setup.InitDid(BobDID)

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, AliceDID)

	aliceDid.VerificationMethod = append(aliceDid.VerificationMethod, &types.VerificationMethod{
		Id:                 AliceKey2,
		Controller:         BobDID,
		Type:               "Ed25519VerificationKey2020",
		PublicKeyMultibase: bodDidDoc.VerificationMethod[0].PublicKeyMultibase,
	})

	aliceDid.Authentication = append(aliceDid.Authentication, AliceKey2)

	aliceKeys := map[string]ed25519.PrivateKey{AliceKey1: privKey}
	_, _ = setup.SendCreateDid(aliceDid, aliceKeys)

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.VerificationMethod = []*types.VerificationMethod{aliceDid.VerificationMethod[0]}
	updatedDidDoc.Authentication = []string{aliceDid.Authentication[0]}
	_, err := setup.SendUpdateDid(updatedDidDoc, aliceKeys)

	// check
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf("signature %v not found: invalid signature detected", BobDID), err.Error())
}

func TestDIDDocVerificationMethodDeleted(t *testing.T) {
	setup := Setup()

	//Init did
	bobKeys, bodDidDoc, _ := setup.InitDid(BobDID)

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, AliceDID)

	aliceDid.Authentication = append(aliceDid.Authentication, AliceKey2)
	aliceDid.VerificationMethod = append(aliceDid.VerificationMethod, &types.VerificationMethod{
		Id:                 AliceKey2,
		Controller:         BobDID,
		Type:               "Ed25519VerificationKey2020",
		PublicKeyMultibase: bodDidDoc.VerificationMethod[0].PublicKeyMultibase,
	})

	aliceKeys := map[string]ed25519.PrivateKey{AliceKey1: privKey}
	_, _ = setup.SendCreateDid(aliceDid, aliceKeys)

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.Authentication = []string{aliceDid.Authentication[0]}
	updatedDidDoc.VerificationMethod = []*types.VerificationMethod{aliceDid.VerificationMethod[0]}
	receivedDid, _ := setup.SendUpdateDid(updatedDidDoc, ConcatKeys(aliceKeys, bobKeys))

	// check
	require.NotEqual(t, len(aliceDid.VerificationMethod), len(receivedDid.VerificationMethod))
	require.True(t, reflect.DeepEqual(aliceDid.VerificationMethod[0], receivedDid.VerificationMethod[0]))
}
