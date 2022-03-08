package types

import (
	"github.com/stretchr/testify/require"
	"testing"
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
