package types

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/multiformats/go-multibase"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestJWKKey struct {
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
}

var ValidJWKKey = TestJWKKey{
	Kty: "RSA",
	N:   "o76AudS2rsCvlz_3D47sFkpuz3NJxgLbXr1cHdmbo9xOMttPMJI97f0rHiSl9stltMi87KIOEEVQWUgMLaWQNaIZThgI1seWDAGRw59AO5sctgM1wPVZYt40fj2Qw4KT7m4RLMsZV1M5NYyXSd1lAAywM4FT25N0RLhkm3u8Hehw2Szj_2lm-rmcbDXzvjeXkodOUszFiOqzqBIS0Bv3c2zj2sytnozaG7aXa14OiUMSwJb4gmBC7I0BjPv5T85CH88VOcFDV51sO9zPJaBQnNBRUWNLh1vQUbkmspIANTzj2sN62cTSoxRhSdnjZQ9E_jraKYEW5oizE9Dtow4EvQ",
	Use: "sig",
	Alg: "RS256",
	E:   "AQAB",
	Kid: "6a8ba5652a7044121d4fedac8f14d14c54e4895b",
}

var NotValidJWKKey = TestJWKKey{
	Kty: "SomeOtherKeyType",
	N:   "o76AudS2rsCvlz_3D47sFkpuz3NJxgLbXr1cHdmbo9xOMttPMJI97f0rHiSl9stltMi87KIOEEVQWUgMLaWQNaIZThgI1seWDAGRw59AO5sctgM1wPVZYt40fj2Qw4KT7m4RLMsZV1M5NYyXSd1lAAywM4FT25N0RLhkm3u8Hehw2Szj_2lm-rmcbDXzvjeXkodOUszFiOqzqBIS0Bv3c2zj2sytnozaG7aXa14OiUMSwJb4gmBC7I0BjPv5T85CH88VOcFDV51sO9zPJaBQnNBRUWNLh1vQUbkmspIANTzj2sN62cTSoxRhSdnjZQ9E_jraKYEW5oizE9Dtow4EvQ",
	Use: "sig",
	Alg: "RS256",
	E:   "AQAB",
	Kid: "6a8ba5652a7044121d4fedac8f14d14c54e4895b",
}

var ValidJWKByte, _ = json.Marshal(ValidJWKKey)
var NotValidJWKByte, _ = json.Marshal(NotValidJWKKey)

var ValidPublicKeyJWK = JSONToPubKeyJWK(string(ValidJWKByte))
var NotValidPublicKeyJWK = JSONToPubKeyJWK(string(NotValidJWKByte))

func TestVerificationMethodValidation(t *testing.T) {
	cases := []struct {
		name              string
		struct_           VerificationMethod
		baseDid           string
		allowedNamespaces []string
		isValid           bool
		errorMsg          string
	}{
		{
			name: "valid method with multibase key",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa#qwe",
				Type:               "Ed25519VerificationKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       nil,
				PublicKeyMultibase: ValidEd25519PubKey,
			},
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "valid method with jwk key",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa#rty",
				Type:               "JsonWebKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "base did: positive",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa#rty",
				Type:               "JsonWebKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			baseDid:  "did:cheqd:aaaaaaaaaaaaaaaa",
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "base did: negative",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa#rty",
				Type:               "JsonWebKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			baseDid:  "did:cheqd:bbbbbbbbbbbbbbbb",
			isValid:  false,
			errorMsg: "id: must have prefix: did:cheqd:bbbbbbbbbbbbbbbb.",
		},
		{
			name: "allowed namespaces: positive",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:mainnet:aaaaaaaaaaaaaaaa#rty",
				Type:               "JsonWebKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			allowedNamespaces: []string{"mainnet", ""},
			isValid:           true,
		},
		{
			name: "allowed namespaces: positive",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:mainnet:aaaaaaaaaaaaaaaa#rty",
				Type:               "JsonWebKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			allowedNamespaces: []string{"testnet"},
			isValid:           false,
			errorMsg:          "controller: did namespace must be one of: testnet; id: did namespace must be one of: testnet.",
		},
		{
			name: "JWK: valid key",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa#qwe",
				Type:               "JsonWebKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			isValid: true,
		},
		{
			name: "JWK: not valid key",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa#qwe",
				Type:               "JsonWebKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       NotValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			isValid:  false,
			errorMsg: "public_key_jwk: can't parse jwk: failed to parse key: invalid key type from JSON (SomeOtherKeyType).",
		},
		{
			name: "all keys and values are required in jwk",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa#qwe",
				Type:               "JsonWebKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       append(ValidPublicKeyJWK, &KeyValuePair{Key: "", Value: ""}),
				PublicKeyMultibase: "",
			},
			isValid:  false,
			errorMsg: "public_key_jwk: (6: (key: cannot be blank; value: cannot be blank.).).",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.struct_.Validate(tc.baseDid, tc.allowedNamespaces)

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, err.Error(), tc.errorMsg)
			}
		})
	}
}

