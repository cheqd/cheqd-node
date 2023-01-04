package types_test

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"

	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/lestrrat-go/jwx/jwk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type VerificationMethodTestCase struct {
	vm                didtypes.VerificationMethod
	baseDid           string
	allowedNamespaces []string
	isValid           bool
	errorMsg          string
}

var _ = DescribeTable("Verification Method validation tests", func(testCase VerificationMethodTestCase) {
	err := testCase.vm.Validate(testCase.baseDid, testCase.allowedNamespaces)

	if testCase.isValid {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(testCase.errorMsg))
	}
},
	Entry(
		"Verification method with expected multibase key",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "Ed25519VerificationKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
			},
			isValid:  true,
			errorMsg: "",
		},
	),
	Entry(
		"Verification method with expected jwk key",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#rty",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidJWK2020VerificationMaterial,
			},
			isValid:  true,
			errorMsg: "",
		},
	),
	Entry(
		"Verification method with expected base58 key",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#uio",
				VerificationMethodType: "Ed25519VerificationKey2018",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidEd25519VerificationKey2018VerificationMaterial,
			},
			isValid:  true,
			errorMsg: "",
		},
	),
	Entry(
		"Id has expected DID as a base",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#rty",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidJWK2020VerificationMaterial,
			},
			baseDid:  "did:cheqd:zABCDEFG123456789abcd",
			isValid:  true,
			errorMsg: "",
		},
	),
	Entry(
		"Id does not have expected DID as a base",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#rty",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidJWK2020VerificationMaterial,
			},
			baseDid:  "did:cheqd:zABCDEFG987654321abcd",
			isValid:  false,
			errorMsg: "id: must have prefix: did:cheqd:zABCDEFG987654321abcd.",
		},
	),
	Entry(
		"Namespace is allowed",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:mainnet:zABCDEFG123456789abcd#rty",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidJWK2020VerificationMaterial,
			},
			allowedNamespaces: []string{"mainnet", ""},
			isValid:           true,
		},
	),

	Entry(
		"Namespace is not allowed",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:mainnet:zABCDEFG123456789abcd#rty",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidJWK2020VerificationMaterial,
			},
			allowedNamespaces: []string{"testnet"},
			isValid:           false,
			errorMsg:          "controller: did namespace must be one of: testnet; id: did namespace must be one of: testnet.",
		},
	),
	Entry(
		"JWK key has expected format",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidJWK2020VerificationMaterial,
			},
			isValid: true,
		},
	),
	Entry(
		"JWK key has unexpected format",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   InvalidJWK2020VerificationMaterial,
			},
			isValid:  false,
			errorMsg: "verification_material: can't parse jwk: failed to parse key",
		},
	),
	Entry(
		"Ed25519 key 2020 has unexpected format",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "Ed25519VerificationKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   InvalidEd25519VerificationKey2020VerificationMaterialBadlength,
			},
			isValid:  false,
			errorMsg: "verification_material: ed25519: bad public key length: 27",
		},
	),
	Entry(
		"Ed25519 key 2018 has unexpected format",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "Ed25519VerificationKey2018",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   InvalidEd25519VerificationKey2018VerificationMaterialBadLength,
			},
			isValid:  false,
			errorMsg: "verification_material: ed25519: bad public key length: 31",
		},
	),
)

var _ = Describe("Validation ed25519 Signature in verification method", func() {
	var pubKey ed25519.PublicKey
	var privKey ed25519.PrivateKey
	var signature []byte
	var err error
	message := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod " +
		"tempor incididunt ut labore et dolore magna aliqua."
	msgBytes := []byte(message)

	pubKey, privKey, err = ed25519.GenerateKey(rand.Reader)
	Expect(err).To(BeNil())

	signature = ed25519.Sign(privKey, msgBytes)

	Context("when ed25519 key 2020 representation is placed", func() {
		It("is valid", func() {
			pubKeyStr := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(pubKey)

			vm := didtypes.VerificationMethod{
				Id:                     "",
				VerificationMethodType: "Ed25519VerificationKey2020",
				Controller:             "",
				VerificationMaterial:   pubKeyStr,
			}

			err = didtypes.VerifySignature(vm, msgBytes, signature)
			Expect(err).To(BeNil())
		})
	})

	Context("when JWK 2020 representation is placed", func() {
		It("is valid", func() {
			jwkKey, err := jwk.New(pubKey)
			Expect(err).To(BeNil())

			jsonPayload, err := json.MarshalIndent(jwkKey, "", "  ")
			Expect(err).To(BeNil())

			vm2 := didtypes.VerificationMethod{
				Id:                     "",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "",
				VerificationMaterial:   string(jsonPayload),
			}

			err = didtypes.VerifySignature(vm2, msgBytes, signature)
			Expect(err).To(BeNil())
		})
	})

	Context("when ed25519 key 2018 representation is placed", func() {
		It("is valid", func() {
			pubKeyStr := testsetup.GenerateEd25519VerificationKey2018VerificationMaterial(pubKey)

			vm := didtypes.VerificationMethod{
				Id:                     "",
				VerificationMethodType: "Ed25519VerificationKey2018",
				Controller:             "",
				VerificationMaterial:   pubKeyStr,
			}

			err = didtypes.VerifySignature(vm, msgBytes, signature)
			Expect(err).To(BeNil())
		})
	})
})

