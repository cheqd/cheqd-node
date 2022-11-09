package types_test

import (
	. "github.com/cheqd/cheqd-node/x/did/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type VerificationMaterialTestCase struct {
	vm       VerificationMaterial
	isValid  bool
	errorMsg string
}

var _ = DescribeTable("Verification Method material validation tests", func(testCase VerificationMaterialTestCase) {
	err := testCase.vm.Validate()

	if testCase.isValid {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(testCase.errorMsg))
	}
},

	Entry(
		"Valid Ed25519VerificationKey2020 verification material",
		VerificationMaterialTestCase{
			vm: Ed25519VerificationKey2020{
				PublicKeyMultibase: ValidEd25519PubKey,
			},
			isValid:  true,
			errorMsg: "",
		}),

	Entry(
		"Valid JsonWebKey2020 verification material",
		VerificationMaterialTestCase{
			vm: JsonWebKey2020{
				PublicKeyJwk: ValidPublicKeyJWK,
			},
			isValid:  true,
			errorMsg: "",
		}),

	Entry(
		"Invalid Ed25519VerificationKey2020 verification material",
		VerificationMaterialTestCase{
			vm: Ed25519VerificationKey2020{
				PublicKeyMultibase: InvalidEd25519PubKey,
			},
			isValid:  false,
			errorMsg: "publicKeyMultibase: ed25519: bad public key length: 18",
		}),

	Entry(
		"Invalid JsonWebKey2020 verification material",
		VerificationMaterialTestCase{
			vm: JsonWebKey2020{
				PublicKeyJwk: InvalidPublicKeyJWK,
			},
			isValid:  false,
			errorMsg: "can't parse jwk: failed to parse key",
		}),
)
