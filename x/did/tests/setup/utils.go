package setup

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	mathrand "math/rand"
	"time"

	"github.com/cheqd/cheqd-node/x/did/types"
	. "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/mr-tron/base58"
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

func GenerateDID(idtype IDType) string {
	prefix := "did:cheqd:" + DID_NAMESPACE + ":"
	mathrand.Seed(time.Now().UnixNano())

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

func BuildEd25519VerificationKey2020VerificationMaterial(publicKey ed25519.PublicKey) string {
	return MustEncodeJson(types.Ed25519VerificationKey2020{
		PublicKeyMultibase: MustEncodeMultibaseBase58(publicKey),
	})
}

func BuildJsonWebKey2020VerificationMaterial(publicKey ed25519.PublicKey) string {
	pubKeyJwk, err := jwk.New(publicKey)
	if err != nil {
		panic(err)
	}

	pubKeyJwkJson, err := json.Marshal(pubKeyJwk)
	if err != nil {
		panic(err)
	}

	return MustEncodeJson(types.JsonWebKey2020{
		PublicKeyJwk: json.RawMessage(pubKeyJwkJson),
	})
}
