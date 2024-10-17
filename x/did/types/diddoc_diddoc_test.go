package types_test

import (
	"encoding/json"
	"fmt"
	"strconv"

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
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
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
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
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
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
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
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "JsonWebKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidJWK2020VerificationMaterial,
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
						Id:                     InvalidTestDID,
						VerificationMethodType: "JsonWebKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidJWK2020VerificationMaterial,
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
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "JsonWebKey2020",
						Controller:             InvalidTestDID,
						VerificationMaterial:   ValidJWK2020VerificationMaterial,
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
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
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
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
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
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
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
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
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
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
					},
					{
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
					},
				},
			},
			isValid:  false,
			errorMsg: "verification_method: there are verification method duplicates.",
		}),
	Entry(
		"Assertion method is valid",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id:         ValidTestDID,
				Controller: []string{ValidTestDID},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
					},
				},
				AssertionMethod: []string{fmt.Sprintf("%s#fragment", ValidTestDID), func() string {
					b, _ := json.Marshal(AssertionMethodJSONUnescaped{
						Id:              fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:            "Ed25519VerificationKey2018",
						Controller:      ValidTestDID,
						PublicKeyBase58: &ValidEd25519VerificationKey2018VerificationMaterial, // arbitrarily chosen, loosely validated
					})
					return strconv.Quote(string(b))
				}()},
			},
			isValid:  true,
			errorMsg: "",
		}),
	Entry(
		"Assertion method has wrong fragment",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id:         ValidTestDID,
				Controller: []string{ValidTestDID},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
					},
				},
				AssertionMethod: []string{fmt.Sprintf("%s#fragment", ValidTestDID), func() string {
					b, _ := json.Marshal(AssertionMethodJSONUnescaped{
						Id:              fmt.Sprintf("%s#fragment-1", ValidTestDID),
						Type:            "Ed25519VerificationKey2018",
						Controller:      ValidTestDID,
						PublicKeyBase58: &ValidEd25519VerificationKey2018VerificationMaterial, // arbitrarily chosen, loosely validated
					})
					return strconv.Quote(string(b))
				}()},
			},
			isValid:  false,
			errorMsg: "assertionMethod should be a valid key reference within the DID document's verification method",
		}),
	Entry(
		"Assertion method has invalid protobuf value",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id:         ValidTestDID,
				Controller: []string{ValidTestDID},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
					},
				},
				AssertionMethod: []string{func() string {
					b, _ := json.Marshal(struct {
						Id           string                 `json:"id"`
						Type         string                 `json:"type"`
						Controller   string                 `json:"controller"`
						InvalidField map[string]interface{} `json:"invalidField"`
					}{
						Id:           fmt.Sprintf("%s#fragment", ValidTestDID),
						Type:         "Ed25519VerificationKey2018",
						Controller:   ValidTestDID,
						InvalidField: map[string]interface{}{"unsupported": []int{1, 2, 3}},
					})
					return strconv.Quote(string(b))
				}()},
			},
			isValid:  false,
			errorMsg: "field invalidField is not protobuf-supported",
		}),
	Entry(
		"Assertion method is missing controller value in JSON",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id:         ValidTestDID,
				Controller: []string{ValidTestDID},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
					},
				},
				AssertionMethod: []string{func() string {
					b, _ := json.Marshal(struct {
						Id   string `json:"id"`
						Type string `json:"type"`
					}{
						Id:   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type: "Ed25519VerificationKey2018",
					})
					return strconv.Quote(string(b))
				}()},
			},
			isValid:  false,
			errorMsg: "assertion_method: (0: (controller: cannot be blank.).).",
		}),
	Entry(
		"Assertion method contains unescaped JSON string",
		DIDDocTestCase{
			didDoc: &DidDoc{
				Id:         ValidTestDID,
				Controller: []string{ValidTestDID},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                     fmt.Sprintf("%s#fragment", ValidTestDID),
						VerificationMethodType: "Ed25519VerificationKey2020",
						Controller:             ValidTestDID,
						VerificationMaterial:   ValidEd25519VerificationKey2020VerificationMaterial,
					},
				},
				AssertionMethod: []string{func() string {
					b, _ := json.Marshal(struct {
						Id   string `json:"id"`
						Type string `json:"type"` // controller is intentionally missing, no additional fields are necessary as the focal point is the unescaped JSON string, i.e. deserialisation should fail first, before any other validation
					}{
						Id:   fmt.Sprintf("%s#fragment", ValidTestDID),
						Type: "Ed25519VerificationKey2018",
					})
					return string(b)
				}()},
			},
			isValid:  false,
			errorMsg: "assertionMethod should be a DIDUrl or an Escaped JSON string",
		}),
)
