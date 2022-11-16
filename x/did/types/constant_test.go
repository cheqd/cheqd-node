package types_test

import (
	"encoding/json"
)

var (
	ValidTestDID         = "did:cheqd:testnet:zABCDEFG123456789abcd"
	ValidTestDID2        = "did:cheqd:testnet:zABCDEFG987654321abcd"
	InvalidTestDID       = "badDid"
	ValidEd25519PubKey   = "zF1hVGXXK9rmx5HhMTpGnGQJiab9qrFJbQXBRhSmYjQWX"
	InvalidEd25519PubKey = "zF1hVGXXK9rmx5HhMTpGnGQJi"
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

var InvalidJWKKey = TestJWKKey{
	Kty: "SomeOtherKeyType",
	N:   "o76AudS2rsCvlz_3D47sFkpuz3NJxgLbXr1cHdmbo9xOMttPMJI97f0rHiSl9stltMi87KIOEEVQWUgMLaWQNaIZThgI1seWDAGRw59AO5sctgM1wPVZYt40fj2Qw4KT7m4RLMsZV1M5NYyXSd1lAAywM4FT25N0RLhkm3u8Hehw2Szj_2lm-rmcbDXzvjeXkodOUszFiOqzqBIS0Bv3c2zj2sytnozaG7aXa14OiUMSwJb4gmBC7I0BjPv5T85CH88VOcFDV51sO9zPJaBQnNBRUWNLh1vQUbkmspIANTzj2sN62cTSoxRhSdnjZQ9E_jraKYEW5oizE9Dtow4EvQ",
	Use: "sig",
	Alg: "RS256",
	E:   "AQAB",
	Kid: "6a8ba5652a7044121d4fedac8f14d14c54e4895b",
}

var (
	ValidPublicKeyJWK, _   = json.Marshal(ValidJWKKey)
	InvalidPublicKeyJWK, _ = json.Marshal(InvalidJWKKey)
)

var (
	ValidEd25519VerificationMaterial   = "{\"publicKeyMultibase\":\"" + ValidEd25519PubKey + "\"}"
	InvalidEd25519VerificationMaterial = "{\"publicKeyMultibase\":\"" + InvalidEd25519PubKey + "\"}"

	ValidJWKKeyVerificationMaterial   = "{\"publicKeyJwk\":" + string(ValidPublicKeyJWK) + "}"
	InvalidJWKKeyVerificationMaterial = "{\"publicKeyJwk\":" + string(InvalidPublicKeyJWK) + "}"
)
