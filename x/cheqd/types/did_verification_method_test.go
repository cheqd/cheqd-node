package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerificationMethodValidation(t *testing.T) {
	cases := []struct {
		name              string
		struct_           VerificationMethod
		baseDid           string
		allowedNamespaces []string
		isValid           bool
		errorMsg          string
	}{
		{
			name: "valid method with multibase key",
			struct_: VerificationMethod{
				Id:                 "did:cheqd:aaaaaaaaaaaaaaaa#qwe",
				Type:               "Ed25519VerificationKey2020",
				Controller:         "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk:       nil,
				PublicKeyMultibase: "multibase",
			},
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "valid method with jwk key",
			struct_: VerificationMethod{
				Id:         "did:cheqd:aaaaaaaaaaaaaaaa#rty",
				Type:       "JsonWebKey2020",
				Controller: "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: []*KeyValuePair{
					{
						Key:   "key",
						Value: "value",
					},
				},
				PublicKeyMultibase: "",
			},
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "base did: positive",
			struct_: VerificationMethod{
				Id:         "did:cheqd:aaaaaaaaaaaaaaaa#rty",
				Type:       "JsonWebKey2020",
				Controller: "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: []*KeyValuePair{
					{
						Key:   "key",
						Value: "value",
					},
				},
				PublicKeyMultibase: "",
			},
			baseDid:  "did:cheqd:aaaaaaaaaaaaaaaa",
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "base did: negative",
			struct_: VerificationMethod{
				Id:         "did:cheqd:aaaaaaaaaaaaaaaa#rty",
				Type:       "JsonWebKey2020",
				Controller: "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: []*KeyValuePair{
					{
						Key:   "key",
						Value: "value",
					},
				},
				PublicKeyMultibase: "",
			},
			baseDid:  "did:cheqd:bbbbbbbbbbbbbbbb",
			isValid:  false,
			errorMsg: "id: must have prefix: did:cheqd:bbbbbbbbbbbbbbbb.",
		},
		{
			name: "allowed namespaces: positive",
			struct_: VerificationMethod{
				Id:         "did:cheqd:mainnet:aaaaaaaaaaaaaaaa#rty",
				Type:       "JsonWebKey2020",
				Controller: "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: []*KeyValuePair{
					{
						Key:   "key",
						Value: "value",
					},
				},
				PublicKeyMultibase: "",
			},
			allowedNamespaces: []string{"mainnet", ""},
			isValid:           true,
		},
		{
			name: "allowed namespaces: positive",
			struct_: VerificationMethod{
				Id:         "did:cheqd:mainnet:aaaaaaaaaaaaaaaa#rty",
				Type:       "JsonWebKey2020",
				Controller: "did:cheqd:bbbbbbbbbbbbbbbb",
				PublicKeyJwk: []*KeyValuePair{
					{
						Key:   "key",
						Value: "value",
					},
				},
				PublicKeyMultibase: "",
			},
			allowedNamespaces: []string{"testnet"},
			isValid:           false,
			errorMsg:          "controller: did namespace must be one of: testnet; id: did namespace must be one of: testnet.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.struct_.Validate(tc.baseDid, tc.allowedNamespaces)

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, err.Error(), tc.errorMsg)
			}
		})
	}
}
