package tests

import (
	"crypto/ed25519"
	"crypto/rand"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDIDDocShortUniqueId(t *testing.T) {
	setup := Setup()

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, "did:cheqd:test:alice")

	aliceKeys := map[string]ed25519.PrivateKey{"did:cheqd:test:alice#key-1": privKey}
	_, err := setup.SendCreateDid(aliceDid, aliceKeys)

	// Checks
	require.Error(t, err)
	require.Equal(t, "Id: is not DID", err.Error())
}


func TestDIDDocVeryLongUniqueId(t *testing.T) {
	setup := Setup()

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, "did:cheqd:test:abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz")

	aliceKeys := map[string]ed25519.PrivateKey{"did:cheqd:test:alice#key-1": privKey}
	_, err := setup.SendCreateDid(aliceDid, aliceKeys)

	// Checks
	require.Error(t, err)
	require.Equal(t, "Id: is not DID", err.Error())
}


func TestDIDDocNotBase58UniqueId(t *testing.T) {
	setup := Setup()

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, "did:cheqd:test:abcdefghijklmnopqrstuvwxyz?")

	aliceKeys := map[string]ed25519.PrivateKey{"did:cheqd:test:alice#key-1": privKey}
	_, err := setup.SendCreateDid(aliceDid, aliceKeys)

	// Checks
	require.Error(t, err)
	require.Equal(t, "Id: is not DID", err.Error())
}


func TestDIDDocPositiveCaseUniqueId16Symbols(t *testing.T) {
	setup := Setup()

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, "did:cheqd:test:1A1zP1eP5QGefi2D")

	aliceKeys := map[string]ed25519.PrivateKey{"did:cheqd:test:1A1zP1eP5QGefi2D#key-1": privKey}
	did, _ := setup.SendCreateDid(aliceDid, aliceKeys)

	// Checks
	require.Equal(t, did.Id, "did:cheqd:test:1A1zP1eP5QGefi2D")
}

func TestDIDDocPositiveCaseUniqueId32Symbols(t *testing.T) {
	setup := Setup()

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, "did:cheqd:test:1A1zP1eP5QGefi2DMPTfTL5SLmv7Divf")

	aliceKeys := map[string]ed25519.PrivateKey{"did:cheqd:test:1A1zP1eP5QGefi2DMPTfTL5SLmv7Divf#key-1": privKey}
	did, _ := setup.SendCreateDid(aliceDid, aliceKeys)

	// Checks
	require.Equal(t, did.Id, "did:cheqd:test:1A1zP1eP5QGefi2DMPTfTL5SLmv7Divf")
}