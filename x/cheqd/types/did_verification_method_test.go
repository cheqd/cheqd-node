package types

import (
	"encoding/json"
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

var	ValidJWKKey = TestJWKKey{
	Kty: "RSA",
	N:   "o76AudS2rsCvlz_3D47sFkpuz3NJxgLbXr1cHdmbo9xOMttPMJI97f0rHiSl9stltMi87KIOEEVQWUgMLaWQNaIZThgI1seWDAGRw59AO5sctgM1wPVZYt40fj2Qw4KT7m4RLMsZV1M5NYyXSd1lAAywM4FT25N0RLhkm3u8Hehw2Szj_2lm-rmcbDXzvjeXkodOUszFiOqzqBIS0Bv3c2zj2sytnozaG7aXa14OiUMSwJb4gmBC7I0BjPv5T85CH88VOcFDV51sO9zPJaBQnNBRUWNLh1vQUbkmspIANTzj2sN62cTSoxRhSdnjZQ9E_jraKYEW5oizE9Dtow4EvQ",
	Use: "sig",
	Alg: "RS256",
	E:   "AQAB",
	Kid: "6a8ba5652a7044121d4fedac8f14d14c54e4895b",
}

var	NotValidJWKKey = TestJWKKey{
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
				Id:         "did:cheqd:aaaaaaaaaaaaaaaa#rty",
				Type:       "JsonWebKey2020",
				Controller: "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "base did: positive",
			struct_: VerificationMethod{
				Id:         "did:cheqd:aaaaaaaaaaaaaaaa#rty",
				Type:       "JsonWebKey2020",
				Controller: "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			baseDid:  "did:cheqd:aaaaaaaaaaaaaaaa",
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "base did: negative",
			struct_: VerificationMethod{
				Id:         "did:cheqd:aaaaaaaaaaaaaaaa#rty",
				Type:       "JsonWebKey2020",
				Controller: "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			baseDid:  "did:cheqd:bbbbbbbbbbbbbbbb",
			isValid:  false,
			errorMsg: "id: must have prefix: did:cheqd:bbbbbbbbbbbbbbbb.",
		},
		{
			name: "allowed namespaces: positive",
			struct_: VerificationMethod{
				Id:         "did:cheqd:mainnet:aaaaaaaaaaaaaaaa#rty",
				Type:       "JsonWebKey2020",
				Controller: "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			allowedNamespaces: []string{"mainnet", ""},
			isValid:           true,
		},
		{
			name: "allowed namespaces: positive",
			struct_: VerificationMethod{
				Id:         "did:cheqd:mainnet:aaaaaaaaaaaaaaaa#rty",
				Type:       "JsonWebKey2020",
				Controller: "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: ValidPublicKeyJWK,
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
				PublicKeyJwk: ValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			isValid:           true,
			errorMsg:          "controller: did namespace must be one of: testnet; id: did namespace must be one of: testnet.",
		},
		{
			name: "JWK: not valid key",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa#qwe",
				Type:               "JsonWebKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: NotValidPublicKeyJWK,
				PublicKeyMultibase: "",
			},
			isValid:           false,
			errorMsg:          "public_key_jwk: invalid format for JWK key, error from validation: failed to unmarshal JWK set: failed to parse sole key in key set: invalid key type from JSON (SomeOtherKeyType).",
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
