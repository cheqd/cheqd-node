package utils

import (
	"encoding/json"
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


func TestValidateMultibase(t *testing.T) {
	cases := []struct {
		name string
		data  string
		encoding multibase.Encoding
		valid bool
	}{
		{"Valid: General pmbkey", "zABCDEFG123456789", multibase.Base58BTC, true},
		{"Not Valid: cannot be empty", "", multibase.Base58BTC, false},
		{"Not Valid: without z but base58", "ABCDEFG123456789", multibase.Base58BTC, false},
		{"Not Valid: without z and not base58", "OIl0ABCDEFG123456789", multibase.Base58BTC, false},
		{"Not Valid: with z but not base58", "zOIl0ABCDEFG123456789", multibase.Base58BTC, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMultibase(tc.data)

			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestValidateBase58(t *testing.T) {
	cases := []struct {
		name string
		data  string
		valid bool
	}{
		{"Valid: General pmbkey", "ABCDEFG123456789", true},
		{"Not Valid: cannot be empty", "", false},
		{"Not Valid: not base58", "OIl0ABCDEFG123456789", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateBase58(tc.data)

			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestValidateJWKEncoding(t *testing.T) {
	cases := []struct {
		name string
		data  string
		valid bool
	}{
		{"Valid: General jwk", string(ValidJWKByte), true},
		{"Not Valid: Bad jwk", string(NotValidJWKByte), false},

	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateJWKEncoding(tc.data)

			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
