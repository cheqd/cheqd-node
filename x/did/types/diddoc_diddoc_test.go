package types_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/did/types"
)

type DIDDocTestCase struct {
	didDoc            *DidDoc
	allowedNamespaces []string
	isValid           bool
	errorMsg          string
}

var _ = DescribeTable("DIDDoc Validation tests", func(testCase DIDDocTestCase) {
	err := testCase.didDoc.Validate(testCase.allowedNamespaces)

	if testCase.isValid {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(testCase.errorMsg))
	}
},

	Entry(
		"DIDDoc is valid",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id: ValidTestDID,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "Ed25519VerificationKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidEd25519VerificationMaterial,
					},
				},
			},
			isValid:  true,
			errorMsg: "",
		}),

	Entry(
		"DIDDoc is invalid",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id: InvalidTestDID,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "Ed25519VerificationKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidEd25519VerificationMaterial,
					},
				},
			},
			isValid:  false,
			errorMsg: "id: unable to split did into method, namespace and id; verification_method: (0: (id: must have prefix: badDid.).).",
		}),

	Entry(
		"Verification method is Ed25519VerificationKey2020",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id: ValidTestDID,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "Ed25519VerificationKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidEd25519VerificationMaterial,
					},
				},
			},
			isValid:  true,
			errorMsg: "",
		}),

	Entry(
		"Verification method is JWK",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id: ValidTestDID,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "JsonWebKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidJWKKeyVerificationMaterial,
					},
				},
			},
			isValid:  true,
			errorMsg: "",
		}),

	Entry("Verification method has wrong ID",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id: ValidTestDID,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   InvalidTestDID,
						Type:                 "JsonWebKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidJWKKeyVerificationMaterial,
					},
				},
			},
			isValid:  false,
			errorMsg: "verification_method: (0: (id: unable to split did into method, namespace and id.).).",
		}),
	Entry(
		"Verification method has wrong controller",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id: ValidTestDID,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "JsonWebKey2020",
						Controller:           InvalidTestDID,
						VerificationMaterial: ValidJWKKeyVerificationMaterial,
					},
				},
			},
			isValid:  false,
			errorMsg: "verification_method: (0: (controller: unable to split did into method, namespace and id.).).",
		}),
	Entry(
		"List of DIDs in controller is allowed",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id:         ValidTestDID,
				Controller: []string{ValidTestDID, ValidTestDID2},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "Ed25519VerificationKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidEd25519VerificationMaterial,
					},
				},
			},
			isValid:  true,
			errorMsg: "",
		}),
	Entry(
		"List of DIDs in controller is not allowed",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Context:    nil,
				Id:         ValidTestDID,
				Controller: []string{ValidTestDID, InvalidTestDID},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "Ed25519VerificationKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidEd25519VerificationMaterial,
					},
				},
			},
			isValid:  false,
			errorMsg: "controller: (1: unable to split did into method, namespace and id.).",
		}),
	Entry(
		"Namespace in controller is not in list of allowed",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id:         ValidTestDID,
				Controller: []string{ValidTestDID},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "Ed25519VerificationKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidEd25519VerificationMaterial,
					},
				},
			},
			allowedNamespaces: []string{"mainnet"},
			isValid:           false,
			errorMsg:          "controller: (0: did namespace must be one of: mainnet.); id: did namespace must be one of: mainnet; verification_method: (0: (controller: did namespace must be one of: mainnet; id: did namespace must be one of: mainnet.).).",
		}),
	Entry(
		"Controller is duplicated",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id:         ValidTestDID,
				Controller: []string{ValidTestDID, ValidTestDID},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "Ed25519VerificationKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidEd25519VerificationMaterial,
					},
				},
			},
			isValid:  false,
			errorMsg: "controller: there should be no duplicates.",
		}),
	Entry(
		"Verification method is duplicated",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id: ValidTestDID,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "Ed25519VerificationKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidEd25519VerificationMaterial,
					},
					{
						Id:                   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:                 "Ed25519VerificationKey2020",
						Controller:           ValidTestDID,
						VerificationMaterial: ValidEd25519VerificationMaterial,
					},
				},
			},
			isValid:  false,
			errorMsg: "verification_method: there are verification method duplicates.",
		}),
)
