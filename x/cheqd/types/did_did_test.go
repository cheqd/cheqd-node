package types_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("DID Validation tests", func() {
	type TestCaseDIDStruct struct {
		did               *Did
		allowedNamespaces []string
		isValid           bool
		errorMsg          string
	}

	DescribeTable("DID Validation tests", func(testCase TestCaseDIDStruct) {
		err := testCase.did.Validate(testCase.allowedNamespaces)

		if testCase.isValid {
			Expect(err).To(BeNil())
		} else {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(testCase.errorMsg))
		}
	},

		Entry(
			"Did is valid",
			TestCaseDIDStruct{
				did: &Did{
					Id: ValidTestDID,
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "Ed25519VerificationKey2020",
							Controller:         ValidTestDID,
							PublicKeyJwk:       nil,
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
				},
				isValid:  true,
				errorMsg: "",
			}),

		Entry(
			"DID is not allowed",
			TestCaseDIDStruct{
				did: &Did{
					Id: InvalidTestDID,
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "Ed25519VerificationKey2020",
							Controller:         ValidTestDID,
							PublicKeyJwk:       nil,
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
				},
				isValid:  false,
				errorMsg: "id: unable to split did into method, namespace and id; verification_method: (0: (id: must have prefix: badDid.).).",
			}),

		Entry(
			"Verification method is Ed25519VerificationKey2020",
			TestCaseDIDStruct{
				did: &Did{
					Id: ValidTestDID,
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "Ed25519VerificationKey2020",
							Controller:         ValidTestDID,
							PublicKeyJwk:       nil,
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
				},
				isValid:  true,
				errorMsg: "",
			}),

		Entry(
			"Verification method is jwk",
			TestCaseDIDStruct{
				did: &Did{
					Id: ValidTestDID,
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "JsonWebKey2020",
							Controller:         ValidTestDID,
							PublicKeyJwk:       ValidPublicKeyJWK,
							PublicKeyMultibase: "",
						},
					},
				},
				isValid:  true,
				errorMsg: "",
			}),

		Entry("Verification method has wrong id",
			TestCaseDIDStruct{
				did: &Did{
					Id: ValidTestDID,
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 InvalidTestDID,
							Type:               "JsonWebKey2020",
							Controller:         ValidTestDID,
							PublicKeyJwk:       ValidPublicKeyJWK,
							PublicKeyMultibase: "",
						},
					},
				},
				isValid:  false,
				errorMsg: "verification_method: (0: (id: unable to split did into method, namespace and id.).).",
			}),
		Entry(
			"Verification method has wrong controller",
			TestCaseDIDStruct{
				did: &Did{
					Id: ValidTestDID,
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "JsonWebKey2020",
							Controller:         InvalidTestDID,
							PublicKeyJwk:       ValidPublicKeyJWK,
							PublicKeyMultibase: "",
						},
					},
				},
				isValid:  false,
				errorMsg: "verification_method: (0: (controller: unable to split did into method, namespace and id.).).",
			}),
		Entry(
			"List of DIDs in cotroller is allowed",
			TestCaseDIDStruct{
				did: &Did{
					Id:         ValidTestDID,
					Controller: []string{ValidTestDID, ValidTestDID2},
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "Ed25519VerificationKey2020",
							Controller:         ValidTestDID,
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
				},
				isValid:  true,
				errorMsg: "",
			}),
		Entry(
			"List of DIDs in cotroller is not allowed",
			TestCaseDIDStruct{
				did: &Did{
					Context:    nil,
					Id:         ValidTestDID,
					Controller: []string{ValidTestDID, InvalidTestDID},
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "Ed25519VerificationKey2020",
							Controller:         ValidTestDID,
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
				},
				isValid:  false,
				errorMsg: "controller: (1: unable to split did into method, namespace and id.).",
			}),
		Entry(
			"Namespace in controler is not in list of allowed",
			TestCaseDIDStruct{
				did: &Did{
					Id:         ValidTestDID,
					Controller: []string{ValidTestDID},
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "Ed25519VerificationKey2020",
							Controller:         ValidTestDID,
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
				},
				allowedNamespaces: []string{"mainnet"},
				isValid:           false,
				errorMsg:          "controller: (0: did namespace must be one of: mainnet.); id: did namespace must be one of: mainnet; verification_method: (0: (controller: did namespace must be one of: mainnet; id: did namespace must be one of: mainnet.).).",
			}),
		Entry(
			"Controller is duplicated",
			TestCaseDIDStruct{
				did: &Did{
					Id:         ValidTestDID,
					Controller: []string{ValidTestDID, ValidTestDID},
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "Ed25519VerificationKey2020",
							Controller:         ValidTestDID,
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
				},
				isValid:  false,
				errorMsg: "controller: there should be no duplicates.",
			}),
		Entry(
			"Verification method is duplicated",
			TestCaseDIDStruct{
				did: &Did{
					Id: ValidTestDID,
					VerificationMethod: []*VerificationMethod{
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "Ed25519VerificationKey2020",
							Controller:         ValidTestDID,
							PublicKeyMultibase: ValidEd25519PubKey,
						},
						{
							Id:                 fmt.Sprintf("%s#fragment", ValidTestDID),
							Type:               "Ed25519VerificationKey2020",
							Controller:         ValidTestDID,
							PublicKeyMultibase: ValidEd25519PubKey,
						},
					},
				},
				isValid:  false,
				errorMsg: "verification_method: there are verification method duplicates.",
			}),
	)
})
