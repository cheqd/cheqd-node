package types_test

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"

	. "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/multiformats/go-multibase"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type VerificationMethodTestCase struct {
	vm                VerificationMethod
	baseDid           string
	allowedNamespaces []string
	isValid           bool
	errorMsg          string
}

var _ = DescribeTable("Verification Method Validation tests", func(testCase VerificationMethodTestCase) {
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
			vm: VerificationMethod{
				Id:                   "did:cheqd:zABCDEFG123456789abcd#qwe",
				Type:                 "Ed25519VerificationKey2020",
				Controller:           "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial: ValidEd25519VerificationMaterial,
			},
			isValid:  true,
			errorMsg: "",
		}),

	Entry(
		"Verification method with expected jwk key",
		VerificationMethodTestCase{
			vm: VerificationMethod{
				Id:                   "did:cheqd:zABCDEFG123456789abcd#rty",
				Type:                 "JsonWebKey2020",
				Controller:           "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial: ValidJWKKeyVerificationMaterial,
			},
			isValid:  true,
			errorMsg: "",
		}),

	Entry(
		"Id has expected DID as a base",
		VerificationMethodTestCase{
			vm: VerificationMethod{
				Id:                   "did:cheqd:zABCDEFG123456789abcd#rty",
				Type:                 "JsonWebKey2020",
				Controller:           "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial: ValidJWKKeyVerificationMaterial,
			},
			baseDid:  "did:cheqd:zABCDEFG123456789abcd",
			isValid:  true,
			errorMsg: "",
		}),

	Entry(
		"Id does not have expected DID as a base",
		VerificationMethodTestCase{
			vm: VerificationMethod{
				Id:                   "did:cheqd:zABCDEFG123456789abcd#rty",
				Type:                 "JsonWebKey2020",
				Controller:           "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial: ValidJWKKeyVerificationMaterial,
			},
			baseDid:  "did:cheqd:zABCDEFG987654321abcd",
			isValid:  false,
			errorMsg: "id: must have prefix: did:cheqd:zABCDEFG987654321abcd.",
		}),

	Entry(
		"Namespace is allowed",
		VerificationMethodTestCase{
			vm: VerificationMethod{
				Id:                   "did:cheqd:mainnet:zABCDEFG123456789abcd#rty",
				Type:                 "JsonWebKey2020",
				Controller:           "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial: ValidJWKKeyVerificationMaterial,
			},
			allowedNamespaces: []string{"mainnet", ""},
			isValid:           true,
		}),

	Entry(
		"Namespace is not allowed",
		VerificationMethodTestCase{
			vm: VerificationMethod{
				Id:                   "did:cheqd:mainnet:zABCDEFG123456789abcd#rty",
				Type:                 "JsonWebKey2020",
				Controller:           "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial: ValidJWKKeyVerificationMaterial,
			},
			allowedNamespaces: []string{"testnet"},
			isValid:           false,
			errorMsg:          "controller: did namespace must be one of: testnet; id: did namespace must be one of: testnet.",
		}),
	Entry(
		"JWK key has expected format",
		VerificationMethodTestCase{
			vm: VerificationMethod{
				Id:                   "did:cheqd:zABCDEFG123456789abcd#qwe",
				Type:                 "JsonWebKey2020",
				Controller:           "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial: ValidJWKKeyVerificationMaterial,
			},
			isValid: true,
		}),
	Entry(
		"JWK key has unexpected format",
		VerificationMethodTestCase{
			vm: VerificationMethod{
				Id:                   "did:cheqd:zABCDEFG123456789abcd#qwe",
				Type:                 "JsonWebKey2020",
				Controller:           "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial: InvalidJWKKeyVerificationMaterial,
			},
			isValid:  false,
			errorMsg: "verification_material: can't parse jwk: failed to parse key",
		}),
	Entry(
		"Ed25519 key has unexpected format",
		VerificationMethodTestCase{
			vm: VerificationMethod{
				Id:                   "did:cheqd:zABCDEFG123456789abcd#qwe",
				Type:                 "Ed25519VerificationKey2020",
				Controller:           "did:cheqd:zABCDEFG987654321abcd",
				VerificationMaterial: InvalidEd25519VerificationMaterial,
			},
			isValid:  false,
			errorMsg: "verification_material: (publicKeyMultibase: ed25519: bad public key length: 18.)",
		}),
)

var _ = Describe("Validation ed25519 Signature in verification method", func() {
	var pubKey ed25519.PublicKey
	var privKey ed25519.PrivateKey
	var err error
	message := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod " +
		"tempor incididunt ut labore et dolore magna aliqua."
	msgBytes := []byte(message)
	var signature []byte

	pubKey, privKey, err = ed25519.GenerateKey(rand.Reader)
	Expect(err).To(BeNil())

	signature = ed25519.Sign(privKey, msgBytes)

	Context("when ed25519 key is placed", func() {
		It("is valid", func() {
			pubKeyStr, err := multibase.Encode(multibase.Base58BTC, pubKey)
			Expect(err).To(BeNil())

			vm := VerificationMethod{
				Id:                   "",
				Type:                 "Ed25519VerificationKey2020",
				Controller:           "",
				VerificationMaterial: "{\"publicKeyMultibase\": \"" + pubKeyStr + "\"}",
			}

			err = VerifySignature(vm, msgBytes, signature)
			Expect(err).To(BeNil())
		})
	})

	Context("when with the same env but JWK is placed", func() {
		It("is valid", func() {
			jwk_, err := jwk.New(pubKey)
			Expect(err).To(BeNil())

			json_, err := json.MarshalIndent(jwk_, "", "  ")
			Expect(err).To(BeNil())

			vm2 := VerificationMethod{
				Id:                   "",
				Type:                 "JsonWebKey2020",
				Controller:           "",
				VerificationMaterial: "{\"publicKeyJwk\": " + string(json_) + "}",
			}

			err = VerifySignature(vm2, msgBytes, signature)
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

			jwk_, err := jwk.New(pubKey)
			Expect(err).To(BeNil())

			json_, err := json.MarshalIndent(jwk_, "", "  ")
			Expect(err).To(BeNil())

			vm := VerificationMethod{
				Id:                   "",
				Type:                 "JsonWebKey2020",
				Controller:           "",
				VerificationMaterial: "{\"publicKeyJwk\": " + string(json_) + "}",
			}

			err = VerifySignature(vm, msgBytes, signature)
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

			jwk_, err := jwk.New(pubKey)
			Expect(err).To(BeNil())

			json_, err := json.Marshal(jwk_)
			Expect(err).To(BeNil())

			vm2 := VerificationMethod{
				Id:                   "",
				Type:                 "JsonWebKey2020",
				Controller:           "",
				VerificationMaterial: "{\"publicKeyJwk\": " + string(json_) + "}",
			}

			err = VerifySignature(vm2, msgBytes, signature)
			Expect(err).To(BeNil())
		})
	})
})
