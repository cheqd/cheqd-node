package utils_test

import (
	"encoding/json"

	. "github.com/cheqd/cheqd-node/x/cheqd/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

var (
	ValidJWKByte, _    = json.Marshal(ValidJWKKey)
	NotValidJWKByte, _ = json.Marshal(NotValidJWKKey)
)

var _ = Describe("Encoding checks", func() {
	DescribeTable("Is valid multibase key",

		func(data string, isValid bool) {
			_err := ValidateMultibase(data)
			if isValid {
				Expect(_err).ShouldNot(HaveOccurred())
			} else {
				Expect(_err).Should(HaveOccurred())
			}
		},

		Entry("Valid: General pmbkey", "zABCDEFG123456789", true),
		Entry("Not Valid: cannot be empty", "", false),
		Entry("Not Valid: without z but base58", "ABCDEFG123456789", false),
		Entry("Not Valid: without z and not base58", "OIl0ABCDEFG123456789", false),
		Entry("Not Valid: with z but not base58", "zOIl0ABCDEFG123456789", false),
	)

	DescribeTable("Validate Base58",

		func(data string, isValid bool) {
			_err := ValidateBase58(data)
			if isValid {
				Expect(_err).ShouldNot(HaveOccurred())
			} else {
				Expect(_err).Should(HaveOccurred())
			}
		},

		Entry("Valid: General pmbkey", "ABCDEFG123456789", true),
		Entry("Not Valid: cannot be empty", "", false),
		Entry("Not Valid: not base58", "OIl0ABCDEFG123456789", false),
	)

	DescribeTable("Validate JWK",

		func(data string, isValid bool) {
			_err := ValidateJWK(data)
			if isValid {
				Expect(_err).ShouldNot(HaveOccurred())
			} else {
				Expect(_err).Should(HaveOccurred())
			}
		},

		Entry("Valid: General jwk", string(ValidJWKByte), true),
		Entry("Not Valid: Bad jwk", string(NotValidJWKByte), false),
	)
})
