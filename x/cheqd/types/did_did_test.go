package types_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("DID Validation tests", func() {
	var struct_ *Did
	var allowedNamespaces []string
	var isValid bool
	var errorMsg string

	BeforeEach(func() {
		struct_ = &Did{}
		allowedNamespaces = []string{}
		isValid = false
		errorMsg = ""
	})

	AfterEach(func() {
		err := struct_.Validate(allowedNamespaces)

		if isValid {
			Expect(err).To(BeNil())
		} else {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(errorMsg))
		}
	})

	When("DID is allowed", func() {
		It("is valid", func() {
			struct_ = &Did{
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
			}
			isValid = true
			errorMsg = ""
		})
	})

	When("DID is not allowed", func() {
		It("should fail with error", func() {
			struct_ = &Did{
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
			}
			isValid = false
			errorMsg = "id: unable to split did into method, namespace and id; verification_method: (0: (id: must have prefix: badDid.).)."
		})
	})

	When("Verification method is Ed25519VerificationKey2020", func() {
		It("is valid", func() {
			struct_ = &Did{
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
			}
			isValid = true
			errorMsg = ""
		})
	})

	When("Verification method is jwk", func() {
		It("is valid", func() {
			struct_ = &Did{
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
			}
			isValid = true
			errorMsg = ""
		})
	})

	When("verification method has wrong id", func() {
		It("should fail with corresponding error", func() {
			struct_ = &Did{
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
			}
			isValid = false
			errorMsg = "verification_method: (0: (id: unable to split did into method, namespace and id.).)."
		})
	})
	When("verification method has wrong controller", func() {
		It("should fail with corresponding error", func() {
			struct_ = &Did{
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
			}
			isValid = false
			errorMsg = "verification_method: (0: (controller: unable to split did into method, namespace and id.).)."
		})
	})

	When("list of DIDs in cotroller is allowed", func() {
		It("s valid", func() {
			struct_ = &Did{
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
			}
			isValid = true
			errorMsg = ""
		})
	})

	When("list of DIDs in cotroller is not allowed", func() {
		It("should fail with corresponding error", func() {
			struct_ = &Did{
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
			}
			isValid = false
			errorMsg = "controller: (1: unable to split did into method, namespace and id.)."
		})
	})

	When("namespace in controleer is not in list of allowed", func() {
		It("should fail", func() {
			struct_ = &Did{
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
			}
			allowedNamespaces = []string{"mainnet"}
			isValid = false
			errorMsg = "controller: (0: did namespace must be one of: mainnet.); id: did namespace must be one of: mainnet; verification_method: (0: (controller: did namespace must be one of: mainnet; id: did namespace must be one of: mainnet.).)."
		})
	})

	When("controller is duplicated", func() {
		It("should fail", func() {
			struct_ = &Did{
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
			}
			isValid = false
			errorMsg = "controller: there should be no duplicates."
		})
	})

	When("verification method is duplicated", func() {
		It("should fail", func() {
			struct_ = &Did{
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
			}
			isValid = false
			errorMsg = "verification_method: there are verification method duplicates."
		})
	})
})
