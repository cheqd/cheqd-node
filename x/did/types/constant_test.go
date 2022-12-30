package types_test

import (
	"encoding/json"
)

var (
	ValidTestDID   = "did:cheqd:testnet:zABCDEFG123456789abcd"
	ValidTestDID2  = "did:cheqd:testnet:zABCDEFG987654321abcd"
	InvalidTestDID = "badDid"
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
	// bytes: 										  ed01c92d1e8f9cfa03f63be3489accb0c2704bb7da3f2e4e94509d8ff9202d564c12
	ValidEd25519VerificationMaterial                = "z6MkszZtxCmA2Ce4vUV132PCuLQmwnaDD5mw2L23fGNnsiX3"
	// bytes: 										  020076a50fe5e0c3616c1b4d85a308c104a1c99d8d3d92c18c1f4e0179202d564c12
	InvalidEd25519VerificationMaterialBadMulticodec = "z3dEYJrMxWigf9boyeJMTRN4Ern8DJMoCXaLK77pzQmxVjf"
	// bytes: 										  ed01c92d1e8f9cfa03f63be3489accb0c2704bb7da3f2e4e94509d8ff9
	InvalidEd25519VerificationMaterialBadlength     = "zBm3emgJHyjidq7HsZFTx3PCjYHayy7SxisBeVCa4"

	ValidJWKKeyVerificationMaterial   = string(ValidPublicKeyJWK)
	InvalidJWKKeyVerificationMaterial = string(InvalidPublicKeyJWK)
)
