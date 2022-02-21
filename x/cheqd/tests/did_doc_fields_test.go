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


func TestDIDDocPositiveCaseUniqueId(t *testing.T) {
	setup := Setup()

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, "did:cheqd:test:abcdefghijkl1234")

	aliceKeys := map[string]ed25519.PrivateKey{"did:cheqd:test:abcdefghijkl1234#key-1": privKey}
	did, _ := setup.SendCreateDid(aliceDid, aliceKeys)

	// Checks
	require.Equal(t, did.Id, "did:cheqd:test:abcdefghijkl1234")
}

func TestDIDDocPositiveCaseUniqueId32Symbols(t *testing.T) {
	setup := Setup()

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, "did:cheqd:test:abcdefghijkl1234abcdefghijkl1234")

	aliceKeys := map[string]ed25519.PrivateKey{"did:cheqd:test:abcdefghijkl1234abcdefghijkl1234#key-1": privKey}
	did, _ := setup.SendCreateDid(aliceDid, aliceKeys)

	// Checks
	require.Equal(t, did.Id, "did:cheqd:test:abcdefghijkl1234abcdefghijkl1234")
}