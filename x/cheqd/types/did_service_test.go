package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceValidation(t *testing.T) {
	cases := []struct {
		name              string
		struct_           Service
		baseDid           string
		allowedNamespaces []string
		isValid           bool
		errorMsg          string
	}{
		{
			name: "positive",
			struct_: Service{
				Id:              "did:cheqd:aaaaaaaaaaaaaaaa#service1",
				Type:            "DIDCommMessaging",
				ServiceEndpoint: "endpoint",
			},
			baseDid:           "did:cheqd:aaaaaaaaaaaaaaaa",
			allowedNamespaces: []string{""},
			isValid:           true,
			errorMsg:          "",
		},
		{
			name: "negative: namespace",
			struct_: Service{
				Id:              "did:cheqd:aaaaaaaaaaaaaaaa#service1",
				Type:            "DIDCommMessaging",
				ServiceEndpoint: "endpoint",
			},
			allowedNamespaces: []string{"mainnet"},
			isValid:           false,
			errorMsg:          "id: did namespace must be one of: mainnet.",
		},
		{
			name: "negative: base did",
			struct_: Service{
				Id:              "did:cheqd:aaaaaaaaaaaaaaaa#service1",
				Type:            "DIDCommMessaging",
				ServiceEndpoint: "endpoint",
			},
			baseDid:  "did:cheqd:baaaaaaaaaaaaaab",
			isValid:  false,
			errorMsg: "id: must have prefix: did:cheqd:baaaaaaaaaaaaaab.",
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

func UUIDTestCases() []struct {
	name        string
	did         string
	expectedDid string
} {
	return []struct {
		name        string
		did         string
		expectedDid string
	}{
		{
			name:        "base58 identifier - not changed",
			did:         "did:cheqd:aaaaaaaaaaaaaaaa",
			expectedDid: "did:cheqd:aaaaaaaaaaaaaaaa",
		},
		{
			name:        "Mixed case UUID",
			did:         "did:cheqd:test:BAbbba14-f294-458a-9b9c-474d188680fd",
			expectedDid: "did:cheqd:test:babbba14-f294-458a-9b9c-474d188680fd",
		},
		{
			name:        "Low case UUID",
			did:         "did:cheqd:test:babbba14-f294-458a-9b9c-474d188680fd",
			expectedDid: "did:cheqd:test:babbba14-f294-458a-9b9c-474d188680fd",
		},
		{
			name:        "Upper case UUID",
			did:         "did:cheqd:test:A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
			expectedDid: "did:cheqd:test:a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		},
	}
}

func TestUpdateUUIDForDID(t *testing.T) {
	for _, tc := range UUIDTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			result := UpdateUUIDForDID(tc.did)

			require.Equal(t, tc.expectedDid, result)
		})
	}
}

func TestUpdateUUIDIdentifiers(t *testing.T) {
	for _, tc := range UUIDTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			did := Did{
				Id:             tc.did,
				Authentication: []string{tc.did + "#key1"},
				VerificationMethod: []*VerificationMethod{
					{
						Id:         tc.did + "#key1",
						Type:       Ed25519VerificationKey2020,
						Controller: tc.did,
					},
				},
			}
			expectedDid := Did{
				Id:             tc.expectedDid,
				Authentication: []string{tc.expectedDid + "#key1"},
				VerificationMethod: []*VerificationMethod{
					{
						Id:         tc.expectedDid + "#key1",
						Type:       Ed25519VerificationKey2020,
						Controller: tc.expectedDid,
					},
				},
			}
			UpdateUUIDIdentifiers(&did)

			require.Equal(t, expectedDid, did)
		})
	}
}