var _ = Describe("Validation ECDSA Signature in verification method", func() {
	Context("ECDSA signature preparations and verification", func() {
		It("is positive case", func() {
			message := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod " +
				"tempor incididunt ut labore et dolore magna aliqua."

			msgBytes := []byte(message)

			hasher := crypto.SHA256.New()
			hasher.Write(msgBytes)
			msgDigest := hasher.Sum(nil)

			privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			Expect(err).To(BeNil())

			pubKey := privKey.PublicKey

			signature, err := ecdsa.SignASN1(rand.Reader, privKey, msgDigest)
			Expect(err).To(BeNil())

			jwkKey, err := jwk.New(pubKey)
			Expect(err).To(BeNil())

			jsonPayload, err := json.MarshalIndent(jwkKey, "", "  ")
			Expect(err).To(BeNil())

			vm := didtypes.VerificationMethod{
				Id:                     "",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "",
				VerificationMaterial:   string(jsonPayload),
			}

			err = didtypes.VerifySignature(vm, msgBytes, signature)
			Expect(err).To(BeNil())
		})
	})
})

var _ = Describe("Validation RSA Signature in verification method", func() {
	Context("RSA signature preparations and verification", func() {
		It("is positive case", func() {
			message := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod " +
				"tempor incididunt ut labore et dolore magna aliqua."
			msgBytes := []byte(message)

			hasher := crypto.SHA256.New()
			hasher.Write(msgBytes)
			msgDigest := hasher.Sum(nil)

			privKey, err := rsa.GenerateKey(rand.Reader, 2048)
			Expect(err).To(BeNil())

			pubKey := privKey.PublicKey

			signature, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, msgDigest, nil)
			Expect(err).To(BeNil())

			jwkKey, err := jwk.New(pubKey)
			Expect(err).To(BeNil())

			jsonPayload, err := json.Marshal(jwkKey)
			Expect(err).To(BeNil())

			vm2 := didtypes.VerificationMethod{
				Id:                     "",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "",
				VerificationMaterial:   string(jsonPayload),
			}

			err = didtypes.VerifySignature(vm2, msgBytes, signature)
			Expect(err).To(BeNil())
		})
	})
})

var _ = DescribeTable("Verification Method material validation tests", func(testCase VerificationMethodTestCase) {
	err := testCase.vm.Validate(testCase.baseDid, testCase.allowedNamespaces)

	if testCase.isValid {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(testCase.errorMsg))
	}
},
	Entry(
		"Valid Ed25519VerificationKey2020 verification material",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "Ed25519VerificationKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
			},
			isValid:  true,
			errorMsg: "",
		},
	),
	Entry(
		"Valid JsonWebKey2020 verification material",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidJWK2020VerificationMaterial,
			},
			isValid:  true,
			errorMsg: "",
		},
	),
	Entry(
		"Valid Ed25519VerificationKey2018 verification material",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "Ed25519VerificationKey2018",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   ValidEd25519VerificationKey2018VerificationMaterial,
			},
			isValid:  true,
			errorMsg: "",
		},
	),
	Entry(
		"Invalid Ed25519VerificationKey2020 verification material",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "Ed25519VerificationKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   InvalidEd25519VerificationKey2020VerificationMaterialBadlength,
			},
			isValid:  false,
			errorMsg: "",
		},
	),
	Entry(
		"Invalid Ed25519VerificationKey2020 verification material",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "Ed25519VerificationKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   InvalidEd25519VerificationKey2020VerificationMaterialBadPrefix,
			},
			isValid:  false,
			errorMsg: "verification_material: invalid two-byte prefix for Ed25519VerificationKey2020. expected: 0xed01 actual: 0x0200",
		},
	),
	Entry(
		"Invalid JsonWebKey2020 verification material",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "JsonWebKey2020",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   InvalidJWK2020VerificationMaterial,
			},
			isValid:  false,
			errorMsg: "can't parse jwk: failed to parse key",
		},
	),
	Entry(
		"Invalid Ed25519VerificationKey2018 verification material",
		VerificationMethodTestCase{
			vm: didtypes.VerificationMethod{
				Id:                     "did:cheqd:zABCDEFG123456789abcd#qwe",
				VerificationMethodType: "Ed25519VerificationKey2018",
				Controller:             "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial:   InvalidEd25519VerificationKey2018VerificationMaterialBadLength,
			},
			isValid:  false,
			errorMsg: "ed25519: bad public key length: 31",
		},
	),
)
