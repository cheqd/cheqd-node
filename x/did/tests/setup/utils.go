package setup

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	mathrand "math/rand"
	"time"

	"cosmossdk.io/math/unsafe"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multibase"
)

func randBase58Seq(bytes int) string {
	b := make([]byte, bytes)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base58.Encode(b)
}

type IDType int

const (
	Base58_16bytes   IDType = iota
	Base58_16symbols IDType = iota
	UUID             IDType = iota
)

var letters = []rune("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}

func ParseJSONToMap(jsonStr string) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GenerateDID(idtype IDType) string {
	prefix := "did:cheqd:" + DidNamespace + ":"
	unsafe.Seed(time.Now().UnixNano())

	switch idtype {
	case Base58_16bytes:
		return prefix + randBase58Seq(16)
	case Base58_16symbols:
		return prefix + randSeq(16)
	case UUID:
		return prefix + uuid.NewString()
	default:
		panic("Unknown ID type")
	}
}

func GenerateKeyPair() KeyPair {
	PublicKey, PrivateKey, _ := ed25519.GenerateKey(rand.Reader)
	return KeyPair{PrivateKey, PublicKey}
}

func GenerateEd25519VerificationKey2020VerificationMaterial(publicKey ed25519.PublicKey) string {
	publicKeyMultibaseBytes := []byte{0xed, 0x01}
	publicKeyMultibaseBytes = append(publicKeyMultibaseBytes, publicKey...)
	keyStr, _ := multibase.Encode(multibase.Base58BTC, publicKeyMultibaseBytes)
	return keyStr
}

func GenerateJSONWebKey2020VerificationMaterial(publicKey ed25519.PublicKey) string {
	pubKeyJwk, err := jwk.New(publicKey)
	if err != nil {
		panic(err)
	}

	pubKeyJwkJSON, err := json.Marshal(pubKeyJwk)
	if err != nil {
		panic(err)
	}

	return string(pubKeyJwkJSON)
}

func GenerateEd25519VerificationKey2018VerificationMaterial(publicKey ed25519.PublicKey) string {
	return base58.Encode(publicKey)
}
