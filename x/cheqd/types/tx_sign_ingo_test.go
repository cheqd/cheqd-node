package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignInfoValidation(t *testing.T) {
	cases := []struct {
		name              string
		struct_           SignInfo
		allowedNamespaces []string
		isValid           bool
		errorMsg          string
	}{
		{
			name: "positive",
			struct_: SignInfo{
				VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
				Signature:            "aaa=",
			},
			isValid:           true,
			errorMsg:          "",
		},
		{
			name: "negative: namespace",
			struct_: SignInfo{
				VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#service1",
				Signature:            "DIDCommMessaging",
			},
			allowedNamespaces: []string{"mainnet"},
			isValid:           false,
			errorMsg:          "verification_method_id: did namespace must be one of: mainnet.",
		},
		{
			name: "negative: signature",
			struct_: SignInfo{
				VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#service1",
				Signature:            "!@#",
			},
			isValid:  false,
			errorMsg: "signature: must be encoded in Base64.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.struct_.Validate(tc.allowedNamespaces)

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, err.Error(), tc.errorMsg)
			}
		})
	}
}

func TestFullSignInfoDublicateValidation(t *testing.T) {
	cases := []struct {
		name              string
		structs_          []*SignInfo
		isValid           bool
	}{
		{
			name: "positive",
			structs_ : []*SignInfo{
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa="},
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "bbb="},
			},
			isValid:           true,
		},
		{
			name: "positive with all different pieces",
			structs_ : []*SignInfo{
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa="},
				{
					VerificationMethodId: "did:cheqd:bbbbbbbbbbbbbbbb#method1",
					Signature:            "bbb="},
			},
			isValid:           true,
		},
		{
			name: "negative",
			structs_ : []*SignInfo{
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa="},
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa="},
			},
			isValid:           false,
		},
		{
			name: "negative with a lot of same elems",
			structs_ : []*SignInfo{
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa="},
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa="},
				{
					VerificationMethodId: "did:cheqd:aaaaaaaaaaaaaaaa#method1",
					Signature:            "aaa="},
			},
			isValid:           false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res_ := IsUniqueSignInfoList(tc.structs_)
			require.Equal(t, res_, tc.isValid)
		})
	}
}
