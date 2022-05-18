package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	ValidTestDID          = "did:cheqd:testnet:123456789abcdefg"
	ValidTestDID2         = "did:cheqd:testnet:gfedcba987654321"
	InvalidTestDID        = "badDid"
	ValidEd25519PubKey    = "zF1hVGXXK9rmx5HhMTpGnGQJiab9qrFJbQXBRhSmYjQWX"
	NotValidEd25519PubKey = "zF1hVGXXK9rmx5HhMTpGnGQJi"
)

func TestDidValidation(t *testing.T) {
	cases := []struct {
		name              string
		struct_           *Did
		allowedNamespaces []string
		isValid           bool
		errorMsg          string
	}{
		{
			name: "Valid: Id: allowed DID",
			struct_: &Did{
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
		},
		{
			name: "Not valid: Id: not allowed DID",
			struct_: &Did{
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
		},
		{
			name: "Valid: Verification Method: all is fine with type Ed25519VerificationKey2020",
			struct_: &Did{
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
		},
		{
			name: "Valid: Verification Method: all is fine with type jwk",
			struct_: &Did{
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
		},
		{
			name: "Not valid: Verification Method: Wrong id",
			struct_: &Did{
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
		},
		{
			name: "Not valid: Verification Method: Wrong controller",
			struct_: &Did{
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
		},
		{
			name: "Valid: Controller: List of DIDs allowed",
			struct_: &Did{
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
		},
		{
			name: "Not valid: Controller: List of DIDs is not allowed",
			struct_: &Did{
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
		},
		{
			name: "Allowed namespaces: Negative",
			struct_: &Did{
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
		},
		{
			name: "Controller duplicated: negative",
			struct_: &Did{
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
		},
		{
			name: "VM duplicated: negative",
			struct_: &Did{
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
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.struct_.Validate(tc.allowedNamespaces)

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			}
		})
	}
}