func TestEd25519SignatureVerification(t *testing.T) {
	message := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod " +
		"tempor incididunt ut labore et dolore magna aliqua."
	msgBytes := []byte(message)

	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	signature := ed25519.Sign(privKey, msgBytes)

	pubKeyStr, err := multibase.Encode(multibase.Base58BTC, pubKey)
	require.NoError(t, err)

	vm := VerificationMethod{
		Id:                 "",
		Type:               "Ed25519VerificationKey2020",
		Controller:         "",
		PublicKeyJwk:       nil,
		PublicKeyMultibase: pubKeyStr,
	}

	err = VerifySignature(vm, msgBytes, signature)
	require.NoError(t, err)

	jwk_, err := jwk.New(pubKey)
	require.NoError(t, err)
	json_, err := json.MarshalIndent(jwk_, "", "  ")
	require.NoError(t, err)
	pubKeyJwk := JSONToPubKeyJWK(string(json_))

	vm2 := VerificationMethod{
		Id:                 "",
		Type:               "JsonWebKey2020",
		Controller:         "",
		PublicKeyJwk:       pubKeyJwk,
		PublicKeyMultibase: "",
	}

	err = VerifySignature(vm2, msgBytes, signature)
	require.NoError(t, err)
}

func TestECDSASignatureVerification(t *testing.T) {
	message := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod " +
		"tempor incididunt ut labore et dolore magna aliqua."
	msgBytes := []byte(message)

	hasher := crypto.SHA256.New()
	hasher.Write(msgBytes)
	msgDigest := hasher.Sum(nil)

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	pubKey := privKey.PublicKey

	signature, err := ecdsa.SignASN1(rand.Reader, privKey, msgDigest)
	require.NoError(t, err)

	jwk_, err := jwk.New(pubKey)
	require.NoError(t, err)
	json_, err := json.MarshalIndent(jwk_, "", "  ")
	require.NoError(t, err)
	pubKeyJwk := JSONToPubKeyJWK(string(json_))

	vm2 := VerificationMethod{
		Id:                 "",
		Type:               "JsonWebKey2020",
		Controller:         "",
		PublicKeyJwk:       pubKeyJwk,
		PublicKeyMultibase: "",
	}

	err = VerifySignature(vm2, msgBytes, signature)
	require.NoError(t, err)
}

func TestRSASignatureVerification(t *testing.T) {
	message := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod " +
		"tempor incididunt ut labore et dolore magna aliqua."
	msgBytes := []byte(message)

	hasher := crypto.SHA256.New()
	hasher.Write(msgBytes)
	msgDigest := hasher.Sum(nil)

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pubKey := privKey.PublicKey

	signature, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, msgDigest, nil)
	require.NoError(t, err)

	jwk_, err := jwk.New(pubKey)
	require.NoError(t, err)
	json_, err := json.MarshalIndent(jwk_, "", "  ")
	require.NoError(t, err)
	pubKeyJwk := JSONToPubKeyJWK(string(json_))

	vm2 := VerificationMethod{
		Id:                 "",
		Type:               "JsonWebKey2020",
		Controller:         "",
		PublicKeyJwk:       pubKeyJwk,
		PublicKeyMultibase: "",
	}

	err = VerifySignature(vm2, msgBytes, signature)
	require.NoError(t, err)
}